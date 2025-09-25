package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TUIState int

type sigMsg syscall.Signal

const (
	RunningState TUIState = iota
	WaitingState
	EditingState
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
	state         TUIState
	quitting      bool
}

var (
	modelStyle = lipgloss.NewStyle().
			Width(30).
			Height(7).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	runningStyle = lipgloss.NewStyle().
			Width(30).
			Height(7).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder()).Background(lipgloss.Color("#5b2599"))
	historyModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(7).
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
	return nil
}

func GetHistory() string {
	today, thisWeek, lastWeek, err := LoadHistort(Filename)
	if err != nil {
		fmt.Println(err)
		return historyModelStyle.Render("Today: --\nThis Week: --\nLast Week: --")
	}
	s := "Today: \n" + today + "\n"
	s += "This Week: \n" + thisWeek + "\n"
	s += "Last Week: \n" + lastWeek
	return historyModelStyle.Render(s)
}

func (m model) View() string {

	main_view := "\n"
	s := ""
	if m.state == RunningState {
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
		main_view += m.laststartTime.Format(layout)
		main_view += "\nTo: "
		main_view += m.lastendTime.Format(layout)
	}

	var style lipgloss.Style
	if m.state == RunningState {
		style = runningStyle
	} else {
		style = modelStyle
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top,
		m.history,
		style.Render(main_view))

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
			if m.state == RunningState {
				AddEndToCSV(Filename, time.Now())
			}
			return m, tea.Quit
		case key.Matches(msg, m.keymap.start):
			m.keymap.stop.SetEnabled(true)
			m.keymap.start.SetEnabled(false)
			m.startTime = time.Now()
			AddStartToCSV(Filename, m.startTime)
			m.state = RunningState
			return m, m.spinner.Tick
		case key.Matches(msg, m.keymap.tab):
			m.tab = 1 - m.tab
			return m, cmd
		case key.Matches(msg, m.keymap.stop):
			m.keymap.stop.SetEnabled(false)
			m.keymap.start.SetEnabled(true)
			m.laststartTime = m.startTime
			m.lastendTime = time.Now()
			AddEndToCSV(Filename, m.lastendTime)
			m.state = WaitingState
			m.history = GetHistory()
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case sigMsg:
		if m.state == RunningState {

			if m.state == RunningState {
				AddEndToCSV(Filename, time.Now())
			}
		}
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.state == RunningState {
				AddEndToCSV(Filename, time.Now())
			}
			return m, tea.Quit
		}
	}

	if m.state == RunningState {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

var Filename string

const layout = "15:04:05 2006/01/02"

func main() {
	// parse CLI flags
	f := flag.String("f", "data.csv", "CSV filename to read/write")
	flag.Parse()
	Filename = *f
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
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	m.state = WaitingState
	m.history = GetHistory()
	m.keymap.tab.SetEnabled(false)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	p := tea.NewProgram(m, tea.WithoutSignalHandler(), tea.WithAltScreen())
	go func() {
		sig := <-sigs
		p.Send(sigMsg(sig.(syscall.Signal)))
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}

}
