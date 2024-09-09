// main.go
package main

import (
	"flag"
	"fmt"
	head "lab12/header"
	parl "lab12/parceline"
	prol "lab12/promptline"
	"log"

	//"os"
	"os/exec"
	"strings"
)

func main() {

	fmt.Println(parl.Parceline())
	fmt.Println(prol.Promptline())
	fmt.Println(head.MAXARGS)

	flag.Parse()

	// var i int
	// var line [1024]byte
	// var ncmds int
	//var prompt string

	/* PLACE SIGNAL CODE HERE */

	//sprintf(prompt, "[%s] ", argv[0])
	//prompt = fmt.Sprintf("[%s] ", os.Args[0])

	//fmt.Println(prompt)

	cmd := exec.Command("cmd", "/c", "dir", "/a:-h/o:-n")
	cmd.Stdin = strings.NewReader("some input")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out.String())
}

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
