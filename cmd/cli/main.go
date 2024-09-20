package main

import (
	"flag"
	"log"
	"os"
)

var filepath string
var threads int

func init() {
	flag.StringVar(&filepath, "path", "", "path ips file")
	flag.IntVar(&threads, "threads", 12, "total file reading threads. Tune it while your cpu will not fully loaded.")
}

func main() {
	flag.Parse()

	var err error
	if filepath != "" {
		var file *os.File
		file, err = os.Open(filepath)
		if err != nil {
			log.Println(err)
		}
		file.Close()
	}

	if filepath == "" || err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	parser := NewsParser(threads, filepath)

	proc := NewIpCollectorProcessor()

	logger := NewLogger()
	log.Println("START")

	parser.Start(
		func(ip IP) {
			proc.HandleIP(ip)
		},
		func(v int) {
			logger.Inc(v)
		})

	logger.Close()

	proc.Close()
	proc.Wait()
	log.Println("FINISH")

	log.Println("TOTAL LINES", logger.Len())
	log.Println("UNIQ LINES", proc.CalcUniqIPs())
}
