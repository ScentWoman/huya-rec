package main

import (
	"flag"
	"log"
	"time"

	rec "github.com/ScentWoman/huya-rec"
)

var (
	room  = flag.String("room", "", "live room url")
	retry = flag.Int("retry", 5, "retry interval in second")
	split = flag.Int("split", 1, "split interval in hour")
	out   = flag.String("o", "", "output path")
)

func init() {
	flag.Parse()

	if *room == "" {
		log.Fatalln("live room needed!")
	}
}

func main() {
	rec.Record(*room, time.Duration(*split)*time.Hour, time.Duration(*retry)*time.Second, *out)
}
