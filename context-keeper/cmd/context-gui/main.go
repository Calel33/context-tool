package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/paginator"
    "github.com/charmbracelet/bubbles/viewport"

    dbpkg "context-keeper/internal/db"
)

// GUI app with list view and basic filters (channel, priority, key search).

type model struct {
    repo        *dbpkg.Repository
    projectPath string
    initErr     error

    // data
    allItems []dbpkg.Item

    // ui state
    tbl       table.Model
    allRows   []table.Row
    search    textinput.Model
    channelIx int // 0 = all
    priorityIx int // 0 = all
    width, height int
    tblHeight int

    // modes: list | edit | confirmDelete
    mode string

    // edit state
    editExisting bool
    keyInput     textinput.Model
    valueInput   textarea.Model
    editChannelIx  int
    editPriorityIx int

    // notifications
    toast string

    // help/keys
    help help.Model
    keys keymap

    // extras
    loadSpin spinner.Model
    pager    paginator.Model
    preview  viewport.Model
    previewActive bool
}

var (
    channels = []string{"all", "general"}
    priorities = []string{"all", "high", "normal", "low"}
)

func initialModel() model {
    cwd, err := os.Getwd()
    if err != nil { return model{initErr: err} }
    dbPath := filepath.Join(cwd, ".context-keeper", "context.db")
    repo, err := dbpkg.NewRepository(dbPath)
    if err != nil { return model{projectPath: cwd, initErr: err} }

    cols := []table.Column{{Title: bold.Render("CHANNEL"), Width: 12}, {Title: bold.Render("PRIORITY"), Width: 10}, {Title: bold.Render("KEY"), Width: 28}, {Title: bold.Render("UPDATED"), Width: 19}}
    tbl := table.New(table.WithColumns(cols))
    styles := table.Styles{
        Header: lipgloss.NewStyle().Foreground(colorPrimary).Bold(true),
        Selected: lipgloss.NewStyle().Foreground(colorPrimary),
    }
    tbl.SetStyles(styles)
    tbl.SetHeight(7)

    ti := textinput.New()
    ti.Placeholder = "Search key... (/ to focus)"
    ti.CharLimit = 128

    // editor inputs
    ki := textinput.New()
    ki.Placeholder = "Key"
    ki.CharLimit = 128

    ta := textarea.New()
    ta.Placeholder = "Value (multiline)"
    ta.SetHeight(8)

    sp := spinner.New()
    sp.Spinner = spinner.Dot
    pg := paginator.New()
    pg.Type = paginator.Dots
    pg.PerPage = 10
    m := model{repo: repo, projectPath: cwd, tbl: tbl, search: ti, keyInput: ki, valueInput: ta, mode: "list", loadSpin: sp, pager: pg}
    m.help = help.New()
    m.keys = newKeymap()
    m.preview = viewport.New(40, 10)
    // initial load
    items, _ := repo.ListItems(nil, nil)
    m.allItems = items
    m.refreshRows()
    return m
}

func (m *model) refreshRows() {
    var rows []table.Row
    ch := ""
    if m.channelIx > 0 && m.channelIx < len(channels) { ch = channels[m.channelIx] }
    pr := ""
    if m.priorityIx > 0 && m.priorityIx < len(priorities) { pr = priorities[m.priorityIx] }
    q := strings.TrimSpace(strings.ToLower(m.search.Value()))
    for _, it := range m.allItems {
        if ch != "" && it.Channel != ch { continue }
        if pr != "" && strings.ToLower(it.Priority) != pr { continue }
        if q != "" && !strings.Contains(strings.ToLower(it.Key), q) { continue }
        rows = append(rows, table.Row{it.Channel, it.Priority, it.Key, it.UpdatedAt.Format("2006-01-02 15:04:05")})
    }
    m.allRows = rows
    m.applyPagination()
}

func (m model) Init() tea.Cmd { return nil }

type toastClearMsg struct{}

