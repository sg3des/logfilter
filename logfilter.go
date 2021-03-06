package logfilter

import (
	"bytes"
	"io"
	"os"
)

type LogPrefix string

type Filter struct {
	writers []logWriter
}

func NewFilter() *Filter {
	return &Filter{}
}

type logWriter struct {
	filename string
	w        io.Writer
	strict   bool
	prefixes []LogPrefix
}

func (l *Filter) AddWriter(w io.Writer, prefixes ...LogPrefix) {
	l.writers = append(l.writers, logWriter{
		w:        w,
		strict:   false,
		prefixes: prefixes,
	})
}

func (l *Filter) AddStrictWriter(w io.Writer, prefixes ...LogPrefix) {
	l.writers = append(l.writers, logWriter{
		w:        w,
		strict:   true,
		prefixes: prefixes,
	})
}

func (l *Filter) AddFileWriter(filename string, prefixes ...LogPrefix) {
	l.writers = append(l.writers, logWriter{
		filename: filename,
		strict:   false,
		prefixes: prefixes,
	})
}

func (l *Filter) AddStrictFileWriter(filename string, prefixes ...LogPrefix) {
	l.writers = append(l.writers, logWriter{
		filename: filename,
		strict:   true,
		prefixes: prefixes,
	})
}

func (lw *logWriter) Check(p []byte, prefix LogPrefix) bool {
	if len(lw.prefixes) == 0 {
		return true
	}

	if prefix == "" {
		return !lw.strict
	}

	for _, l := range lw.prefixes {
		if prefix == l {
			return true
		}
	}

	return false
}

func (lw *logWriter) Write(p []byte) (int, error) {
	if lw.filename != "" && lw.w == nil {
		var err error
		lw.w, err = os.OpenFile(lw.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return 0, err
		}
	}

	return lw.w.Write(p)
}

func (l *Filter) getPrefix(p []byte) LogPrefix {
	x := bytes.IndexByte(p, '[')
	if x >= 0 {
		y := bytes.IndexByte(p[x:], ']')
		if y >= 0 {
			return LogPrefix(p[x+1 : x+y])
		}
	}

	return LogPrefix("")
}

func (l *Filter) Write(p []byte) (int, error) {
	if len(l.writers) == 0 {
		return os.Stdout.Write(p)
	}

	prefix := l.getPrefix(p)

	for _, lw := range l.writers {
		if lw.Check(p, prefix) {
			n, err := lw.Write(p)
			if err != nil {
				return n, err
			}
		}
	}

	return len(p), nil
}
