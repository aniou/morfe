
package main

import (
        "fmt"
        "time"
        "github.com/aniou/go65c816/emulator/platform"
        "github.com/aniou/go65c816/lib/mylog"
)


func main() {
        logger := mylog.New()
        p := platform.New()
        p.Init(logger)
        p.LoadHex("/home/aniou/c256/go65c816/data/bench1.hex")
        p.CPU.PC = 0x0000
        p.CPU.RK = 0x03

	var running bool = false

	start := time.Now()
	for ! running {
		_, running = p.CPU.Step()
	}
	elapsed := time.Since(start)
	fmt.Printf("CPU stop after %d steps, took %s\n", p.CPU.AllCycles, elapsed)

}



