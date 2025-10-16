package db

import (
    "os"
    "path/filepath"
    "testing"
)

func TestRepositoryCRUD(t *testing.T) {
    dir := t.TempDir()
    dbPath := filepath.Join(dir, "context.db")
    repo, err := NewRepository(dbPath)
    if err != nil { t.Fatalf("NewRepository error: %v", err) }
    t.Cleanup(func(){ _ = repo.Close() })

    // Save
    if err := repo.SaveItem("k1", "v1", "general", "normal"); err != nil {
        t.Fatalf("SaveItem: %v", err)
    }

    // Get
    it, err := repo.GetItemByKey("k1")
    if err != nil { t.Fatalf("GetItemByKey: %v", err) }
    if it.Value != "v1" || it.Channel != "general" || it.Priority != "normal" {
        t.Fatalf("unexpected item: %+v", it)
    }

    // Update
    if err := repo.SaveItem("k1", "v1b", "testing", "high"); err != nil {
        t.Fatalf("SaveItem update: %v", err)
    }
    it2, err := repo.GetItemByKey("k1")
    if err != nil { t.Fatalf("GetItemByKey after update: %v", err) }
    if it2.Value != "v1b" || it2.Channel != "testing" || it2.Priority != "high" {
        t.Fatalf("unexpected updated item: %+v", it2)
    }

    // List
    items, err := repo.ListItems(nil, nil)
    if err != nil { t.Fatalf("ListItems: %v", err) }
    if len(items) != 1 || items[0].Key != "k1" {
        t.Fatalf("unexpected list: %+v", items)
    }

    ch := "testing"
    lim := 1
    items2, err := repo.ListItems(&ch, &lim)
    if err != nil { t.Fatalf("ListItems filtered: %v", err) }
    if len(items2) != 1 || items2[0].Channel != "testing" {
        t.Fatalf("unexpected filtered list: %+v", items2)
    }

    // Delete
    if err := repo.DeleteItem("k1"); err != nil { t.Fatalf("DeleteItem: %v", err) }
    if _, err := repo.GetItemByKey("k1"); err == nil {
        t.Fatalf("expected error after delete")
    }

    // Ensure DB file exists in temp dir only
    if _, err := os.Stat(dbPath); err != nil {
        t.Fatalf("db file missing: %v", err)
    }
}
