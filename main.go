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
	log.Println("시작시간:", start, "종료시간:", time.Now(), "총 소요시간:", time.Since(start))
}
