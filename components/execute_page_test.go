package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"net"
	"testing"
)

func TestExecutePage(t *testing.T) {
	t.Parallel()

	model := func() ExecutePage {
		return ExecutePage{}
	}

	socketModel := func() ExecutePage {
		conn := &socket.DummySocket{
			Output: []byte(`some commands response`),
		}

		return NewExecutePage(func() net.Conn { return conn })
	}

	t.Run("New", func(t *testing.T) {
		m := NewExecutePage(nil)
		assert.NotNil(t, m)
		assert.NotNil(t, m.input)
		assert.NotNil(t, m.keys)
	})

	t.Run("Init", func(t *testing.T) {
		cmd := socketModel().Init()
		assert.Nil(t, cmd)
	})

	t.Run("Update Unknown Message", func(t *testing.T) {
		m := model()
		m, cmd := m.Update(tea.KeyMsg{})

		assert.NotNil(t, m)
		assert.Nil(t, cmd)
	})

	t.Run("Update Command", func(t *testing.T) {
		m, cmd := socketModel().Update(haproxy.Command{
			Name: "foo", Help: "foo help text", Args: "<a>/<b>",
		})

		assert.Nil(t, cmd)
		assert.Equal(t, "foo ", m.input.Prompt)
		assert.Equal(t, "<a>/<b>", m.input.Placeholder)
	})

	t.Run("Update Goto Commands Page", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyBackspace, Runes: []rune{}})

		assert.NotNil(t, cmd)
		assert.IsType(t, ActivateCommandsPageCmd(), cmd)
	})

	t.Run("Update Quit", func(t *testing.T) {
		t.Skip("not implemented")
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		assert.NotNil(t, cmd)
		assert.IsType(t, tea.QuitMsg{}, cmd())
	})

	t.Run("Update Execute", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{}})

		assert.NotNil(t, cmd)
		assert.Equal(t, ExecuteResponse("some commands response"), cmd())
	})

	t.Run("View", func(t *testing.T) {
		m, _ := socketModel().Update(haproxy.Command{
			Name: "foo", Help: "bar help text", Args: "<a>/<b>",
		})

		res := m.View()

		assert.Contains(t, res, "bar help text")
		assert.Contains(t, res, "foo")
	})

	t.Run("Supports", func(t *testing.T) {
		m := model()

		assert.True(t, m.Supports(ExecuteResponse(""), false))
		assert.True(t, m.Supports(haproxy.Command{}, false))
		assert.True(t, m.Supports(tea.KeyMsg{}, true))
		assert.False(t, m.Supports(tea.KeyMsg{}, false))
	})
}
