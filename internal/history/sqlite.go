package history

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func OpenDB(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("history: open db: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, tagline TEXT DEFAULT '',
		path TEXT NOT NULL DEFAULT '', provider TEXT NOT NULL DEFAULT '',
		model TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY, project_id TEXT NOT NULL, phase TEXT DEFAULT '',
		role TEXT NOT NULL, content TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);
	CREATE TABLE IF NOT EXISTS provider_models (
		provider TEXT NOT NULL, model_id TEXT NOT NULL,
		display_name TEXT NOT NULL DEFAULT '',
		cached_at TEXT NOT NULL DEFAULT (datetime('now')),
		PRIMARY KEY (provider, model_id)
	);
	CREATE INDEX IF NOT EXISTS idx_conversations_project ON conversations(project_id, created_at);
	CREATE INDEX IF NOT EXISTS idx_projects_updated ON projects(updated_at DESC);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *Store) CreateProject(p Project) error {
	_, err := s.db.Exec(
		`INSERT INTO projects (id, name, tagline, path, provider, model) VALUES (?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Tagline, p.Path, p.Provider, p.Model,
	)
	return err
}

func (s *Store) UpdateProject(p Project) error {
	_, err := s.db.Exec(
		`UPDATE projects SET name=?, tagline=?, path=?, provider=?, model=?, updated_at=datetime('now') WHERE id=?`,
		p.Name, p.Tagline, p.Path, p.Provider, p.Model, p.ID,
	)
	return err
}

func (s *Store) ListProjects() ([]Project, error) {
	rows, err := s.db.Query(
		`SELECT id, name, tagline, path, provider, model, created_at, updated_at FROM projects ORDER BY updated_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Tagline, &p.Path, &p.Provider, &p.Model, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (s *Store) GetProject(id string) (*Project, error) {
	var p Project
	err := s.db.QueryRow(
		`SELECT id, name, tagline, path, provider, model, created_at, updated_at FROM projects WHERE id=?`, id,
	).Scan(&p.ID, &p.Name, &p.Tagline, &p.Path, &p.Provider, &p.Model, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) AddMessage(m Message) error {
	_, err := s.db.Exec(
		`INSERT INTO conversations (id, project_id, phase, role, content) VALUES (?, ?, ?, ?, ?)`,
		m.ID, m.ProjectID, m.Phase, m.Role, m.Content,
	)
	return err
}

func (s *Store) GetMessages(projectID string) ([]Message, error) {
	rows, err := s.db.Query(
		`SELECT id, project_id, phase, role, content, created_at FROM conversations WHERE project_id=? ORDER BY created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Phase, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (s *Store) DeleteProject(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM conversations WHERE project_id = ?`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM projects WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) Close() error { return s.db.Close() }

type ProviderModel struct {
	Provider    string `json:"provider"`
	ModelID     string `json:"model_id"`
	DisplayName string `json:"display_name"`
}

func (s *Store) UpsertProviderModel(provider, modelID, displayName string) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO provider_models (provider, model_id, display_name, cached_at) VALUES (?, ?, ?, datetime('now'))`,
		provider, modelID, displayName,
	)
	return err
}

func (s *Store) ListCachedModels(provider string) ([]ProviderModel, error) {
	rows, err := s.db.Query(
		`SELECT provider, model_id, display_name FROM provider_models WHERE provider = ? ORDER BY model_id`,
		provider,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []ProviderModel
	for rows.Next() {
		var m ProviderModel
		if err := rows.Scan(&m.Provider, &m.ModelID, &m.DisplayName); err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, rows.Err()
}
