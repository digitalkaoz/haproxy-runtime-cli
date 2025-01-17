package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"haproxy-runtime-cli/haproxy"
	"haproxy-runtime-cli/socket"
	"net"
	"testing"
)

func TestStatusPage(t *testing.T) {
	t.Parallel()

	model := func() StatusPage {
		return StatusPage{}
	}

	socketModel := func() StatusPage {
		conn := &socket.DummySocket{
			Output: []byte(`1
			# be_id be_name srv_id srv_name srv_addr srv_op_state srv_admin_state srv_uweight srv_iweight srv_time_since_last_change srv_check_status srv_check_result srv_check_health srv_check_state srv_agent_state bk_f_forced_id srv_f_forced_id srv_fqdn srv_port srvrecord srv_use_ssl srv_check_port srv_check_addr srv_agent_addr srv_agent_port
			4 default 1 haproxy 209.126.35.1 2 0 20 20 9 9 3 4 6 0 0 0 haproxy.com 443 - 1 0 - - 0
			5 other 1 haproxy 209.126.35.1 2 0 1 1 9 15 3 4 6 0 0 0 haproxy.com 443 - 1 0 - - 0
			`),
		}

		return NewStatusPage(func() net.Conn { return conn })
	}

	t.Run("New", func(t *testing.T) {
		m := NewStatusPage(nil)
		assert.NotNil(t, m)
		assert.NotNil(t, m.table)
		assert.NotNil(t, m.keys)
	})

	t.Run("Init", func(t *testing.T) {
		cmd := socketModel().Init()
		assert.NotNil(t, cmd)
	})

	t.Run("Init with fetch backends", func(t *testing.T) {
		cmd := socketModel().Init()

		res := cmd()

		assert.IsType(t, []haproxy.Backend{}, res)
		assert.Len(t, res.([]haproxy.Backend), 2)
	})

	t.Run("Update Unknown Message", func(t *testing.T) {
		m := model()
		m, cmd := m.Update(tea.KeyMsg{})

		assert.NotNil(t, m)
		assert.Nil(t, cmd)
	})

	t.Run("Update Backends", func(t *testing.T) {
		m, cmd := socketModel().Update([]haproxy.Backend{
			{Name: "foo", Id: 1, Servers: []haproxy.Server{{Name: "foo", Id: 1, Fqdn: "foo.com", Port: 443}}},
		})

		assert.Nil(t, cmd)
		assert.Len(t, m.table.Rows(), 2)
	})

	t.Run("Update Resize", func(t *testing.T) {
		nm, cmd := model().Update(tea.WindowSizeMsg{Width: 100, Height: 100})

		assert.Nil(t, cmd)
		assert.Equal(t, 92, nm.table.Height())
		assert.Equal(t, 96, nm.table.Width())
	})

	t.Run("Update Goto Commands Page", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyBackspace, Runes: []rune{}})

		assert.NotNil(t, cmd)
		assert.IsType(t, ActivateCommandsPageCmd(), cmd)
	})

	t.Run("Update Reload", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

		assert.NotNil(t, cmd)
		assert.IsType(t, []haproxy.Backend{}, cmd())
	})

	t.Run("Update Quit", func(t *testing.T) {
		_, cmd := socketModel().Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		assert.NotNil(t, cmd)
		assert.IsType(t, tea.QuitMsg{}, cmd())
	})

	t.Run("View", func(t *testing.T) {
		m, _ := socketModel().Update([]haproxy.Backend{
			{Name: "bar", Id: 1, Servers: []haproxy.Server{{Name: "foo", Id: 1, Fqdn: "foo.com", Port: 443, Address: net.ParseIP("127.0.0.1")}}},
		})

		res := m.View()

		assert.Contains(t, res, "127.0.0.1")
		assert.Contains(t, res, "bar")
		assert.Contains(t, res, "foo")
		assert.Contains(t, res, "foo.com:443")
	})

	t.Run("Supports", func(t *testing.T) {
		m := model()

		assert.True(t, m.Supports([]haproxy.Backend{}, false))
		assert.True(t, m.Supports(tea.WindowSizeMsg{}, false))
		assert.True(t, m.Supports(tea.KeyMsg{}, true))
		assert.False(t, m.Supports(tea.KeyMsg{}, false))
	})
}
