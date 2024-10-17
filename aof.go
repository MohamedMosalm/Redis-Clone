package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	wr   *bufio.Writer
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		wr:   bufio.NewWriter(f),
	}

	go func() {
		for {

			aof.mu.Lock()
			aof.wr.Flush()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second * 5)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.wr.Flush()

	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	data := value.Marshal()

	_, err := aof.wr.Write(data)
	return err
}
