package integration

import (
    "bytes"
    "encoding/json"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "testing"
)

func repoRoot(t *testing.T) string {
    _, file, _, ok := runtime.Caller(0)
    if !ok { t.Fatal("no caller") }
    // this file: <repo>/context-keeper/integration/context_cli_test.go
    dir := filepath.Dir(file)
    // go up one to repo root of module
    return filepath.Dir(dir)
}

func run(t *testing.T, dir string, args ...string) (string, string, error) {
    t.Helper()
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Dir = dir
    var out, errB bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errB
    e := cmd.Run()
    return out.String(), errB.String(), e
}

func TestCLI_EndToEnd_JSONAndHuman(t *testing.T) {
    root := repoRoot(t)
    tmp := t.TempDir()

    // Build test binary in repo root
    bin := filepath.Join(root, "context-integration-test.exe")
    if out, errS, err := run(t, root, "go", "build", "-o", bin, "./cmd/context"); err != nil {
        t.Fatalf("build failed: %v, stdout=%s, stderr=%s", err, out, errS)
    }

    // init --json in temp working dir
    out, errS, err := run(t, tmp, bin, "init", "--json")
    if err != nil { t.Fatalf("init failed: %v, stderr=%s", err, errS) }
    var initResp map[string]any
    if err := json.Unmarshal([]byte(out), &initResp); err != nil { t.Fatalf("init json: %v, out=%s", err, out) }
    if initResp["success"] != true { t.Fatalf("init success false: %v", initResp) }

    // save (human)
    out, errS, err = run(t, tmp, bin, "save", "--key", "main_file", "--value", "src/index.js", "--priority", "high")
    if err != nil { t.Fatalf("save failed: %v, stderr=%s", err, errS) }
    if !strings.Contains(out, "Context saved for key") { t.Fatalf("unexpected save output: %s", out) }

    // get --json
    out, errS, err = run(t, tmp, bin, "get", "--key", "main_file", "--json")
    if err != nil { t.Fatalf("get failed: %v, stderr=%s", err, errS) }
    var getResp map[string]any
    if err := json.Unmarshal([]byte(out), &getResp); err != nil { t.Fatalf("get json: %v, out=%s", err, out) }
    if getResp["value"] != "src/index.js" { t.Fatalf("get value mismatch: %v", getResp) }

    // list --json
    out, errS, err = run(t, tmp, bin, "list", "--json")
    if err != nil { t.Fatalf("list failed: %v, stderr=%s", err, errS) }
    var listResp []map[string]any
    if err := json.Unmarshal([]byte(out), &listResp); err != nil { t.Fatalf("list json: %v, out=%s", err, out) }
    if len(listResp) == 0 { t.Fatalf("list empty: %v", listResp) }

    // list human (contains headers)
    out, errS, err = run(t, tmp, bin, "list")
    if err != nil { t.Fatalf("list human failed: %v, stderr=%s", err, errS) }
    if !strings.Contains(out, "CHANNEL") || !strings.Contains(out, "KEY") {
        t.Fatalf("list human missing headers: %s", out)
    }

    // get human contains the value
    out, errS, err = run(t, tmp, bin, "get", "--key", "main_file")
    if err != nil { t.Fatalf("get human failed: %v, stderr=%s", err, errS) }
    if !strings.Contains(out, "src/index.js") {
        t.Fatalf("get human missing value: %s", out)
    }

    // delete --json
    out, errS, err = run(t, tmp, bin, "delete", "--key", "main_file", "--json")
    if err != nil { t.Fatalf("delete failed: %v, stderr=%s", err, errS) }
    var delResp map[string]any
    if err := json.Unmarshal([]byte(out), &delResp); err != nil { t.Fatalf("delete json: %v, out=%s", err, out) }
    if delResp["action"] != "deleted" { t.Fatalf("delete action mismatch: %v", delResp) }

    // ensure DB and files are inside tmp/.context-keeper
    if _, err := os.Stat(filepath.Join(tmp, ".context-keeper", "context.db")); err != nil { t.Fatalf("missing db: %v", err) }
    if _, err := os.Stat(filepath.Join(tmp, ".context-keeper", "AGENT_INSTRUCTIONS.md")); err != nil { t.Fatalf("missing instructions: %v", err) }

    // cleanup binary
    _ = os.Remove(bin)
}
