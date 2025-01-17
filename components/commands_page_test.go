package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"net"
	"reflect"
	"testing"
)

func TestCommandsPage(t *testing.T) {
	t.Parallel()

	model := func() CommandsPage {
		l := list.Model{}
		l.SetDelegate(createListDelegate())
		l.SetItems(make([]list.Item, 1))
		return CommandsPage{commands: l}
	}

	socketModel := func() CommandsPage {
		conn := &socket.DummySocket{
			Output: []byte("The following commands are valid at this level\nhelp : foo"),
		}

		return NewCommandsPage(func() net.Conn { return conn })
	}

	parsedModel := func() CommandsPage {
		m, _ := socketModel().Update(haproxy.ParsedHelp{haproxy.Command{Name: "help", Help: "foo"}})
		return m
	}

	t.Run("New", func(t *testing.T) {
		m := NewCommandsPage(nil)
		assert.NotNil(t, m)
		assert.NotNil(t, m.commands)
		assert.NotNil(t, m.keys)
	})

	t.Run("Init with Commands", func(t *testing.T) {
		assert.Nil(t, model().Init())
	})

	t.Run("Init Initial", func(t *testing.T) {
		cmd := CommandsPage{}.Init()

		assert.NotNil(t, cmd)
		batch := cmd()
		assert.Len(t, batch, 2)
	})

	t.Run("Init Help Command", func(t *testing.T) {
		cmd := socketModel().Init()

		res := cmd()
		_ = res.(tea.BatchMsg)[0]()
		help := res.(tea.BatchMsg)[1]()

		assert.Equal(t, haproxy.ParsedHelp{{
			Name: "help",
			Help: "foo",
			Args: "",
		}}, help)
	})

	t.Run("View List", func(t *testing.T) {
		view := parsedModel().View()
		assert.Contains(t, view, "1 item")
		assert.Contains(t, view, "help")
		assert.Contains(t, view, "foo")
		assert.Contains(t, view, "s goto status page")
	})

	t.Run("Update Unknown Message", func(t *testing.T) {
		m := model()
		m, cmd := m.Update(tea.KeyMsg{})

		assert.NotNil(t, m)
		assert.Nil(t, cmd)
	})

	t.Run("Update ParsedHelp", func(t *testing.T) {
		m, cmd := model().Update(haproxy.ParsedHelp{haproxy.Command{Name: "help"}})

		assert.Nil(t, cmd)
		assert.Len(t, m.commands.Items(), 1)
	})

	t.Run("Update Resize", func(t *testing.T) {
		nm, cmd := model().Update(tea.WindowSizeMsg{Width: 100, Height: 100})

		assert.Nil(t, cmd)
		assert.Equal(t, 96, nm.commands.Height())
		assert.Equal(t, 100, nm.commands.Width())
	})

	t.Run("Update Goto Status Page", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

		assert.NotNil(t, cmd)
		assert.IsType(t, ActivateStatusPageCmd(), cmd)
	})

	t.Run("Update Quit", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		assert.NotNil(t, cmd)
		assert.IsType(t, tea.QuitMsg{}, cmd())
	})

	t.Run("Update Goto Execute Page", func(t *testing.T) {
		_, cmd := parsedModel().Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{}})

		assert.NotNil(t, cmd)

		//type juggling as the type is tea.sequenceMsg (which is private), we cast it to the underlying type []tea.Cmd
		cmds := reflect.ValueOf(cmd()).Convert(reflect.TypeOf([]tea.Cmd{})).Interface().([]tea.Cmd)
		assert.IsType(t, ActivateExecutePageCmd(), cmds[0])
		assert.IsType(t, haproxy.Command{}, cmds[1]())
	})

	t.Run("Supports", func(t *testing.T) {
		m := model()

		assert.True(t, m.Supports(haproxy.ParsedHelp{}, false))
		assert.True(t, m.Supports(tea.WindowSizeMsg{}, false))
		assert.True(t, m.Supports(list.FilterMatchesMsg{}, false))
		assert.True(t, m.Supports(tea.KeyMsg{}, true))
		assert.False(t, m.Supports(tea.KeyMsg{}, false))
	})
}
