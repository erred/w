package stdlog

import (
	"errors"
	"log"

	"github.com/go-logr/logr"
)

func New(l logr.Logger, errorLevel bool) *log.Logger {
	w := &writer{
		l, errorLevel,
	}
	return log.New(w, "", 0)
}

type writer struct {
	l          logr.Logger
	errorLevel bool
}

func (w *writer) Write(b []byte) (int, error) {
	if w.errorLevel {
		w.l.Error(errors.New(string(b)), "")
	} else {
		w.l.Info(string(b))
	}
	return len(b), nil
}
