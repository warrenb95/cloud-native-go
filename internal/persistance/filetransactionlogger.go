package persistance

import "os"

type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
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
