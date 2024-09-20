package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"slices"
	"sync"
)

type Parser struct {
	threads  int
	filepath string
}

func NewsParser(threads int, filepath string) *Parser {
	return &Parser{threads: threads, filepath: filepath}
}

func (p *Parser) Start(handleFn func(ip IP), logFn func(v int)) {
	wg := &sync.WaitGroup{}

	bufChan := parseFileToChunks(p.filepath)
	for range p.threads {

		wg.Add(1)
		go func() {
			defer wg.Done()

			for lBuf := range bufChan {
				scan := bufio.NewScanner(bytes.NewReader(lBuf))

				chunkLinesCnt := 0

				for scan.Scan() {
					ip := parseIP(scan.Text())
					handleFn(ip)
					chunkLinesCnt++
				}

				logFn(chunkLinesCnt)
			}
		}()
	}
	wg.Wait()
}

func parseFileToChunks(filepath string) <-chan []byte {
	out := make(chan []byte)

	go func() {
		defer close(out)

		file, err := os.Open(filepath)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(file *os.File) {
			errFClose := file.Close()
			if errFClose != nil {
				panic(errFClose)
			}
		}(file)

		fReader := bufio.NewReader(file)

		var lastLine []byte
		for {
			buf := make([]byte, 1024*1024*10)
			copy(buf, lastLine)

			nBytes, errRead := fReader.Read(buf[len(lastLine):])
			if errors.Is(errRead, io.EOF) {
				break
			}
			if errRead != nil {
				panic(errRead)
			}

			buf = buf[:len(lastLine)+nBytes]

			lastLine = make([]byte, 0, 16)

			maxIdx := len(buf) - 1
			for i := range 16 {
				idx := maxIdx - i
				if idx < 0 {
					buf = nil
					break
				}

				if buf[idx] != 0x0A {
					lastLine = append(lastLine, buf[idx])
				} else {
					buf = buf[:len(buf)-i-1]
					break
				}
			}
			slices.Reverse(lastLine)
			if len(buf) == 0 {
				continue
			}

			out <- buf
		}
	}()

	return out
}
