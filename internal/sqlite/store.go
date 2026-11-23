// Package sqlite is the only place that talks to the database. Pure-Go driver, no cgo, no ORM.
package sqlite

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"flow/internal/domain/flow"
	_ "modernc.org/sqlite"
)

//go:embed migrations/0001_init.sql
var migration0001 string

// Store is the persistence boundary.
type Store struct {
	db *sql.DB
	// Ports is injected so validation can check port declarations without this
	// package importing the node registry. Keeps dependencies pointing inward.
	Ports flow.PortNamer
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(migration0001); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Save(f *flow.Flow) error {
	if err := f.Validate(s.Ports); err != nil {
		return err
	}
	f.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO flows (id, name, updated_at, data) VALUES (?,?,?,?)
		 ON CONFLICT(id) DO UPDATE SET name=excluded.name, updated_at=excluded.updated_at, data=excluded.data`,
		f.ID, f.Name, f.UpdatedAt, string(data),
	)
	return err
}

func (s *Store) Get(id string) (*flow.Flow, error) {
	var data string
	if err := s.db.QueryRow(`SELECT data FROM flows WHERE id = ?`, id).Scan(&data); err != nil {
		return nil, err
	}
	var f flow.Flow
	return &f, json.Unmarshal([]byte(data), &f)
}

func (s *Store) List() ([]flow.Flow, error) {
	rows, err := s.db.Query(`SELECT data FROM flows ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []flow.Flow
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		var f flow.Flow
		if err := json.Unmarshal([]byte(data), &f); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (s *Store) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM flows WHERE id = ?`, id)
	return err
}

// RecordRun persists a finished run, so history survives a restart.
func (s *Store) RecordRun(runID, flowID, status string) error {
	_, err := s.db.Exec(
		`INSERT INTO runs (id, flow_id, status, started_at, finished_at) VALUES (?,?,?,?,?)
		 ON CONFLICT(id) DO UPDATE SET status=excluded.status, finished_at=excluded.finished_at`,
		runID, flowID, status, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339),
	)
	return err
}
