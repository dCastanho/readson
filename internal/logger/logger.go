package logger

import (
	"io"
	"log"
)

type logger struct {
	nodeLogger  *log.Logger
	blockLogger *log.Logger
	initiated   bool
}

func (l *logger) Node(v ...interface{}) {
	if l.initiated {
		l.nodeLogger.Println(v...)
	}
}

func (l *logger) Block(v ...interface{}) {
	if l.initiated {
		l.blockLogger.Println(v...)
	}
}

type Logger interface {
	Node(v ...interface{})
	Block(v ...interface{})
}

var DefaultLogger Logger

func DeployLogger(initiated bool, filepath io.Writer) {
	DefaultLogger = &logger{
		nodeLogger:  log.New(filepath, "EVALUATING: ", 0),
		blockLogger: log.New(filepath, "BLOCK: ", 0),
		initiated:   initiated,
	}
}
