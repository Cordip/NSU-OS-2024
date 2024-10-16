package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	head "../lab12/header"
	parce "../lab12/parceline"
	prom "../lab12/promptline"
)

func main() {

	var i int
	var line []byte
	var ncmds int
	var prompt string

	/* PLACE SIGNAL CODE HERE */

	prompt = fmt.Sprintf("[%s] ", os.Args[0])

	for prom.Promptline(prompt, line, unsafe.Sizeof(line)) > 0 {
		ncmds = parce.Parceline(line)
		if ncmds <= 0 {
			continue
		}
		for i = 0; i < ncmds; i++ {
			_, err := os.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{Env: append(os.Environ(), "DAEMON=true"),
				Sys: &syscall.SysProcAttr{
					Setsid: true,
				},
				Files: []uintptr{0, 1, 2}, // print message to the same pty
			})
			if err != nil {
				panic("ForkExec Problem")
			}
		}
	}
}
