package haproxy

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const rawHelp = `The following commands are valid at this level:
  abort ssl ca-file <cafile>              : abort a transaction for a CA file
  abort ssl cert <certfile>               : abort a transaction for a certificate file
  abort ssl crl-file <crlfile>            : abort a transaction for a CRL file
  add acl [@<ver>] <acl> <pattern>        : add an acl entry
  add map [@<ver>] <map> <key> <val>      : add a map entry (payload supported instead of key/val)
  add server <bk>/<srv>                   : create a new server
  add ssl ca-file <cafile> <payload>      : add a certificate into the CA file
  add ssl crt-list <list> <cert> [opts]*  : add to crt-list file <list> a line <cert> or a payload
  clear acl [@<ver>] <acl>                : clear the contents of this acl
  clear counters [all]                    : clear max statistics counters (or all counters)
  clear map [@<ver>] <map>                : clear the contents of this map
  clear table <table> [<filter>]*         : remove an entry from a table (filter: data/key)
  commit acl @<ver> <acl>                 : commit the ACL at this version
  commit map @<ver> <map>                 : commit the map at this version
  commit ssl ca-file <cafile>             : commit a CA file
  commit ssl cert <certfile>              : commit a certificate file
  commit ssl crl-file <crlfile>           : commit a CRL file
  debug counters [?|all|bug|cnt|chk|glt]* : dump/reset rare event counters
  debug dev hash   [msg]                  : return msg hashed if anon is set
  del acl <acl> [<key>|#<ref>]            : delete acl entries matching <key>
  del map <map> [<key>|#<ref>]            : delete map entries matching <key>
  del server <bk>/<srv>                   : remove a dynamically added server
  del ssl ca-file <cafile>                : delete an unused CA file
  del ssl cert <certfile>                 : delete an unused certificate file
  del ssl crl-file <crlfile>              : delete an unused CRL file
  del ssl crt-list <list> <cert[:line]>   : delete a line <cert> from crt-list file <list>
  disable agent                           : disable agent checks
  disable dynamic-cookie backend <bk>     : disable dynamic cookies on a specific backend
  disable frontend <frontend>             : temporarily disable specific frontend
  disable health                          : disable health checks
  disable server (DEPRECATED)             : disable a server for maintenance (use 'set server' instead)
  dump ssl cert <certfile>                : dump the SSL certificates in PEM format
  dump stats-file                         : dump stats for restore
  echo <text>                             : print text to the output
  enable agent                            : enable agent checks
  enable dynamic-cookie backend <bk>      : enable dynamic cookies on a specific backend
  enable frontend <frontend>              : re-enable specific frontend
  enable health                           : enable health checks
  enable server  (DEPRECATED)             : enable a disabled server (use 'set server' instead)
  get acl <acl> <value>                   : report the patterns matching a sample for an ACL
  get map <acl> <value>                   : report the keys and values matching a sample for a map
  get var <name>                          : retrieve contents of a process-wide variable
  get weight <bk>/<srv>                   : report a server's current weight
  new ssl ca-file <cafile>                : create a new CA file to be used in a crt-list
  new ssl cert <certfile>                 : create a new certificate file to be used in a crt-list or a directory
  new ssl crl-file <crlfile>               : create a new CRL file to be used in a crt-list
  operator                                : lower the level of the current CLI session to operator
  prepare acl <acl>                       : prepare a new version for atomic ACL replacement
  prepare map <acl>                       : prepare a new version for atomic map replacement
  set anon global-key <value>             : change the global anonymizing key
  set anon off                            : deactivate the anonymized mode
  set anon on [value]                     : activate the anonymized mode
  set dynamic-cookie-key backend <bk> <k> : change a backend secret key for dynamic cookies
  set map <map> [<key>|#<ref>] <value>    : modify a map entry
  set maxconn frontend <frontend> <value> : change a frontend's maxconn setting
  set maxconn global <value>              : change the per-process maxconn setting
  set maxconn server <bk>/<srv>           : change a server's maxconn setting
  set profiling <what> {auto|on|off}      : enable/disable resource profiling (tasks,memory)
  set rate-limit <setting> <value>        : change a rate limiting value
  set server <bk>/<srv> [opts]            : change a server's state, weight, address or ssl
  set severity-output [none|number|string]: set presence of severity level in feedback information
  set ssl ca-file <cafile> <payload>      : replace a CA file
  set ssl cert <certfile> <payload>       : replace a certificate file
  set ssl crl-file <crlfile> <payload>    : replace a CRL file
  set ssl ocsp-response <resp|payload>       : update a certificate's OCSP Response from a base64-encode DER
  set ssl tls-key [id|file] <key>         : set the next TLS key for the <id> or <file> listener to <key>
  set table <table> key <k> [data.* <v>]* : update or create a table entry's data
  set timeout [cli] <delay>               : change a timeout setting
  set weight <bk>/<srv>  (DEPRECATED)     : change a server's weight (use 'set server' instead)
  show acl [@<ver>] <acl>]                : report available acls or dump an acl's contents
  show activity [-1|0|thread_num]         : show per-thread activity stats (for support/developers)
  show anon                               : display the current state of anonymized mode
  show backend                            : list backends in the current running config
  show cache                              : show cache status
  show cli level                          : display the level of the current CLI session
  show cli sockets                        : dump list of cli sockets
  show dev                                : show debug info for developers
  show env [var]                          : dump environment variables known to the process
  show errors [<px>] [request|response]   : report last request and/or response errors for each proxy
  show events [<sink>] [-w] [-n]          : show event sink state
  show fd [-!plcfbsd]* [num]              : dump list of file descriptors in use or a specific one
  show info [desc|json|typed|float]*      : report information about the running process
  show map [@ver] [map]                   : report available maps or dump a map's contents
  show peers [dict|-] [section]           : dump some information about all the peers or this peers section
  show pools [by*] [match <pfx>] [nb]     : report information about the memory pools usage
  show profiling [<what>|<#lines>|<opts>]*: show profiling state (all,status,tasks,memory)
  show resolvers [id]                     : dumps counters from all resolvers section and associated name servers
  show schema json                        : report schema used for stats
  show servers conn [<backend>]           : dump server connections status (all or for a single backend)
  show servers state [<backend>]          : dump volatile server information (all or for a single backend)
  show sess [help|<id>|all|susp|older...] : report the list of current streams or dump this exact stream
  show ssl ca-file [<cafile>[:<index>]]   : display the SSL CA files used in memory, or the details of a <cafile>, or a single certificate of index <index> of a CA file <cafile>
  show ssl cert [<certfile>]              : display the SSL certificates used in memory, or the details of a file
  show ssl crl-file [<crlfile[:<index>>]] : display the SSL CRL files used in memory, or the details of a <crlfile>, or a single CRL of index <index> of CRL file <crlfile>
  show ssl crt-list [-n] [<list>]         : show the list of crt-lists or the content of a crt-list file <list>
  show ssl ocsp-response [[text|base64] id]  : display the IDs of the OCSP responses used in memory, or the details of a single OCSP response (in text or base64 format)
  show ssl ocsp-updates                      : display information about the next 'nb' ocsp responses that will be updated automatically
  show ssl providers                      : show loaded SSL providers
  show startup-logs                       : report logs emitted during HAProxy startup
  show stat [desc|json|no-maint|typed|up]*: report counters for each proxy and server
  show table <table> [<filter>]*          : report table usage stats or dump this table's contents (filter: data/key)
  show tasks                              : show running tasks
  show threads                            : show some threads debugging information
  show tls-keys [id|*]                    : show tls keys references or dump tls ticket keys when id specified
  show trace [<module>]                   : show live tracing state
  show version                            : show version of the current process
  shutdown frontend <frontend>            : stop a specific frontend
  shutdown session [id]                   : kill a specific session
  shutdown sessions server <bk>/<srv>     : kill sessions on a server
  trace [<module>|0] [cmd [args...]]      : manage live tracing (empty to list, 0 to stop all)
  update ssl ocsp-response <certfile>     : send ocsp request and update stored ocsp response
  user                                    : lower the level of the current CLI session to user
  wait {-h|<delay_ms>} cond [args...]     : wait the specified delay or condition (-h to see list)
  help [<command>]                        : list matching or all commands
  prompt [timed]                          : toggle interactive mode with prompt
  quit                                    : disconnect`

