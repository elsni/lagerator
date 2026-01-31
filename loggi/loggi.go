package loggi

import "fmt"

var Log = NewLoggi(true)

type Loggi struct {
	debug   bool
	loglist []string
}

// NewLoggi creates a logger with optional debug output.
func NewLoggi(debug bool) Loggi {
	l := Loggi{debug: debug}
	l.loglist = make([]string, 0, 120)
	l.loglist = append(l.loglist, "Loggi gestartet -----------")
	return l
}

// Log appends a message to the log buffer.
func (l *Loggi) Log(message string) {
	l.loglist = append(l.loglist, message)
}

// Print prints the log buffer when debug and doit are true.
func (l *Loggi) Print(doit bool) {
	if l.debug && doit {
		for _, m := range l.loglist {
			fmt.Println(m)
		}
	}
}
