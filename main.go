package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"haproxy-runtime-cli/socket"
	"haproxy-runtime-cli/styles"
	"log"
	"net"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	args := os.Args[1:]
	parseCommandLine(args)

	openSocket := func() net.Conn {
		return socket.Listen(args[0])
	}

	p := tea.NewProgram(NewRuntimeApi(openSocket), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(styles.ErrorStyle.Render(fmt.Sprintf("Alas, there's been an error: %v", err)))
	}
}

func parseCommandLine(args []string) {
	if len(args) == 0 {
		log.Fatal(styles.ErrorStyle.Render("Please specify a haproxy socket as argument"))
	}

	if args[0] == "-v" || args[0] == "--version" {
		fmt.Println(styles.ActiveStyle.Render(fmt.Sprintf("%s, commit %s, built at %s by digitalkaoz", version, commit, date)))
		os.Exit(0)
	}

	stat, err := os.Stat(args[0])
	if err != nil {
		log.Fatal(styles.ErrorStyle.Render(err.Error()))
	}

	if stat.IsDir() {
		log.Fatal(styles.ErrorStyle.Render(fmt.Sprintf(`%s is not a valid haproxy socket`, args[0])))
	}
}