const sampleBackends = `1
# be_id be_name srv_id srv_name srv_addr srv_op_state srv_admin_state srv_uweight srv_iweight srv_time_since_last_change srv_check_status srv_check_result srv_check_health srv_check_state srv_agent_state bk_f_forced_id srv_f_forced_id srv_fqdn srv_port srvrecord srv_use_ssl srv_check_port srv_check_addr srv_agent_addr srv_agent_port
4 default 1 haproxy 209.126.35.1 2 0 20 20 9 9 3 4 6 0 0 0 haproxy.com 443 - 1 0 - - 0
4 default 2 apache 151.101.2.132 2 0 80 80 9 9 3 4 6 0 0 0 apache.org 443 - 1 0 - - 0
5 other 1 haproxy 209.126.35.1 2 0 1 1 9 15 3 4 6 0 0 0 haproxy.com 443 - 1 0 - - 0
5 other 2 apache 151.101.2.132 0 0 1 1 7 17 2 0 6 0 0 0 apache.org 443 - 1 0 - - 0
`

func TestParseHelp(t *testing.T) {
	input := rawHelp

	parsedHelp := ParseHelp(&input)

	assert.NotEmpty(t, parsedHelp)
	assert.Len(t, parsedHelp, 113)
}

func TestParsedHelpContains(t *testing.T) {
	help := ParsedHelp{
		Command{
			Name: "clear acl",
			Help: "clear the contents of this acl",
			Args: "[@<ver>] <acl>",
		},
	}

	assert.True(t, help.Contains("clear acl"))
	assert.False(t, help.Contains("unknown command"))
}