func clearToastAfter(d time.Duration) tea.Cmd {
    return tea.Tick(d, func(time.Time) tea.Msg { return toastClearMsg{} })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case reloadDoneMsg:
        if msg.err != nil { m.toast = "Error: " + msg.err.Error(); return m, clearToastAfter(2*time.Second) }
        m.allItems = msg.items
        m.refreshRows()
        m.toast = "Reloaded"
        return m, clearToastAfter(2*time.Second)
    case openDoneMsg:
        if msg.err != nil { m.toast = "Error: " + msg.err.Error(); return m, clearToastAfter(2*time.Second) }
        // reopen repo now
        newDB := filepath.Join(msg.project, ".context-keeper", "context.db")
        if m.repo != nil { _ = m.repo.Close() }
        repo, err := dbpkg.NewRepository(newDB)
        if err != nil { m.toast = "Error: " + err.Error(); return m, clearToastAfter(2*time.Second) }
        m.repo = repo
        m.projectPath = msg.project
        m.allItems = msg.items
        m.refreshRows()
        m.mode = "list"
        m.toast = "Project opened"
        return m, clearToastAfter(2*time.Second)
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.tblHeight = max(5, m.height-10)
        m.tbl.SetHeight(m.tblHeight)
        m.preview.Width = max(30, m.width-8)
        m.preview.Height = max(8, m.height-8)
        m.applyPagination()
        return m, nil
    case toastClearMsg:
        m.toast = ""
        return m, nil
    case tea.KeyMsg:
        if m.mode == "edit" {
            // route to inputs
            switch msg.String() {
            case "tab":
                if m.keyInput.Focused() { m.keyInput.Blur(); m.valueInput.Focus() } else { m.valueInput.Blur(); m.keyInput.Focus() }
                return m, nil
            case "[": // cycle channel back
                m.editChannelIx = (m.editChannelIx + len(channels) - 1) % len(channels)
                return m, nil
            case "]": // cycle channel forward
                m.editChannelIx = (m.editChannelIx + 1) % len(channels)
                return m, nil
            case "-": // cycle priority back
                m.editPriorityIx = (m.editPriorityIx + len(priorities) - 1) % len(priorities)
                return m, nil
            case "=": // cycle priority forward
                m.editPriorityIx = (m.editPriorityIx + 1) % len(priorities)
                return m, nil
            case "s": // save
                key := strings.TrimSpace(m.keyInput.Value())
                val := m.valueInput.Value()
                ch := channels[max(1, m.editChannelIx)] // avoid "all"
                pr := priorities[max(1, m.editPriorityIx)]
                if key != "" && val != "" {
                    if err := m.repo.SaveItem(key, val, ch, pr); err != nil {
                        m.toast = "Error: " + err.Error()
                        return m, clearToastAfter(2*time.Second)
                    } else {
                        items, _ := m.repo.ListItems(nil, nil)
                        m.allItems = items
                        m.refreshRows()
                        m.mode = "list"
                        m.toast = "Saved"
                        return m, clearToastAfter(2*time.Second)
                    }
                } else {
                    m.toast = "Key and Value required"
                    return m, clearToastAfter(2*time.Second)
                }
                return m, nil
            case "esc":
                m.mode = "list"
                return m, nil
            }
            var cmd tea.Cmd
            if m.keyInput.Focused() { m.keyInput, cmd = m.keyInput.Update(msg) } else { m.valueInput, cmd = m.valueInput.Update(msg) }
            return m, cmd
        }
        if m.mode == "confirmDelete" {
            switch msg.String() {
            case "y":
                if row := m.tbl.SelectedRow(); row != nil {
                    key := row[2]
                    if err := m.repo.DeleteItem(key); err != nil { m.toast = "Error: " + err.Error(); m.mode = "list"; return m, clearToastAfter(2*time.Second) } else {
                        items, _ := m.repo.ListItems(nil, nil)
                        m.allItems = items
                        m.refreshRows()
                        m.toast = "Deleted"
                        m.mode = "list"
                        return m, clearToastAfter(2*time.Second)
                    }
                }
                m.mode = "list"
                return m, nil
            case "n", "esc":
                m.mode = "list"
                return m, nil
            }
            return m, nil
        }
        if m.search.Focused() {
            var cmd tea.Cmd
            m.search, cmd = m.search.Update(msg)
            m.refreshRows()
            if msg.Type == tea.KeyEsc { m.search.Blur() }
            return m, cmd
        }
        switch msg.String() {
        case "p": // toggle preview of selected item
            if row := m.tbl.SelectedRow(); row != nil {
                it, _ := m.repo.GetItemByKey(row[2])
                m.preview.SetContent(it.Value)
                m.previewActive = !m.previewActive
                return m, nil
            }
        case "pgdn", "ctrl+f":
            m.pager.NextPage(); m.applyPagination(); return m, nil
        case "pgup", "ctrl+b":
            m.pager.PrevPage(); m.applyPagination(); return m, nil
        case "o": // Open Project (path input via prompt)
            m.mode = "openProject"
            m.keyInput.Placeholder = "Enter project folder path and press Enter"
            m.keyInput.SetValue("")
            m.keyInput.Focus()
            return m, nil
        case "/":
            m.search.Focus(); return m, nil
        case "tab": // cycle channel
            m.channelIx = (m.channelIx + 1) % len(channels)
            m.refreshRows(); return m, nil
        case "shift+tab": // cycle priority
            m.priorityIx = (m.priorityIx + 1) % len(priorities)
            m.refreshRows(); return m, nil
        case "r": // reload from repo
            m.toast = "Reloading..."
            return m, tea.Batch(clearToastAfter(1*time.Second), m.cmdReload())
        case "enter": // edit selected
            if row := m.tbl.SelectedRow(); row != nil {
                key := row[2]
                // fetch latest for value
                it, _ := m.repo.GetItemByKey(key)
                m.editExisting = true
                m.keyInput.SetValue(it.Key)
                m.keyInput.Blur() // disable editing key when updating
                m.valueInput.SetValue(it.Value)
                m.editChannelIx = indexOf(channels, it.Channel)
                if m.editChannelIx <= 0 { m.editChannelIx = 1 }
                m.editPriorityIx = indexOf(priorities, strings.ToLower(it.Priority))
                if m.editPriorityIx <= 0 { m.editPriorityIx = 1 }
                m.mode = "edit"
                m.valueInput.Focus()
                return m, nil
            }
        case "n": // new item
            m.editExisting = false
            m.keyInput.SetValue("")
            m.keyInput.Focus()
            m.valueInput.SetValue("")
            m.editChannelIx = 1
            m.editPriorityIx = 2 // default normal
            m.mode = "edit"
            return m, nil
        case "d": // delete selected
            if row := m.tbl.SelectedRow(); row != nil {
                m.mode = "confirmDelete"
                return m, nil
            }
        case "ctrl+c", "q":
            return m, tea.Quit
        }
    }
    var cmd tea.Cmd
    if m.mode == "openProject" {
        switch k := msg.(type) {
        case tea.KeyMsg:
            if k.Type == tea.KeyEnter {
                p := strings.TrimSpace(m.keyInput.Value())
                if p != "" {
                    if m.repo != nil { _ = m.repo.Close() }
                    return m, m.cmdOpenProject(p)
                }
                return m, nil
            } else if k.Type == tea.KeyEsc {
                m.mode = "list"
                return m, nil
            }
        }
        m.keyInput, cmd = m.keyInput.Update(msg)
        return m, cmd
    }
    // preview viewport update
    if m.previewActive {
        m.preview, _ = m.preview.Update(msg)
    }
    m.tbl, cmd = m.tbl.Update(msg)
    return m, cmd
}

