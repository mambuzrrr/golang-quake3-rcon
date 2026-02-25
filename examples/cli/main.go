package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"q3rcon/examples"
	"q3rcon/q3rcon"
)

func main() {
	cfg, err := examples.LoadConfig("config.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	c, err := q3rcon.New(
		cfg.Address,
		cfg.Password,
		q3rcon.Debug(cfg.Debug),
		q3rcon.Timeout(examples.MsOrDefault(cfg.TimeoutMs, 2000)),
		q3rcon.QuietWindow(examples.MsOrDefault(cfg.QuietMs, 180)),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer c.Close()

	fmt.Println("initialized. write rcon commands (enter to send). type 'exit' to quit.")

	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for {
		fmt.Print("> ")
		if !in.Scan() {
			break
		}

		cmd := strings.TrimSpace(in.Text())
		if cmd == "" {
			continue
		}
		if strings.EqualFold(cmd, "exit") || strings.EqualFold(cmd, "quit") {
			return
		}

		resp, err := c.Send(cmd)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println("server:", resp)
	}

	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
