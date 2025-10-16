package db

import (
    "database/sql"
    "errors"
    _ "modernc.org/sqlite"
    "time"
)

type Item struct {
    Key      string
    Value    string
    Channel  string
    Priority string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Repository struct {
    db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
    // modernc.org/sqlite driver name is "sqlite"; use file: DSN for portability
    dsn := "file:" + dbPath
    db, err := sql.Open("sqlite", dsn)
    if err != nil { return nil, err }
    if err := createSchema(db); err != nil { _ = db.Close(); return nil, err }
    return &Repository{db: db}, nil
}

func createSchema(db *sql.DB) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS context_items (
        key TEXT PRIMARY KEY,
        value TEXT NOT NULL,
        channel TEXT NOT NULL DEFAULT 'general',
        priority TEXT NOT NULL DEFAULT 'normal',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`)
    return err
}

func (r *Repository) Close() error { return r.db.Close() }

func (r *Repository) SaveItem(key, value, channel, priority string) error {
    if key == "" { return errors.New("key required") }
    _, err := r.db.Exec(`INSERT INTO context_items(key, value, channel, priority)
        VALUES(?, ?, ?, ?)
        ON CONFLICT(key) DO UPDATE SET value=excluded.value, channel=excluded.channel, priority=excluded.priority, updated_at=CURRENT_TIMESTAMP`,
        key, value, channel, priority)
    return err
}

func (r *Repository) GetItemByKey(key string) (*Item, error) {
    row := r.db.QueryRow(`SELECT key, value, channel, priority, created_at, updated_at FROM context_items WHERE key = ?`, key)
    var it Item
    var created, updated string
    if err := row.Scan(&it.Key, &it.Value, &it.Channel, &it.Priority, &created, &updated); err != nil { return nil, err }
    it.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
    it.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
    return &it, nil
}

func (r *Repository) ListItems(channel *string, limit *int) ([]Item, error) {
    q := `SELECT key, value, channel, priority, created_at, updated_at FROM context_items`
    args := []any{}
    if channel != nil {
        q += ` WHERE channel = ?`
        args = append(args, *channel)
    }
    q += ` ORDER BY updated_at DESC`
    if limit != nil && *limit > 0 {
        q += ` LIMIT ?`
        args = append(args, *limit)
    }
    rows, err := r.db.Query(q, args...)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Item
    for rows.Next() {
        var it Item
        var created, updated string
        if err := rows.Scan(&it.Key, &it.Value, &it.Channel, &it.Priority, &created, &updated); err != nil { return nil, err }
        out = append(out, it)
    }
    return out, rows.Err()
}

func (r *Repository) DeleteItem(key string) error {
    _, err := r.db.Exec(`DELETE FROM context_items WHERE key = ?`, key)
    return err
}
