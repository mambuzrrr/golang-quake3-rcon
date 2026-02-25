package q3rcon

import (
	"log"
	"time"
)

type Options struct {
	Password string

	Timeout     time.Duration
	QuietWindow time.Duration
	ReadBuffer  int
	MaxPackets  int // 0 = unlimited

	Debug bool
	Logf  func(format string, args ...any)
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		Timeout:     1200 * time.Millisecond,
		QuietWindow: 140 * time.Millisecond,
		ReadBuffer:  64 * 1024,
		MaxPackets:  0,

		Debug: false,
		Logf:  log.Printf,
	}
}

func Timeout(d time.Duration) Option     { return func(o *Options) { o.Timeout = d } }
func QuietWindow(d time.Duration) Option { return func(o *Options) { o.QuietWindow = d } }
func Debug(enabled bool) Option          { return func(o *Options) { o.Debug = enabled } }
func MaxPackets(n int) Option {
	return func(o *Options) {
		if n >= 0 {
			o.MaxPackets = n
		}
	}
}

func ReadBuffer(n int) Option {
	return func(o *Options) {
		if n > 0 {
			o.ReadBuffer = n
		}
	}
}

func Logger(logf func(format string, args ...any)) Option {
	return func(o *Options) { o.Logf = logf }
}
