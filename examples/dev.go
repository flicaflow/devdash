package main

import (
	"fmt"
	"time"

	dd "github.com/flicaflow/devdash"
)

func main() {
	dd.Start()

	dd.Message("A", "Hallo")
	time.Sleep(2 * time.Second)
	dd.Message("A", "Hallo Welt")
	time.Sleep(2 * time.Second)

	go func() {
		for i := time.Duration(0); i < 100; i++ {
			dd.Message("B", fmt.Sprint("Hallo Others ", i))
			time.Sleep(i * time.Second)
		}

	}()

	for i := time.Duration(0); i < 100; i++ {
		dd.Message("A", fmt.Sprint("Hallo Welt ", i))
		time.Sleep(2 * (i * time.Second))
	}

	fmt.Println("vim-go")
}
