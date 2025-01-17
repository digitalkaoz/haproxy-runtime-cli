package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"haproxy-runtime-cli/styles"
	"net"
	"strings"
)

type ExecutePage struct {
	command  haproxy.Command
	socket   func() net.Conn
	keys     executePageKeyMap
	help     help.Model
	input    textinput.Model
	response ExecuteResponse
}

type executePageKeyMap struct {
	GotoCommands key.Binding
	Execute      key.Binding
}

type ExecuteResponse string

type ActivateExecutePage bool

func NewExecutePage(socket func() net.Conn) ExecutePage {
	ti := createInput()

	return ExecutePage{
		socket: socket,
		keys:   createExecuteKeyMap(),
		help:   help.New(),
		input:  ti,
	}
}

func (e ExecutePage) Init() tea.Cmd {
	return nil
}

func (e ExecutePage) Update(msg tea.Msg) (ExecutePage, tea.Cmd) {
	switch msg := msg.(type) {
	case ExecuteResponse:
		e.response = msg
	case haproxy.Command:
		e.command = msg
		e.input.Prompt = e.command.Name + " "
		e.input.Placeholder = e.command.Args
		e.input.SetValue("")
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keys.GotoCommands):
			if e.input.Value() == "" {
				return e, ActivateCommandsPageCmd()
			}
		case key.Matches(msg, e.keys.Execute):
			//TODO properly execute command and show result
			return e, executeCommand(e)
		}
	}

	var cmd tea.Cmd
	e.input, cmd = e.input.Update(msg)

	return e, cmd
}

func executeCommand(e ExecutePage) func() tea.Msg {
	return socket.ExecCmd[ExecuteResponse](
		e.socket,
		fmt.Sprintf("%s %s", e.command.Name, e.input.Value()),
		func(s *string) ExecuteResponse { return ExecuteResponse(*s) },
	)
}

func (e ExecutePage) View() string {
	return description(e.command) + "\n" +
		input(e.input) + "\n" +
		response(e.response) + "\n" +
		e.help.ShortHelpView([]key.Binding{e.keys.GotoCommands, e.keys.Execute})
}

func (e ExecutePage) Supports(msg tea.Msg, isActive bool) bool {
	switch msg.(type) {
	case ExecuteResponse, haproxy.Command:
		return true
	case tea.KeyMsg:
		if isActive {
			return true
		}
	}

	return false
}

func ActivateExecutePageCmd() tea.Cmd {
	return func() tea.Msg {
		return ActivateExecutePage(true)
	}
}

func createExecuteKeyMap() executePageKeyMap {
	return executePageKeyMap{
		GotoCommands: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "command list"),
		),
		Execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "execute"),
		),
	}
}

func input(t textinput.Model) string {
	total := strings.Count(t.Value(), " ")
	curr := strings.Split(t.Placeholder, " ")
	if total > len(curr) {
		total = len(curr)
	}

	return t.View() + t.PlaceholderStyle.Render(strings.Join(curr[total:], " "))
}

func description(c haproxy.Command) string {
	return styles.ActiveStyle.MarginTop(1).Render(c.Help)
}

func response(r ExecuteResponse) string {
	return styles.ResponseStyle.Render(string(r))
}

func createInput() textinput.Model {
	ti := textinput.New()
	//ti.Width = 80
	ti.PromptStyle = styles.ComplementStyle.MarginTop(1)
	ti.TextStyle = styles.ComplementStyle.Bold(false)

	ti.Focus()

	return ti
}
