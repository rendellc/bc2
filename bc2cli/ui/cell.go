package ui

import (
	"fmt"
	"rendellc/bc2/calc"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func color(s string) lipgloss.Color {
	return lipgloss.Color(s)
}

var cellStyle = lipgloss.NewStyle()
var cellInputStyle = lipgloss.NewStyle()
var cellResultStyle = lipgloss.NewStyle()

type cell struct {
	input  textinput.Model
	result *calc.InterpreterLineResult
}

func createCell() cell {
	t := textinput.New()
	t.Placeholder = ""
	t.Prompt = ""

	t.Blur()
	return cell{
		input:  t,
		result: nil,
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

func (c *cell) SetResult(result calc.InterpreterLineResult) {
	if c.result == nil {
		c.result = new(calc.InterpreterLineResult)
	}

	*c.result = result
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

func (c cell) Length() int {
	return len(c.input.Value())
}

func (c cell) Update(msg tea.Msg) (cell, tea.Cmd) {
	cellInput, cmd := c.input.Update(msg)
	c.input = cellInput
	return c, cmd
}

func (c cell) View(cellWidth int) string {
	hasContent := len(c.input.Value()) > 0
	inputView := c.input.View()

	resultView := ""
	if c.result != nil {
		isError := c.result.ResultType() == calc.LineResultError
		if hasContent && !isError {
			resultView = "= " + c.result.Message()
		} 
	}

	return cellStyle.Render(fmt.Sprintf("%s\t\t\t\t%s",
		cellInputStyle.Render(inputView),
		cellResultStyle.Render(resultView)))
}