var (
    colorPrimary = lipgloss.Color("212")
    colorSuccess = lipgloss.Color("42")
    colorError   = lipgloss.Color("196")
    colorMuted   = lipgloss.Color("240")

    titleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary)
    mutedStyle = lipgloss.NewStyle().Foreground(colorMuted)
    bold = lipgloss.NewStyle().Bold(true)
)

func (m model) View() string {
    if m.mode == "openProject" {
        header := titleStyle.Render("Open Project Folder")
        tip := mutedStyle.Render("Enter absolute path; Esc to cancel")
        return header + "\n" + tip + "\n" + m.keyInput.View() + "\n"
    }
    if m.initErr != nil {
        return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: ") + m.initErr.Error() + "\nPress q to quit."
    }
    header := titleStyle.Render("Context Keeper GUI")
    proj := mutedStyle.Render(fmt.Sprintf("Project: %s", m.projectPath))
    if m.mode == "edit" {
        ch := channels[max(1, m.editChannelIx)]
        pr := priorities[max(1, m.editPriorityIx)]
        tips := mutedStyle.Render("Tab switch field, [/] channel, -/= priority, s save, Esc cancel")
        keyView := m.keyInput.View()
        if m.editExisting { keyView = mutedStyle.Render("Key: ") + m.keyInput.Value() + mutedStyle.Render(" (immutable)") }
        form := keyView + "\n" + m.valueInput.View() + "\n" + fmt.Sprintf("Channel: %s  Priority: %s", ch, pr)
        return header + "\n" + proj + "\n\n" + form + "\n" + tips + viewHelpShort(m) + note(m.toast) + "\n"
    }
    if m.mode == "confirmDelete" {
        // render base view dimmed and overlay a centered modal
        base := header + "\n" + proj + "\n\n" + m.tbl.View() + viewHelpShort(m) + note(m.toast) + "\n"
        modal := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorMuted).Padding(1, 2).Render("Confirm delete? [y]=Yes [n/esc]=No")
        overlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
        dim := lipgloss.NewStyle().Foreground(colorMuted).Render(base)
        return dim + "\n" + overlay
    }
    filt := fmt.Sprintf("Channel: %s  Priority: %s  (/ focus search, Tab channel, Shift+Tab priority, enter edit, n new, d delete, r reload, p preview, PgUp/PgDn pages)",
        channels[m.channelIx], priorities[m.priorityIx])
    body := m.tbl.View()
    if m.previewActive {
        box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorMuted).Padding(0, 1)
        body = lipgloss.JoinHorizontal(lipgloss.Top, body, "  ", box.Render(m.preview.View()))
    }
    pagerView := mutedStyle.Render(m.pager.View())
    return header + "\n" + proj + "\n" + m.search.View() + "\n" + filt + "\n\n" + body + "\n" + pagerView + viewHelpShort(m) + note(m.toast) + "\n"
}

