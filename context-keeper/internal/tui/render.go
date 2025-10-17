package tui

import (
    "fmt"

    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/lipgloss"
    dbpkg "context-keeper/internal/db"
)

var (
    bold = lipgloss.NewStyle().Bold(true)
)

// Design tokens (no hard-coded usage outside tokens)
var (
    colorPrimary = lipgloss.Color("212")
    colorSuccess = lipgloss.Color("42")
    colorError   = lipgloss.Color("196")
    colorMuted   = lipgloss.Color("240")

    keyStyle     = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary)
    successStyle = lipgloss.NewStyle().Foreground(colorSuccess)
    errorStyle   = lipgloss.NewStyle().Foreground(colorError)
    borderStyle  = lipgloss.NewStyle().Foreground(colorMuted)
)

// Public minimal helpers used by CLI rendering paths
func Success(msg string) string { return successStyle.Render("✅ " + msg) }
func Deleted(msg string) string { return successStyle.Render("🗑️ " + msg) }
func Highlight(s string) string { return keyStyle.Render(s) }

// TableStyles centralizes Bubbles table styles per tokens
func TableStyles() table.Styles {
    return table.Styles{
        Header:  lipgloss.NewStyle().Bold(true).Foreground(colorPrimary),
        Cell:    lipgloss.NewStyle(),
        Selected: lipgloss.NewStyle().Foreground(colorPrimary),
        Focused:  lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(colorMuted),
    }
}

func RenderList(items []dbpkg.Item) error {
    cols := []table.Column{{Title: bold.Render("CHANNEL"), Width: 12}, {Title: bold.Render("PRIORITY"), Width: 10}, {Title: bold.Render("KEY"), Width: 30}}
    var rows []table.Row
    for _, it := range items {
        rows = append(rows, table.Row{it.Channel, it.Priority, keyStyle.Render(it.Key)})
    }
    t := table.New(
        table.WithColumns(cols),
        table.WithRows(rows),
        table.WithStyles(TableStyles()),
    )
    view := t.View()
    fmt.Println(view)
    return nil
}
