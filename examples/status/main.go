package main

import (
	"fmt"
	"os"

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

	resp, err := c.Send("status")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(resp)
}
