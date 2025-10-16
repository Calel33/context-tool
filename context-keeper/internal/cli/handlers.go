package cli

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"

    dbpkg "context-keeper/internal/db"
    "context-keeper/internal/tui"
)

func HandleInit(jsonOut bool) error {
    cwd, err := os.Getwd()
    if err != nil { return err }
    dir := filepath.Join(cwd, ".context-keeper")
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        if err := os.MkdirAll(dir, 0o755); err != nil { return err }
    }
    dbPath := filepath.Join(dir, "context.db")
    repo, err := dbpkg.NewRepository(dbPath)
    if err != nil { return err }
    _ = repo.Close()

    instrPath := filepath.Join(dir, "AGENT_INSTRUCTIONS.md")
    if _, err := os.Stat(instrPath); os.IsNotExist(err) {
        // generate from template if present
        cwdTpl := filepath.Join(cwd, "AGENT_INSTRUCTIONS.md.template")
        content := defaultInstructions()
        if b, err := os.ReadFile(cwdTpl); err == nil { content = string(b) }
        if err := os.WriteFile(instrPath, []byte(content), 0o644); err != nil { return err }
    }
    if jsonOut {
        out := map[string]any{"success": true, "message": "Context Keeper initialized successfully.", "storagePath": filepath.ToSlash(filepath.Join("./.context-keeper", "context.db")), "instructionsPath": filepath.ToSlash(filepath.Join("./.context-keeper", "AGENT_INSTRUCTIONS.md"))}
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        return enc.Encode(out)
    }
    fmt.Println(tui.Success("Context Keeper initialized successfully."))
    fmt.Println("- Storage created at: ./.context-keeper/context.db")
    fmt.Println("- Instructions for your AI agent are in: ./.context-keeper/AGENT_INSTRUCTIONS.md")
    return nil
}

func HandleSave(ctx context.Context, repo *dbpkg.Repository, key, value, channel, priority string, jsonOut bool) error {
    if err := repo.SaveItem(key, value, channel, priority); err != nil { return err }
    if jsonOut {
        out := map[string]any{"success": true, "key": key, "action": "saved"}
        return json.NewEncoder(os.Stdout).Encode(out)
    }
    fmt.Println(tui.Success(fmt.Sprintf("Context saved for key: '%s'", key)))
    return nil
}

func HandleGet(ctx context.Context, repo *dbpkg.Repository, key string, jsonOut bool) error {
    it, err := repo.GetItemByKey(key)
    if err != nil { return err }
    if jsonOut {
        out := map[string]any{"key": it.Key, "value": it.Value, "channel": it.Channel, "priority": it.Priority}
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        return enc.Encode(out)
    }
    fmt.Printf("Value for key '%s':\n%s\n", tui.Highlight(it.Key), it.Value)
    return nil
}

func HandleList(ctx context.Context, repo *dbpkg.Repository, channel *string, limit *int, jsonOut bool) error {
    items, err := repo.ListItems(channel, limit)
    if err != nil { return err }
    if jsonOut {
        type outItem struct{ Key, Value, Channel, Priority string }
        arr := make([]outItem, 0, len(items))
        for _, it := range items { arr = append(arr, outItem{it.Key, it.Value, it.Channel, it.Priority}) }
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        return enc.Encode(arr)
    }
    return tui.RenderList(items)
}

func HandleDelete(ctx context.Context, repo *dbpkg.Repository, key string, jsonOut bool) error {
    if key == "" { return errors.New("key required") }
    if err := repo.DeleteItem(key); err != nil { return err }
    if jsonOut {
        out := map[string]any{"success": true, "key": key, "action": "deleted"}
        return json.NewEncoder(os.Stdout).Encode(out)
    }
    fmt.Println(tui.Deleted(fmt.Sprintf("Context for key '%s' deleted.", key)))
    return nil
}

func defaultInstructions() string {
    return "# Context Keeper Instructions\n\nUse 'context init|save|get|list|delete'. Prefer --json for machine parsing."
}
