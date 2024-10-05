package main

import (
	"fmt"
	"lab12/parceline"
	"lab12/promptline"
	"os"
	"syscall"
	"unsafe"
	//"os"
)

// main.go

func main() {

	// fmt.Println(parl.Parceline())
	// fmt.Println(prol.Promptline())
	// fmt.Println(head.MAXARGS)

	// flag.Parse()

	var i int
	var line [1024]byte
	var ncmds int
	var prompt string

	/* PLACE SIGNAL CODE HERE */

	//sprintf(prompt, "[%s] ", argv[0])
	prompt = fmt.Sprintf("[%s] ", os.Args[0])

	for promptline.Promptline(prompt, line, unsafe.Sizeof(line)) > 0 {
		ncmds = parceline.Parceline(line)
		if ncmds <= 0 {
			continue
		}
	}
	daemonENV := []string{"DAEMON=true"}
	for i = 0; i < ncmds; i++ {
		syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{Env: append(os.Environ(), daemonENV...),
			Sys: &syscall.SysProcAttr{
				Setsid: true,
			},
			Files: []uintptr{0, 1, 2}, // print message to the same pty
		})
	}

	//fmt.Println(prompt)

	// 	cmd := exec.Command("bash", "/c", "dir", "/a:-h/o:-n")
	// 	cmd.Stdin = strings.NewReader("some input")
	// 	var out strings.Builder
	// 	cmd.Stdout = &out
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%s\n", out.String())
	// }

	/*
	   fmt.Println("Hello World!")
	   	fmt.Println(parl.Parceline())
	   	fmt.Println(prol.Promptline())
	   	//var cm head.Command
	   	//fmt.Println(cm)
	   	//cm2 := head.Command{cmdargs: []string{"ULTRAKILL", "SOME"}}
	   	cm2 := head.Command{Cmdargs: [head.MAXARGS]string{"taasdasdasdasdsdvjsdfijvsdfijvnsdfvbisus", "fas"}, Cmdflag: '1'}
	   	fmt.Println(cm2)
	   	fmt.Println(head.Cmds)
	   	fmt.Println(head.Infile)
	   	parl.Change()
	   	fmt.Println(head.Infile)
	   	//fmt.Println(head.MAXARGS)
	*/
}
