package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func color(s string) lipgloss.Color {
	return lipgloss.Color(s)
}

var cellStyle = lipgloss.NewStyle().Padding(2)
var cellInputStyle = lipgloss.NewStyle()
var cellResultStyle = lipgloss.NewStyle().Foreground(color("51"))

type cell struct {
	input  textinput.Model
	result string
}

func createCell() cell {
	t := textinput.New()
	t.Placeholder = ""
	t.Prompt = ""

	t.Blur()
	return cell{
		input:  t,
		result: "",
	}
}

func (c *cell) SetValue(value string) {
	c.input.SetValue(value)
}

func (c cell) Position() int {
	return c.input.Position()
}

func (c *cell) SetCursor(pos int) {
	c.input.SetCursor(pos)
}

func (c *cell) SetResult(result string) {
	c.result = result
}

func (c cell) Value() string {
	return c.input.Value()
}

func (c *cell) Focus() {
	c.input.Focus()
}

func (c *cell) Blur() {
	c.input.Blur()
}

func (c cell) Init() tea.Cmd {
	return nil
}

func (c cell) Update(msg tea.Msg) (cell, tea.Cmd) {
	cellInput, cmd := c.input.Update(msg)
	c.input = cellInput
	return c, cmd
}

func (c cell) View(cellWidth int) string {
	inputStyleRender := cellInputStyle.Width(cellWidth).Render
	return cellInputStyle.Render(fmt.Sprintf("%s\t\t%s",
		inputStyleRender(c.input.View()),
		cellResultStyle.Render(c.result)))
}
