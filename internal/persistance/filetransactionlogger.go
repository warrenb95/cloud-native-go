package persistance

import (
	"bufio"
	"fmt"
	"os"
)

type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log filr: %w", err)
	}

	return &FileTransactionLogger{
		file: file,
	}, nil
}

func (f *FileTransactionLogger) Run() {
	// Can't just assign the chan to the struct as the channels in the struct are one way.
	// i.e. can't do f.events = make(chan<- Event, 16) because we can't then recieve on this chan below.
	events := make(chan Event, 16)
	f.events = events
	errors := make(chan error, 1)
	f.errors = errors

	defer close(events)
	defer close(errors)
	defer f.file.Close()

	go func() {
		for e := range events {
			f.lastSequence++
			_, err := fmt.Fprintf(f.file, "%d\t%d\t%s\t%s\n",
				f.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (f *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(f.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event // Use value and not pointer so we dont't hace to recreate in for loop.

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if f.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction out of sequence order")
				return
			}
			f.lastSequence = e.Sequence
			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func (f *FileTransactionLogger) WritePut(key string, value interface{}) {
	f.events <- Event{EventType: EventPut, Key: key, Value: value.([]byte)}
}

func (f *FileTransactionLogger) WriteDelete(key string) {
	f.events <- Event{EventType: EventDelete, Key: key}
}

func (f *FileTransactionLogger) Err() <-chan error {
	return f.errors
}
