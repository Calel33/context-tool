package tui

import (
    "fmt"

    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/lipgloss"
    dbpkg "context-keeper/internal/db"
)

var (
    colorPrimary = lipgloss.Color("212")
    colorSuccess = lipgloss.Color("42")
    colorError   = lipgloss.Color("196")
    colorMuted   = lipgloss.Color("240")

    bold = lipgloss.NewStyle().Bold(true)
    keyStyle = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary)
    successStyle = lipgloss.NewStyle().Foreground(colorSuccess)
    errorStyle = lipgloss.NewStyle().Foreground(colorError)
    borderStyle = lipgloss.NewStyle().Foreground(colorMuted)
)

func Success(msg string) string { return successStyle.Render("✅ " + msg) }
func Deleted(msg string) string { return successStyle.Render("🗑️ " + msg) }
func Highlight(s string) string { return keyStyle.Render(s) }

func RenderList(items []dbpkg.Item) error {
    cols := []table.Column{{Title: bold.Render("CHANNEL"), Width: 12}, {Title: bold.Render("PRIORITY"), Width: 10}, {Title: bold.Render("KEY"), Width: 30}}
    var rows []table.Row
    for _, it := range items {
        rows = append(rows, table.Row{it.Channel, it.Priority, keyStyle.Render(it.Key)})
    }
    t := table.New(table.WithColumns(cols), table.WithRows(rows))
    view := t.View()
    fmt.Println(view)
    return nil
}
