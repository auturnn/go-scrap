package main

import (
	"log"
	"time"

	"github.com/auturnn/go-scrap/run"
)

func main() {
	start := time.Now()
	// runtime.GOMAXPROCS(12)
	log.Println("Start Log...")
	run.Run()
	log.Println("시작시간:", start.Format("2006-01-02 15:04:05"), "\n종료시간:",
		time.Now().Format("2006-01-02 15:04:05"), "\n총 소요시간:", time.Since(start))
}
