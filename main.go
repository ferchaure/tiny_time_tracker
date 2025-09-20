package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	keymap        keymap
	help          help.Model
	spinner       spinner.Model
	startTime     time.Time
	laststartTime time.Time
	lastendTime   time.Time
	history       string
	tab           uint
	running       bool
	quitting      bool
}

var (
	modelStyle = lipgloss.NewStyle().
			Width(30).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	historyModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Left, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type keymap struct {
	start key.Binding
	stop  key.Binding
	e     key.Binding
	quit  key.Binding
	tab   key.Binding
}

func (m model) Init() tea.Cmd {
	m.running = false
	m.history = GetHistory()
	return nil
}

func GetHistory() string {
	s := "Today: \n"
	s += "This Week:\n"
	//s +=  compute this week
	s += "Last Week:\n"
	//s +=  compute this week
	return historyModelStyle.Render(s)
}

func (m model) View() string {

	main_view := ""
	s := ""
	if m.running {
		main_view += m.spinner.View()
		seconds := int(time.Since(m.startTime).Seconds())
		h := seconds / 3600
		minutes := (seconds % 3600) / 60
		main_view += fmt.Sprintf("\n%02d:%02d:%02d", h, minutes, seconds%60)
		main_view += "\n\n--Current interval-- \nFrom: "
		main_view += m.startTime.Format("15:04:05 2006/01/02")

	} else {
		main_view += "▶"
		//current timer:
		main_view += "\n\n--Last interval-- \nFrom: "
		main_view += m.laststartTime.Format("15:04:05 2006/01/02")
		main_view += "\nTo: "
		main_view += m.lastendTime.Format("15:04:05 2006/01/02")
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top,
		m.history,
		modelStyle.Render(main_view))

	s += helpStyle.Render("\n" + m.helpView())
	return s
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.e,
		m.keymap.tab,
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.start):
			m.keymap.stop.SetEnabled(!m.running)
			m.keymap.start.SetEnabled(m.running)
			m.startTime = time.Now()
			m.running = true
			return m, m.spinner.Tick
		case key.Matches(msg, m.keymap.tab):
			m.tab = 1 - m.tab
			return m, cmd
		case key.Matches(msg, m.keymap.stop):
			m.keymap.stop.SetEnabled(!m.running)
			m.keymap.start.SetEnabled(m.running)
			m.laststartTime = m.startTime
			m.lastendTime = time.Now()
			m.running = false
			return m, nil
		}

	}
	if m.running {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
func main() {
	m := model{
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "stop"),
			),
			e: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit last interval"),
			),
			tab: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "change tab"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#2E2AEB"))
	m.running = false
	m.history = GetHistory()
	m.keymap.tab.SetEnabled(false)
	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion()).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}

}
