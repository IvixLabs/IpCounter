package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type Logger struct {
	linesCnt   int
	muLinexCnt *sync.Mutex
	close      context.CancelFunc
}

func NewLogger() *Logger {
	logger := &Logger{
		muLinexCnt: &sync.Mutex{},
	}
	logger.start()

	return logger
}

func (l *Logger) Close() {
	l.close()
}

func (l *Logger) Len() int {
	return l.linesCnt
}

func (l *Logger) Inc(v int) {
	l.muLinexCnt.Lock()
	l.linesCnt += v
	l.muLinexCnt.Unlock()
}

func (l *Logger) start() {
	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("LINES", l.linesCnt)
			case <-ctx.Done():
				return
			}

		}
	}()
	l.close = cancel
}
