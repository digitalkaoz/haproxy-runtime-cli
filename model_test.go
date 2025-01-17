package main

import (
	"github.com/stretchr/testify/assert"
	"haproxy-runtime-cli/components"
	"haproxy-runtime-cli/haproxy"
	"net"
	"testing"
)

func TestNewModel(t *testing.T) {
	m := NewRuntimeApi(func() net.Conn {
		return nil
	})

	assert.NotNil(t, m)
	assert.NotNil(t, m.statusPage)
	assert.NotNil(t, m.commandsPage)
	assert.NotNil(t, m.executePage)
	assert.Equal(t, sessionState(2), m.page)
}

func TestInit(t *testing.T) {
	cmd := NewRuntimeApi(func() net.Conn { return nil }).Init()

	assert.NotNil(t, cmd)
	cmds := cmd()

	assert.Len(t, cmds, 3) //tea.Cmd from sub component inits
}

func TestViewStatus(t *testing.T) {
	m := NewRuntimeApi(func() net.Conn { return nil })
	res := m.View()

	// default page is status page
	assert.Contains(t, res, "haproxy-runtime-cli")
	assert.Contains(t, res, "Backend") // a column from status page

	nm, _ := m.Update(components.ActivateStatusPage(true))
	res = nm.View()

	assert.Contains(t, res, "haproxy-runtime-cli")
	assert.Contains(t, res, "Backend")

}

func TestViewCommands(t *testing.T) {
	m := NewRuntimeApi(func() net.Conn { return nil })
	nm, _ := m.Update(components.ActivateCommandsPage(true))
	res := nm.View()

	assert.Contains(t, res, "haproxy-runtime-cli")
	assert.Contains(t, res, "No items") // a column from status page
}

func TestViewExecute(t *testing.T) {
	m := NewRuntimeApi(func() net.Conn { return nil })
	nm, _ := m.Update(components.ActivateExecutePage(true))
	res := nm.View()

	assert.Contains(t, res, "haproxy-runtime-cli")
	assert.Contains(t, res, "enter execute") // a column from status page
}

func TestUpdateWithKnownCommands(t *testing.T) {
	m := NewRuntimeApi(func() net.Conn { return nil })

	// supported by status
	_, cmds := m.Update([]haproxy.Backend{
		{Id: 1, Name: "backend1"},
	})

	assert.Nil(t, cmds)

	// supported by commands
	nm, _ := m.Update(components.ActivateCommandsPage(true))
	_, cmds = nm.Update(haproxy.ParsedHelp{
		{Name: "foo", Help: "bar"},
	})

	assert.Nil(t, cmds)

	// supported by execute
	nm, _ = m.Update(components.ActivateExecutePage(true))
	_, cmds = nm.Update(haproxy.Command{
		Name: "foo", Help: "bar",
	})

	assert.Nil(t, cmds)
}