func TestParsedHelpGetCommand(t *testing.T) {
	cmd := Command{
		Name: "clear acl",
		Help: "clear the contents of this acl",
		Args: "[@<ver>] <acl>",
	}

	help := ParsedHelp{cmd}

	assert.Equal(t, cmd, *help.GetCommand("clear acl"))
	assert.Nil(t, help.GetCommand("unknown"))
}

func TestParseHelpCommand(t *testing.T) {
	input := rawHelp

	parsedHelp := ParseHelp(&input)

	cmd := parsedHelp.GetCommand("clear acl")

	assert.Equal(t, "clear acl", cmd.Name)
	assert.Equal(t, "clear the contents of this acl", cmd.Help)
	assert.Equal(t, "[@<ver>] <acl>", cmd.Args)
}

func TestSplitArgumentsFromCommand(t *testing.T) {
	cmd := "clear acl [arg1] [arg2]"
	assert.Equal(t, 10, findArgumentsStart(cmd))
	assert.Equal(t, "clear acl", strings.TrimSpace(cmd[:10]))

	cmd = "add ssl ca-file <cafile> <payload>"
	assert.Equal(t, 16, findArgumentsStart(cmd))
	assert.Equal(t, "add ssl ca-file", strings.TrimSpace(cmd[:16]))

	cmd = "wait {-h|<delay_ms>} cond [args...]"
	assert.Equal(t, 5, findArgumentsStart(cmd))
	assert.Equal(t, "wait", strings.TrimSpace(cmd[:5]))
}

func TestParseBackends(t *testing.T) {
	res := ParseBackends(sampleBackends)

	assert.Len(t, res, 2)
	assert.Len(t, res[0].Servers, 2)
	assert.Len(t, res[1].Servers, 2)
}
