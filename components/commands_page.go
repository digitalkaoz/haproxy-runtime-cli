package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"haproxy-runtime-cli/styles"
	"net"
)

type CommandsPage struct {
	commands list.Model
	socket   func() net.Conn
	keys     commandsPageKeyMap
}

type commandsPageKeyMap struct {
	GotoStatusPage  key.Binding
	GotoExecutePage key.Binding
}

type ActivateCommandsPage bool

func NewCommandsPage(socket func() net.Conn) CommandsPage {
	keys := createCommandsKeyMap()
	return CommandsPage{
		socket:   socket,
		commands: createList(keys),
		keys:     keys,
	}
}

func (c CommandsPage) Init() tea.Cmd {
	if len(c.commands.Items()) > 0 {
		return nil
	}

	return tea.Batch(
		c.commands.StartSpinner(),
		fetchHelpCommand(c.socket),
	)
}

func (c CommandsPage) Update(msg tea.Msg) (CommandsPage, tea.Cmd) {
	switch msg := msg.(type) {

	case haproxy.ParsedHelp:
		c.commands.StopSpinner()
		c.commands.SetItems(helpToList(msg))
		return c, nil
	case tea.WindowSizeMsg:
		resizeList(&c.commands, msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.GotoStatusPage):
			return c, ActivateStatusPageCmd()
		case key.Matches(msg, c.keys.GotoExecutePage):
			return c, tea.Sequence(ActivateExecutePageCmd(), func() tea.Msg {
				return c.commands.SelectedItem()
			})
		}
	}

	return updateTeaList(c, msg)
}

func (c CommandsPage) View() string {
	return c.commands.View()
}

func (c CommandsPage) Supports(msg tea.Msg, isActive bool) bool {
	switch msg.(type) {
	case haproxy.ParsedHelp, tea.WindowSizeMsg, list.FilterMatchesMsg:
		return true
	case tea.KeyMsg:
		if isActive {
			return true
		}
	}

	return false
}

func ActivateCommandsPageCmd() tea.Cmd {
	return func() tea.Msg {
		return ActivateCommandsPage(true)
	}
}

func errorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}

func createCommandsKeyMap() commandsPageKeyMap {
	return commandsPageKeyMap{
		GotoStatusPage: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "goto status page"),
		),
		GotoExecutePage: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "goto command execute page"),
		),
	}
}

func updateTeaList(c CommandsPage, msg tea.Msg) (CommandsPage, tea.Cmd) {
	var cmd tea.Cmd
	c.commands, cmd = c.commands.Update(msg)

	return c, cmd
}

func createList(keys commandsPageKeyMap) list.Model {
	l := list.New(nil, createListDelegate(), 160, 40)
	l.FilterInput.ShowSuggestions = true
	l.Styles.ActivePaginationDot = l.Styles.ActivePaginationDot.Foreground(styles.ActiveColor).Bold(true)
	l.Styles.PaginationStyle = l.Styles.PaginationStyle.PaddingLeft(0)
	l.Paginator.ActiveDot = l.Styles.ActivePaginationDot.String()
	l.Styles.HelpStyle = l.Styles.HelpStyle.PaddingLeft(0)

	l.StartSpinner()
	l.SetShowTitle(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.GotoStatusPage,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.GotoStatusPage,
		}
	}

	return l
}

func resizeList(model *list.Model, msg tea.WindowSizeMsg) {
	_, _ = styles.PageStyle.GetFrameSize()
	_, hv := styles.HeaderStyle.GetFrameSize()
	if msg.Width == 0 || msg.Height == 0 {
		return
	}
	model.SetSize(msg.Width, msg.Height-(hv*2))
}

func fetchHelpCommand(sock func() net.Conn) tea.Cmd {
	return func() tea.Msg {
		out, err := socket.Exec(sock, "help")
		if err != nil {
			return errorCmd(err)
		}
		return haproxy.ParseHelp(out)
	}
}

func helpToList(help haproxy.ParsedHelp) []list.Item {
	var cmds = make([]list.Item, len(help))
	for i, r := range help {
		cmds[i] = r
	}

	return cmds
}

func createListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = styles.ActiveStyle.PaddingLeft(1).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(styles.ActiveColor)
	d.Styles.SelectedDesc = styles.ComplementStyle.PaddingLeft(1).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(styles.ActiveColor)

	return d
}