func note(s string) string {
    if strings.TrimSpace(s) == "" { return "" }
    if strings.HasPrefix(s, "Error:") {
        return "\n" + lipgloss.NewStyle().Foreground(colorError).Render(s)
    }
    if s == "Saved" || s == "Deleted" || s == "Project opened" {
        return "\n" + lipgloss.NewStyle().Foreground(colorSuccess).Render(s)
    }
    return "\n" + mutedStyle.Render(s)
}

func indexOf(arr []string, v string) int {
    for i, s := range arr { if s == v { return i } }
    return -1
}

func max(a, b int) int { if a > b { return a } ; return b }

func (m *model) applyPagination() {
    per := m.tblHeight
    if per <= 0 { per = 7 }
    m.pager.PerPage = per
    total := len(m.allRows)
    if total == 0 {
        m.tbl.SetRows(nil)
        m.pager.SetTotalPages(1)
        return
    }
    pages := total / per
    if total%per != 0 { pages++ }
    if pages <= 0 { pages = 1 }
    m.pager.SetTotalPages(pages)
    page := m.pager.Page
    if page >= pages { page = pages - 1; m.pager.Page = page }
    start := page * per
    end := start + per
    if end > total { end = total }
    m.tbl.SetRows(m.allRows[start:end])
}

// async commands
type reloadDoneMsg struct{ items []dbpkg.Item; err error }

func (m model) cmdReload() tea.Cmd {
    return func() tea.Msg {
        items, err := m.repo.ListItems(nil, nil)
        return reloadDoneMsg{items: items, err: err}
    }
}

type openDoneMsg struct{ items []dbpkg.Item; err error; project string }

func (m model) cmdOpenProject(p string) tea.Cmd {
    return func() tea.Msg {
        newDB := filepath.Join(p, ".context-keeper", "context.db")
        repo, err := dbpkg.NewRepository(newDB)
        if err != nil { return openDoneMsg{err: err} }
        items, e2 := repo.ListItems(nil, nil)
        if e2 != nil { _ = repo.Close(); return openDoneMsg{err: e2} }
        // swap repo in a follow-up Update to avoid race; pass project path
        // we will set m.repo on handling this msg
        _ = repo.Close() // open later for simplicity
        return openDoneMsg{items: items, err: nil, project: p}
    }
}

// key bindings and help
type keymap struct {
    Search, NextChan, NextPrio, Edit, New, Delete, OpenProj, Reload, Back, Quit key.Binding
}

func newKeymap() keymap {
    return keymap{
        Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
        NextChan: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next channel")),
        NextPrio: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("S-Tab", "next priority")),
        Edit:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit")),
        New:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
        Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
        OpenProj: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open project")),
        Reload:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
        Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
        Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
    }
}

func viewHelpShort(m model) string {
    // minimal inline help; non-blocking and uses tokens
    return "\n" + m.help.ShortHelpView([]key.Binding{m.keys.Search, m.keys.Edit, m.keys.New, m.keys.Delete, m.keys.Quit})
}

func main() {
    m := initialModel()
    p := tea.NewProgram(m)
    if _, err := p.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "fatal:", err)
        os.Exit(1)
    }
    if m.repo != nil {
        _ = m.repo.Close()
    }
}
