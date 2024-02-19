// type status int

// const (
// 	todo status = iota
// 	inProgress
// 	done
// )

// type Task struct {
// 	title string
// 	description string
// 	Status status
// }

// func (t *Task) FilterValue() string {
// 	return t.title
// }

// func (t *Task) Title() string {
// 	return t.title
// }

// func (t *Task) Description() string {
// 	return t.description
// }

// type Model struct {
// 	list list.Model
// 	err error
// }

// func (m *Model) initList() {
// 	m.list = list.New([]list.Item{}, list.NewDefaultDelegate())
// }

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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

// type item struct {
// 	title     string
// 	desc      string
// 	extra     string
// 	startDate time.Time
// 	dueDate   time.Time
// }

// func (i item) Title() string       { return i.title }
// func (i item) Description() string { return i.desc }
// func (i item) Extra() string       { return i.extra }
// func (i item) StartDate() string   { return i.startDate.String() }
// func (i item) DueDate() string     { return i.dueDate.String() }
// func (i item) FilterValue() string { return i.title }

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

func CreateItem(title, desc string) item {

	now := time.Now()
	yyyy, mm, dd := now.Date()

	return item{
		title:       title,
		description: desc,
		status:      todo,
		startDate:   time.Now(),
		dueDate:     time.Date(yyyy, mm, dd+1, 8, 0, 0, 0, now.Location()),
	}
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

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return m, tea.Quit
// 		}
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)
// 	}

// 	var cmd tea.Cmd
// 	m.list, cmd = m.list.Update(msg)
// 	return m, cmd
// }

// func (m model) View() string {
// 	return docStyle.Render(m.list.View())
// }

func main() {
	items := []list.Item{
		// item{title: "Raspberry Pi’s", extra: "I have ’em all over my house!", startDate: time.Now()},
		// item{title: "Nutella", desc: "It's good on toast"},
		// item{title: "Bitter melon", desc: "It cools you down"},
		// item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		// item{title: "Eight hours of sleep", desc: "I had this once"},
		// item{title: "Cats", desc: "Usually"},
		// item{title: "Plantasia, the album", desc: "My plants love it too"},
		// item{title: "Pour over coffee", desc: "It takes forever to make though"},
		// item{title: "VR", desc: "Virtual reality...what is there to say?"},
		// item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		// item{title: "Linux", desc: "Pretty much the best OS"},
		// item{title: "Business school", desc: "Just kidding"},
		// item{title: "Pottery", desc: "Wet clay is a great feeling"},
		// item{title: "Shampoo", desc: "Nothing like clean hair"},
		// item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		// item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		// item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		// item{title: "Stickers", desc: "The thicker the vinyl the better"},
		// item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		// item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		// item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		// item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		// item{title: "Terrycloth", desc: "In other words, towel fabric"},
		item{title: "Hamburguer", description: "4 guys"},
		item{title: "Ramen", description: "Wok Restaurant"},
		item{title: "Fries", description: "McDonald's"},
		item{title: "Churros", description: "La churrería"},
	}

	// m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	// m.list.Title = "My Fave Things"

	// p := tea.NewProgram(m, tea.WithAltScreen())

	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
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
