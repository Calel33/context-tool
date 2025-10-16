package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"

    dbpkg "context-keeper/internal/db"
    cipkg "context-keeper/internal/cli"
)

func main() {
    if len(os.Args) < 2 {
        printRootHelp()
        os.Exit(1)
    }

    sub := os.Args[1]
    switch sub {
    case "init":
        fs := flag.NewFlagSet("init", flag.ExitOnError)
        jsonOut := fs.Bool("json", false, "JSON output")
        _ = fs.Parse(os.Args[2:])
        exitIfErr(cipkg.HandleInit(*jsonOut))
    case "save":
        fs := flag.NewFlagSet("save", flag.ExitOnError)
        key := fs.String("key", "", "Key")
        value := fs.String("value", "", "Value")
        channel := fs.String("channel", "general", "Channel")
        priority := fs.String("priority", "normal", "Priority (high|normal|low)")
        jsonOut := fs.Bool("json", false, "JSON output")
        _ = fs.Parse(os.Args[2:])
        if *key == "" || *value == "" {
            exitIfErr(errors.New("--key and --value are required"))
        }
        repo, err := openRepoInCWD()
        exitIfErr(err)
        defer repo.Close()
        exitIfErr(cipkg.HandleSave(context.Background(), repo, *key, *value, *channel, *priority, *jsonOut))
    case "get":
        fs := flag.NewFlagSet("get", flag.ExitOnError)
        key := fs.String("key", "", "Key")
        jsonOut := fs.Bool("json", false, "JSON output")
        _ = fs.Parse(os.Args[2:])
        if *key == "" {
            exitIfErr(errors.New("--key is required"))
        }
        repo, err := openRepoInCWD()
        exitIfErr(err)
        defer repo.Close()
        exitIfErr(cipkg.HandleGet(context.Background(), repo, *key, *jsonOut))
    case "list":
        fs := flag.NewFlagSet("list", flag.ExitOnError)
        channel := fs.String("channel", "", "Channel filter")
        limit := fs.Int("limit", 0, "Max items")
        jsonOut := fs.Bool("json", false, "JSON output")
        _ = fs.Parse(os.Args[2:])
        repo, err := openRepoInCWD()
        exitIfErr(err)
        defer repo.Close()
        var chPtr *string
        if *channel != "" {
            chPtr = channel
        }
        var limPtr *int
        if *limit > 0 {
            limPtr = limit
        }
        exitIfErr(cipkg.HandleList(context.Background(), repo, chPtr, limPtr, *jsonOut))
    case "delete":
        fs := flag.NewFlagSet("delete", flag.ExitOnError)
        key := fs.String("key", "", "Key")
        jsonOut := fs.Bool("json", false, "JSON output")
        _ = fs.Parse(os.Args[2:])
        if *key == "" {
            exitIfErr(errors.New("--key is required"))
        }
        repo, err := openRepoInCWD()
        exitIfErr(err)
        defer repo.Close()
        exitIfErr(cipkg.HandleDelete(context.Background(), repo, *key, *jsonOut))
    case "--help", "-h", "help":
        printRootHelp()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", sub)
        printRootHelp()
        os.Exit(1)
    }
}

func exitIfErr(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }
}

func openRepoInCWD() (*dbpkg.Repository, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    dbPath := filepath.Join(cwd, ".context-keeper", "context.db")
    return dbpkg.NewRepository(dbPath)
}

func printRootHelp() {
    fmt.Print("context <command> [options]\n\nCommands:\n  init\n  save --key k --value v [--channel c] [--priority p]\n  get --key k [--json]\n  list [--channel c] [--limit N] [--json]\n  delete --key k\n")
}
