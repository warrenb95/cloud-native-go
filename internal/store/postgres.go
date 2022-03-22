package store

import (
	"database/sql"
	"fmt"
)

type PostgresTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	db     *sql.DB
}

type PostgresConfig struct {
	DBName   string
	Host     string
	User     string
	Password string
}

func NewPostgresTransactionLogger(params PostgresConfig) (*PostgresTransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		params.Host, params.DBName, params.User, params.Password)

	db, err := sql.Open("postgress", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{
		db: db,
	}

	exists, err := logger.tableExist()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}

	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions
		(event_type, key, value)
		VALUES($1, $2, $3)`

		for e := range events {
			tx, err := l.db.Begin()
			if err != nil {
				continue
			}
			defer tx.Rollback()

			_, err = tx.Exec(
				query,
				e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}

			tx.Commit()
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outError)
		defer close(outEvent)

		query := `SELECT sequence, event_type, key, value FROM transactions
		ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- err
			return
		}

		defer rows.Close()
		e := Event{}
		for rows.Next() {
			err = rows.Scan(
				&e.Sequence, &e.EventType, &e.Key, &e.Value,
			)
			if err != nil {
				outError <- err
			}

			outEvent <- e
		}

		if rows.Err() != nil {
			outError <- rows.Err()
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) tableExist() (bool, error) {
	const table = "transactions"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	query := `CREATE TABLE transactions {
			sequence BIGSERIAL PRIMARY KEY,
			event_type INT64 NOT NULL,
			key TEXT NOT NULL,
			value TEXT NOT NULL
		};`

	_, err := l.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create new table: %w", err)
	}

	return nil
}

func (l *PostgresTransactionLogger) WritePut(key string, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}
