package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	// "time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var docStyle = lipgloss.NewStyle().Margin(1, 2)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Background((lipgloss.Color("#7D56F4")))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type status int

const (
	todo status = iota
	inProgress
	completed
)


type item struct {
	title       string
	description string
	status      status
	startDate   time.Time
	dueDate     time.Time
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) DueDate() string     { return i.dueDate.Format("02-01-2006") }
func (i item) Status() string {
	switch i.status {
	case 0:
		return "To do"
	case 1:
		return "In progress"
	case 2:
		return "Completed"
	}
	return ""
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
		if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)
	strRest := fmt.Sprintf("\t  %s", i.description)

	fn := itemStyle.Render
	fn2 := itemStyle.Render
	if index == m.Index() {
	fn = func(s ...string) string {
	return selectedItemStyle.Render("> " + strings.Join(s, " "))
	}
		fn2 = func(s ...string) string {
	return selectedItemStyle.Render(" " + strings.Join(s, " "))
	}
	}

	fmt.Fprintln(w, fn(str))
fmt.Fprint(w, fn2(strRest))
}

type model struct {
	list     list.Model
	choice   item
	quitting bool
	err      error
}

func New() *model {
	return &model{err: nil}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return &m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return &m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice.title = i.title
				m.choice.description = i.description
			}
			return &m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return &m, cmd
}

func (m model) View() string {
	if m.choice.title != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.\n Enjoy %s!!", m.choice.title, m.choice.description))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}



type ListCfg struct {
	width     int
	height    int
	itemList  []list.Item
	listTitle string
}

func (m *model) initList(listCfg ListCfg) {
	m.list = list.New(listCfg.itemList, itemDelegate{}, listCfg.width, listCfg.height)
	m.list.Title = listCfg.listTitle
}

func main() {

	itemsInit := []list.Item{
		item{title: "Hamburguer", description: "4 guys"},
		item{title: "Ramen", description: "Wok Restaurant"},
		item{title: "Fries", description: "McDonald's"},
		item{title: "Churros", description: "La churrería"},
	}

	// listConfig := ListCfg{
	// 	width:     300,
	// 	height:    200,
	// 	itemList:  itemsInit,
	// 	listTitle: "Restaurant List",
	// }

	const defaultWidth = 20

	// m := New()
	// m.initList(listConfig)
	// m.list.SetShowStatusBar(false)
	// m.list.SetFilteringEnabled(false)
	// m.list.Styles.Title = titleStyle
	// m.list.Styles.PaginationStyle = paginationStyle
	// m.list.Styles.HelpStyle = helpStyle
	// p := tea.NewProgram(m)

	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
	l := list.New(itemsInit, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
