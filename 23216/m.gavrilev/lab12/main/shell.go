package main

import (
	"fmt"
	"os"
	"syscall"

	parce "lab12/parceline"
	prom "lab12/promptline"
)

func main() {

	var i int
	const lineLen int64 = 100
	line := make([]byte, lineLen, lineLen)
	var ncmds int
	var prompt string
	/* PLACE SIGNAL CODE HERE */

	prompt = fmt.Sprintf("[%s] ", os.Args[0])

	for prom.Promptline(prompt, line, lineLen) > 0 {
		ncmds = parce.Parceline(line)
		if ncmds <= 0 {
			continue
		}
		for i = 0; i < ncmds; i++ {
			_, err := syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{Env: append(os.Environ(), "DAEMON=true"),
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
