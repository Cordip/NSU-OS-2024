package promptline

import (
	"bufio"
	"os"
)

func Promptline(prompt string, line []byte , sizline uintptr) int {

	var len int = 0

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)

	_, err := out.WriteString(prompt)
	if err != nil {
		panic("Write Problem")
	}

	for {
		n, err := in.Read(line)
		if err != nil {
			panic("Read Problem")
		}
		len += n
		line[len] = 0
		if(line[len-2] != '\\' && line[len-1] != '\n'){
			line[len] = ' '
			line[len-1] = ' '
			line[len-2] = ' '
			continue
		}
		return(len)
	}
}