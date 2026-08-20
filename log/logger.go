package log

import (
	stdlog "log"
	"os"
)

/*
 * Solarman logging interface
 *
 * Users can implement their own loagger and use it
 */
type Logger interface {
	Infof(format string, args ...any)
	Info(msg string)
	Errorf(format string, args ...any)
	Error(msg string)
	Warnf(format string, args ...any)
	Warn(nsg string)
	Debugf(format string, args ...any)
	Debug(msg string)
}

// SilentLoger dummy logger
type SilentLoger struct{}

func (s *SilentLoger) Infof(format string, args ...any)  {}
func (s *SilentLoger) Info(msg string)                   {}
func (s *SilentLoger) Errorf(format string, args ...any) {}
func (s *SilentLoger) Error(msg string)                  {}
func (s *SilentLoger) Warnf(format string, args ...any)  {}
func (s *SilentLoger) Warn(msg string)                   {}
func (s *SilentLoger) Debugf(format string, args ...any) {}
func (s *SilentLoger) Debug(msg string)                  {}

type solarmanLogger struct {
	li *stdlog.Logger
	le *stdlog.Logger
	lw *stdlog.Logger
	ld *stdlog.Logger
}

func (s *solarmanLogger) Infof(format string, args ...any)  { s.li.Printf(format, args...) }
func (s *solarmanLogger) Info(msg string)                   { s.li.Print(msg) }
func (s *solarmanLogger) Errorf(format string, args ...any) { s.le.Printf(format, args...) }
func (s *solarmanLogger) Error(msg string)                  { s.le.Print(msg) }
func (s *solarmanLogger) Warnf(format string, args ...any)  { s.lw.Printf(format, args...) }
func (s *solarmanLogger) Warn(msg string)                   { s.lw.Print(msg) }
func (s *solarmanLogger) Debugf(format string, args ...any) { s.ld.Printf(format, args...) }
func (s *solarmanLogger) Debug(msg string)                  { s.ld.Print(msg) }

func NewSolarmanLogger() Logger {
	l := new(solarmanLogger)
	l.li = stdlog.New(os.Stdout, "<solarman>  [INFO] ", stdlog.Lmicroseconds|stdlog.LstdFlags)
	l.le = stdlog.New(os.Stdout, "<solarman> [ERROR] ", stdlog.Lmicroseconds|stdlog.LstdFlags)
	l.lw = stdlog.New(os.Stdout, "<solarman>  [WARN] ", stdlog.Lmicroseconds|stdlog.LstdFlags)
	l.ld = stdlog.New(os.Stdout, "<solarman> [DEBUG] ", stdlog.Lmicroseconds|stdlog.LstdFlags)

	return l
}
