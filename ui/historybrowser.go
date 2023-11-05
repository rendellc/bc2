package ui

import (
	"fmt"
	"io"
	"log"
	"rendellc/bc2/storage"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 20
const defaultWidth = 20

var blurredSelectedItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
var focusSelectedItemStyle = blurredSelectedItemStyle.Copy().Foreground(lipgloss.Color("176"))

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = blurredSelectedItemStyle
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item storage.ScriptInfo

func (i item) Title() string       { return storage.ScriptInfo(i).Name() }
func (i item) FilterValue() string { return storage.ScriptInfo(i).Name() }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {

			return selectedItemStyle.Render("> " + strings.Join(s, " "))
			// return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprintf(w, fn(str))
}

type historySelectMsg item

func newHistorySelectCmd(item item) tea.Cmd {
	return func() tea.Msg {
		return historySelectMsg(item)
	}
}

type historyBrowser struct {
	storage  *storage.Storage
	list     list.Model
	hasFocus bool
}

func CreateHistoryBrowser(storage *storage.Storage) historyBrowser {
	scriptInfos := storage.GetScriptInfos()
	items := []list.Item{}
	for _, scriptInfo := range scriptInfos {
		items = append(items, item(scriptInfo))
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "History"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.KeyMap.Quit.Unbind()

	return historyBrowser{
		storage:  storage,
		list:     l,
		hasFocus: false,
	}
}

func (h historyBrowser) Init() tea.Cmd {
	log.Printf("Initialize historybrowser %v", h.storage)
	return nil
}

func (h historyBrowser) Update(msg tea.Msg) (historyBrowser, tea.Cmd) {
	switch msg := msg.(type) {
	case focusChangeMsg:
		h.hasFocus = (focusElement(msg) == focusHistory)

		if h.hasFocus {
			selectedItemStyle = focusSelectedItemStyle
		} else {
			selectedItemStyle = blurredSelectedItemStyle
		}

	case tea.WindowSizeMsg:
		h.list.SetWidth(msg.Width / 2)
		return h, nil
	case tea.KeyMsg:
		if !h.hasFocus {
			return h, nil
		}
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := h.list.SelectedItem().(item)
			if !ok {
				return h, nil
			}

			log.Printf("Selected: %s", i)
			return h, newHistorySelectCmd(i)
		}
	}

	var cmd tea.Cmd
	h.list, cmd = h.list.Update(msg)

	return h, cmd
}

func (h historyBrowser) View() string {
	return h.list.View()
}
