package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/wasilibs/go-re2"
)

func worker(name string, msgch chan *nats.Msg) {
	log.Printf("consumer running with name %s\n", name)
	rx := `@`
	r := re2.MustCompile(rx)

	count := 0
	for m := range msgch {
		count++
		// println(m.Size(), len(m.Data))

		// 假设每个包匹配100次
		for i := 0; i < 100; i++ {
			r.FindIndex(m.Data)
		}

		if count%1000 == 0 {
			println(m.Size(), m.Header, string(m.Data))
			log.Printf("consumer %s count: %d, find index %v\n", name, count, 0)
		}
	}
}
