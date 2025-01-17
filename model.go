package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"haproxy-runtime-cli/components"
	"haproxy-runtime-cli/styles"
	"net"
)

type sessionState uint

const (
	commandsListPage sessionState = iota
	executePage
	statusPage
)

type RuntimeAPI struct {
	page         sessionState
	commandsPage components.CommandsPage
	statusPage   components.StatusPage
	executePage  components.ExecutePage
}

type ActivePage interface {
	Supports(tea.Msg, bool) bool
}

func NewRuntimeApi(socket func() net.Conn) RuntimeAPI {
	return RuntimeAPI{
		page:         statusPage,
		commandsPage: components.NewCommandsPage(socket),
		statusPage:   components.NewStatusPage(socket),
		executePage:  components.NewExecutePage(socket),
	}
}

func (m RuntimeAPI) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle(fmt.Sprintf("haproxy-runtime-cli")),
		m.commandsPage.Init(),
		m.statusPage.Init(),
		m.executePage.Init(),
	)
}

func (m RuntimeAPI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case error:
		panic(msg)
	case components.ActivateStatusPage:
		m.page = statusPage
		return m, nil
	case components.ActivateCommandsPage:
		m.page = commandsListPage
		return m, nil
	case components.ActivateExecutePage:
		m.page = executePage
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// delegate to sub components
	if m.statusPage.Supports(msg, m.page == statusPage) {
		m.statusPage, cmd = m.statusPage.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.commandsPage.Supports(msg, m.page == commandsListPage) {
		m.commandsPage, cmd = m.commandsPage.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.executePage.Supports(msg, m.page == executePage) {
		m.executePage, cmd = m.executePage.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RuntimeAPI) View() string {
	s := header() + "\n"

	switch m.page {
	case statusPage:
		s += m.statusPage.View()
	case commandsListPage:
		s += m.commandsPage.View()
	case executePage:
		s += m.executePage.View()
	}

	return styles.PageStyle.Render(s)
}

func header() string {
	return styles.HeaderStyle.Render(styles.AppName)
}
