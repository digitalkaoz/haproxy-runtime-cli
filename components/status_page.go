package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"haproxy-runtime-cli/styles"
	"net"
	"strconv"
)

var tblStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type StatusPage struct {
	socket   func() net.Conn
	keys     statusPageKeyMap
	backends []haproxy.Backend
	table    table.Model
}

type statusPageKeyMap struct {
	GotoCommands key.Binding
	Reload       key.Binding
	Quit         key.Binding
}

type ActivateStatusPage bool

func NewStatusPage(socket func() net.Conn) StatusPage {
	km := createStatusKeyMap()
	return StatusPage{
		socket: socket,
		keys:   km,
		table:  createTable(km),
	}
}

func (s StatusPage) Init() tea.Cmd {
	return fetchBackends(s.socket)
}

func (s StatusPage) Update(msg tea.Msg) (StatusPage, tea.Cmd) {
	switch msg := msg.(type) {
	case []haproxy.Backend:
		s.backends = msg
		s.table.SetRows(backendsToRows(s.backends))
		s.table = recalculateTableSize(s.table)
	case tea.WindowSizeMsg:
		s.table.UpdateViewport()
		s.table.SetWidth(msg.Width - styles.PageStyle.GetHorizontalMargins())
		s.table.SetHeight(msg.Height - styles.PageStyle.GetVerticalMargins() - 3 - 3)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.keys.GotoCommands):
			return s, ActivateCommandsPageCmd()
		case key.Matches(msg, s.keys.Reload):
			return s, fetchBackends(s.socket)
		}
	}

	return updateTeaTable(s, msg)
}

func (s StatusPage) View() string {
	return tblStyle.Render(s.table.View()) + "\n" + s.table.HelpView()
}

func (s StatusPage) Supports(msg tea.Msg, isActive bool) bool {
	switch msg.(type) {
	case []haproxy.Backend, tea.WindowSizeMsg:
		return true
	case tea.KeyMsg:
		if isActive {
			return true
		}
	}

	return false
}

func ActivateStatusPageCmd() tea.Cmd {
	return func() tea.Msg {
		return ActivateStatusPage(true)
	}
}

func createStatusKeyMap() statusPageKeyMap {
	return statusPageKeyMap{
		GotoCommands: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "command list"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

func updateTeaTable(s StatusPage, msg tea.Msg) (StatusPage, tea.Cmd) {
	var cmd tea.Cmd
	s.table, cmd = s.table.Update(msg)

	return s, cmd
}

func recalculateTableSize(model table.Model) table.Model {
	sizes := make([]int, len(model.Columns()))

	for _, r := range model.Rows() {
		for column, cell := range r {
			if sizes[column] < len(cell) {
				sizes[column] = len(cell)
			}
		}
	}

	// set new column widths
	cols := model.Columns()
	for i := range cols {
		cols[i].Width = sizes[i] + 2
		if cols[i].Width < len(cols[i].Title) {
			cols[i].Width = len(cols[i].Title) + 2
		}
	}
	model.SetColumns(cols)

	return model
}

func createTable(km statusPageKeyMap) table.Model {
	s := table.DefaultStyles()
	//selected := s.Selected
	s.Selected = styles.ActiveStyle

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	//TODO responsive table, calculate max width based on values
	tbl := table.New(
		table.WithHeight(10),
		table.WithColumns([]table.Column{
			{Title: "Backend", Width: 25},
			{Title: "Name", Width: 25},
			{Title: "Weight", Width: 6},
			{Title: "State", Width: 8},
			{Title: "IP", Width: 16},
			{Title: "CHECK", Width: 8},
			{Title: "", Width: 8},
			{Title: "FQDN", Width: 30},
			{Title: "SSL", Width: 5},
		}),
		table.WithFocused(true),
		table.WithStyles(s),
		table.WithAdditionalShortHelpKeys([]key.Binding{km.GotoCommands, km.Reload, km.Quit}),
	)

	return tbl
}

func backendsToRows(backends []haproxy.Backend) []table.Row {
	var rows []table.Row

	for _, b := range backends {
		rows = append(rows, table.Row{
			b.Name,
		})
		for _, s := range b.Servers {
			addr := ""
			if s.Address != nil {
				addr = s.Address.String()
			}

			rows = append(rows, table.Row{
				"",
				s.Name,
				strconv.Itoa(s.CalculatedWeight),
				s.State,
				addr,
				s.CheckState,
				s.SrvCheckResult,
				fmt.Sprintf(`%s:%d`, s.Fqdn, s.Port),
				strconv.FormatBool(s.UseSSL)})
		}
	}

	return rows
}

func fetchBackends(s func() net.Conn) tea.Cmd {
	return socket.ExecCmd[[]haproxy.Backend](
		s,
		"show servers state",
		func(s *string) []haproxy.Backend { return haproxy.ParseBackends(*s) },
	)
}
