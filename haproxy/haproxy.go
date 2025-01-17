package haproxy

import (
	"errors"
	"maps"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ParsedHelp []Command

func (p ParsedHelp) Contains(command string) bool {
	for _, c := range p {
		if c.Name == command {
			return true
		}
	}

	return false
}

func (p ParsedHelp) GetCommand(command string) *Command {
	for _, c := range p {
		if c.Name == command {
			return &c
		}
	}

	return nil
}

type Command struct {
	Name string
	Help string
	Args string
}

type Backend struct {
	Id      int
	Name    string
	Servers []Server
}

const (
	UP            = "UP"
	DOWN          = "DOWN"
	MAINT         = "MAINT"
	DRAIN         = "DRAIN"
	AVAILABLE     = "AVAILABLE"
	STOPPED       = "STOPPED"
	RUNNING       = "RUNNING"
	STOPPING      = "STOPPING"
	STARTING      = "STARTING"
	DISABLED      = "DISABLED"
	READY         = "READY"
	CHECK         = "CHECK"   // agent using checks
	UNKNOWN       = "UNKNOWN" // check info unknown
	NEUTRAL       = "NEUTRAL" // Valid check but no status information.
	FAILED        = "FAILED"
	PASSED        = "PASSED"
	RES_COND_PASS = "RES_COND_PASS" //Check reports the server doesn't want new sessions.
	ENABLED       = "ENABLED"
	PAUSED        = "PAUSED"
	AGENT         = "AGENT" //Check is an agent check (otherwise it's a health check).
)

type Server struct {
	Id               int
	Name             string
	Address          net.IP
	State            string
	AdminState       string
	UserWeight       int
	CalculatedWeight int
	LastStateChanged time.Time
	SrvCheckStatus   string
	SrvCheckResult   string
	ChecksSucceeded  int
	CheckState       string
	SrvAgentState    string
	BkFForcedID      int
	SrvFForcedID     int
	Fqdn             string
	Port             int
	SrvRecord        string
	UseSSL           bool
	CheckPort        int
	CheckAddr        string
	AgentAddr        string
	AgentPort        int
}

func (c Command) FilterValue() string {
	return c.Name
}

func (c Command) Title() string {
	return c.Name
}
func (c Command) Description() string {
	return c.Help
}

func ParseHelp(rawHelp *string) ParsedHelp {
	cmds := strings.Split(
		*rawHelp,
		"\n",
	)

	var parsedHelp ParsedHelp

	for _, cmd := range cmds[1:] {
		if cmd == "" {
			continue
		}
		parts := strings.Split(cmd, ":")
		help := strings.TrimSpace(parts[1])
		command := strings.TrimSpace(parts[0])
		var args string

		if strings.Contains(cmd, "DEPRECATED") {
			continue
		}

		argStart := findArgumentsStart(command)
		if argStart > 0 {
			args = strings.TrimSpace(command[argStart:])
			command = strings.TrimSpace(command[:argStart])
		}

		parsedHelp = append(parsedHelp, Command{
			Name: command,
			Help: help,
			Args: args,
		})
	}

	return parsedHelp
}

func ParseBackends(input string) []Backend {
	backends := make(map[string]Backend)
	lines := strings.Split(input, "\n")

	// Skip lines with comments (#) or empty lines
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into fields
		fields := strings.Fields(line)

		if len(fields) == 1 {
			//TODO what is this "1" ?
			continue
		}

		backend := Backend{
			Id:   strToInt(fields[0]), // be_id
			Name: fields[1],           // ba_name
		}

		// see https://docs.haproxy.org/3.1/management.html for number permutations
		server := Server{
			Id:               strToInt(fields[2]),               // srv_id
			Name:             fields[3],                         //srv_name
			Address:          strToIp(fields[4]),                // srv_addr
			State:            strToServerState(fields[5]),       // srv_op_state
			AdminState:       strToAdminState(fields[6]),        // srv_admin_state
			UserWeight:       strToInt(fields[7]),               // srv_uweight
			CalculatedWeight: strToInt(fields[8]),               // srv_iweight
			LastStateChanged: strToTime(fields[9]),              // srv_time_since_last_change
			SrvCheckStatus:   strToCheckState(fields[10]),       // srv_check_status
			SrvCheckResult:   strToCheckResultState(fields[11]), // srv_check_result
			ChecksSucceeded:  strToInt(fields[12]),              // srv_check_health
			CheckState:       strToCheckInfoState(fields[13]),   // srv_check_state
			SrvAgentState:    strToAgentState(fields[14]),       // srv_agent_state
			BkFForcedID:      strToInt(fields[15]),
			SrvFForcedID:     strToInt(fields[16]),
			Fqdn:             fields[17],
			Port:             strToInt(fields[18]),
			SrvRecord:        fields[19],
			UseSSL:           strToBool(fields[20]),
			CheckPort:        strToInt(fields[21]),
			CheckAddr:        fields[22],
			AgentAddr:        fields[23],
			AgentPort:        strToInt(fields[24]),
		}
		if _, ok := backends[backend.Name]; !ok {
			backends[backend.Name] = backend
		}

		backend = backends[backend.Name]
		backend.Servers = append(backend.Servers, server)
		backends[backend.Name] = backend
	}

	return slices.Collect(maps.Values(backends))
}

func strToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func strToIp(s string) net.IP {
	if s == "-" {
		return nil
	}

	i := net.ParseIP(s)

	if i == nil {
		panic(errors.New("invalid ip"))
	}

	return i
}

func strToBool(s string) bool {
	if s == "1" {
		return true
	}

	return false
}

func strToServerState(s string) string {
	switch s {
	case "0":
		return STOPPED
	case "1":
		return STARTING
	case "2":
		return RUNNING
	case "3":
		return STOPPING
	default:
		return s
	}
}

func strToAgentState(s string) string {
	//bitmask
	const (
		running = 1 << iota
		configured
		enabled
		paused
		is_agent
	)
	state := strToInt(s)

	if (state & (configured | enabled)) == (configured | enabled) {
		return ENABLED
	}

	if state&paused != 0 {
		return PAUSED
	}

	if state&is_agent != 0 {
		return AGENT
	}

	return s
}

func strToCheckState(s string) string {
	switch s {
	case "0":
		return DOWN
	case "1":
		return UP
	default:
		return s
	}
}

func strToAdminState(s string) string {
	const (
		f_maint = 1 << iota
		i_maint
		c_maint
		f_drain
		i_drain
		r_maint
		h_maint
	)
	state := strToInt(s)

	if state&(f_maint|i_maint|c_maint|h_maint|r_maint|h_maint) != 0 {
		return MAINT
	}

	if state&(f_drain|i_drain) != 0 {
		return DRAIN
	}

	return s
}
func strToCheckInfoState(s string) string {
	//bitmask
	const (
		running = 1 << iota
		configured
		enabled
		paused
	)
	state := strToInt(s)

	if (state & (configured | enabled)) == (configured | enabled) {
		return ENABLED
	}

	if state&paused != 0 {
		return PAUSED
	}

	return DISABLED
}

func strToCheckResultState(s string) string {
	switch s {
	case "0":
		return UNKNOWN
	case "1":
		return NEUTRAL
	case "2":
		return FAILED
	case "3":
		return PASSED
	case "4":
		return RES_COND_PASS
	default:
		return s
	}
}

func strToTime(s string) time.Time {
	return time.Unix(time.Now().Unix()-int64(strToInt(s)), 0)
}

func findArgumentsStart(command string) int {
	brace1 := strings.Index(command, "[")
	brace2 := strings.Index(command, "<")
	brace4 := strings.Index(command, "{")

	var all []int
	for _, b := range []int{brace1, brace2, brace4} {
		if b > 0 {
			all = append(all, b)
		}
	}

	if len(all) == 0 {
		return -1
	}

	return slices.Min(all)
}
