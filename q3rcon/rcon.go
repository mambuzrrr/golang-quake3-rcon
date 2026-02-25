package q3rcon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const header = "\xff\xff\xff\xff"

type Client struct {
	mu   sync.Mutex
	conn *net.UDPConn
	addr *net.UDPAddr
	opt  Options
}

// New makes a tiny udp rcon client.
// addr is "ip:port" (your game port), password is rcon_password.
func New(addr, password string, opts ...Option) (*Client, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, errors.New("q3rcon: addr is empty")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("q3rcon: resolve addr: %w", err)
	}

	o := DefaultOptions()
	o.Password = password
	for _, apply := range opts {
		apply(&o)
	}

	// dial udp: we get an ephemeral local port, server replies to that.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("q3rcon: dial udp: %w", err)
	}

	c := &Client{conn: conn, addr: udpAddr, opt: o}

	if c.opt.Debug && c.opt.Logf != nil {
		c.opt.Logf("q3rcon: connected to %s", udpAddr.String())
	}

	return c, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

// Send sends a command and waits for the response.
// it collects multiple udp packets until the server goes "quiet" for QuietWindow.
func (c *Client) Send(cmd string) (string, error) {
	return c.SendContext(context.Background(), cmd)
}

// SendContext is the same but cancellable.
func (c *Client) SendContext(ctx context.Context, cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", errors.New("q3rcon: cmd is empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return "", errors.New("q3rcon: client is closed")
	}

	if c.opt.Debug && c.opt.Logf != nil {
		c.opt.Logf("q3rcon: >> %q", cmd)
	}

	if _, err := c.conn.Write(buildPacket(c.opt.Password, cmd)); err != nil {
		return "", fmt.Errorf("q3rcon: write: %w", err)
	}

	raw, err := c.readUntilQuiet(ctx)
	if err != nil {
		return "", err
	}

	out := clean(string(raw))

	if c.opt.Debug && c.opt.Logf != nil {
		preview := out
		if len(preview) > 200 {
			preview = preview[:200] + "…"
		}
		c.opt.Logf("q3rcon: << %q", preview)
	}

	return out, nil
}

func buildPacket(password, cmd string) []byte {
	// classic q3 rcon format:
	// 0xff 0xff 0xff 0xff + "rcon <pass> <cmd>\n"
	payload := "rcon " + password + " " + cmd + "\n"
	return append([]byte(header), []byte(payload)...)
}

func (c *Client) readUntilQuiet(ctx context.Context) ([]byte, error) {
	timeout := c.opt.Timeout
	if timeout <= 0 {
		timeout = 1200 * time.Millisecond
	}

	quiet := c.opt.QuietWindow
	if quiet <= 0 {
		quiet = 140 * time.Millisecond
	}

	readBufSize := c.opt.ReadBuffer
	if readBufSize <= 0 {
		readBufSize = 64 * 1024
	}

	overallDeadline := time.Now().Add(timeout)

	buf := make([]byte, readBufSize)
	var out []byte

	packets := 0

	for {
		// context cancel is allowed to return partial data (nice for tooling).
		select {
		case <-ctx.Done():
			if len(out) > 0 {
				return out, nil
			}
			return nil, ctx.Err()
		default:
		}

		if time.Now().After(overallDeadline) {
			if len(out) > 0 {
				return out, nil
			}
			return nil, errors.New("q3rcon: timeout (no response)")
		}

		// first packet: wait up to the overall deadline
		// after that: we only wait a short "quiet window"
		var rd time.Time
		if len(out) == 0 {
			rd = overallDeadline
		} else {
			rd = time.Now().Add(quiet)
		}
		_ = c.conn.SetReadDeadline(rd)

		n, err := c.conn.Read(buf)
		if err != nil {
			// after we have *some* data, a timeout basically means "ok done"
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if len(out) > 0 {
					return out, nil
				}
				return nil, errors.New("q3rcon: timeout (no response)")
			}
			return nil, fmt.Errorf("q3rcon: read: %w", err)
		}

		packets++
		if c.opt.MaxPackets > 0 && packets > c.opt.MaxPackets {
			return out, errors.New("q3rcon: max packets reached")
		}

		if c.opt.Debug && c.opt.Logf != nil {
			c.opt.Logf("q3rcon: << packet %d (%d bytes)", packets, n)
		}

		out = append(out, buf[:n]...)
	}
}

func clean(s string) string {
	if s == "" {
		return ""
	}

	// responses often come like:
	// \xff\xff\xff\xffprint\n.... (and repeated per packet)
	parts := strings.Split(s, header)
	if len(parts) > 1 {
		var b strings.Builder
		for _, p := range parts {
			if p == "" {
				continue
			}
			b.WriteString(p)
		}
		s = b.String()
	}

	s = strings.ReplaceAll(s, "\r\n", "\n")

	// common prefix
	s = strings.TrimPrefix(s, "print\n")
	s = strings.TrimPrefix(s, "print")

	// trim annoying tails
	s = strings.Trim(s, "\x00")
	s = strings.TrimSpace(s)

	return s
}
