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

type historyItem storage.StoreMeta

func (i historyItem) Title() string       { return storage.StoreMeta(i).Name() }
func (i historyItem) FilterValue() string { return storage.StoreMeta(i).Name() }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(historyItem)
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

type historySelectMsg historyItem

func newHistorySelectCmd(item historyItem) tea.Cmd {
	return func() tea.Msg {
		return historySelectMsg(item)
	}
}

type historyBrowser struct {
	storage  storage.Store
	list     *list.Model
	hasFocus bool
}

func (h historyBrowser) refreshScriptBrowserCmd() tea.Cmd {
	scriptInfos, err := h.storage.GetScriptInfos()
	if err != nil {
		log.Printf("refresh script info failed: %s", err.Error())
		return nil
	}

	items := make([]list.Item,len(scriptInfos))
	for i, scriptInfo := range scriptInfos {
		items[i] = historyItem(scriptInfo)
	}

	log.Printf("Refreshing history browser with %d items", len(scriptInfos))

	return h.list.SetItems(items)
}

func CreateHistoryBrowser(storage storage.Store) historyBrowser {
	items := []list.Item{}
	listModel := new(list.Model)
	*listModel = list.New(items, itemDelegate{}, defaultWidth, listHeight)
	listModel.SetShowTitle(false)
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)
	listModel.Styles.Title = titleStyle
	listModel.Styles.PaginationStyle = paginationStyle
	listModel.Styles.HelpStyle = helpStyle
	listModel.KeyMap.Quit.Unbind()

	return historyBrowser{
		storage:  storage,
		list:     listModel,
		hasFocus: false,
	}
}

func (h historyBrowser) Init() tea.Cmd {
	log.Printf("Initialize historybrowser %v", h.storage)
	return h.refreshScriptBrowserCmd()
}

func (h historyBrowser) Update(msg tea.Msg) (historyBrowser, tea.Cmd) {
	switch msg := msg.(type) {
	case fileSavedMsg:
		return h, h.refreshScriptBrowserCmd()
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
			i, ok := h.list.SelectedItem().(historyItem)
			if !ok {
				return h, nil
			}

			log.Printf("Selected: %s", i)
			return h, newHistorySelectCmd(i)
		}
	}

	var cmd tea.Cmd
	*h.list, cmd = h.list.Update(msg)

	return h, cmd
}

func (h historyBrowser) View() string {
	return "History\n\n" + h.list.View()
}
