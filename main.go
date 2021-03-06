package main

import (
	"log"
	"runtime"
	"time"

	"github.com/auturnn/go-scrap/run"
)

func main() {
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("Start Log... cpu:", runtime.NumCPU())
	run.Run()
	log.Println("시작시간:", start.Format("2006-01-02 15:04:05"), "\n종료시간:",
		time.Now().Format("2006-01-02 15:04:05"), "\n총 소요시간:", time.Since(start))
}
