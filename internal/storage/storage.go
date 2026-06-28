// Package storage хранит историю сроков истечения доменов в локальной SQLite-базе,
// расположенной рядом с бинарником. Это позволяет замечать продление домена:
// когда срок оплаты увеличивается по сравнению с предыдущей проверкой.
package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DefaultFileName - имя файла базы данных рядом с бинарником.
const DefaultFileName = "kz-domain-monitor.db"

// Store - хранилище истории сроков истечения доменов.
type Store struct {
	db *sql.DB
}

// DefaultPath возвращает путь к файлу базы данных рядом с исполняемым файлом.
// При недоступности пути исполняемого файла используется текущая директория.
func DefaultPath() string {
	exe, err := os.Executable()
	if err != nil {
		return DefaultFileName
	}
	return filepath.Join(filepath.Dir(exe), DefaultFileName)
}

// Open открывает (создавая при необходимости) базу данных по указанному пути.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS domain_history (
		name TEXT PRIMARY KEY,
		expiration_date TEXT NOT NULL,
		updated_at TEXT NOT NULL
	)`); err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close закрывает базу данных.
func (s *Store) Close() error {
	return s.db.Close()
}

// GetExpiration возвращает сохранённый срок истечения домена.
// Если записи нет, возвращается (nil, nil).
func (s *Store) GetExpiration(name string) (*time.Time, error) {
	var raw string
	err := s.db.QueryRow(`SELECT expiration_date FROM domain_history WHERE name = ?`, name).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

// SetExpiration сохраняет (или обновляет) срок истечения домена.
func (s *Store) SetExpiration(name string, date time.Time) error {
	_, err := s.db.Exec(`INSERT INTO domain_history (name, expiration_date, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET expiration_date = excluded.expiration_date, updated_at = excluded.updated_at`,
		name, date.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	return err
}
