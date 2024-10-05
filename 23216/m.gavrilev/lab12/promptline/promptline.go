package promptline

//"fmt"
import (
	"bufio"
	"fmt"
	head "lab12/header"
	"io"
	"os"
	"syscall"
)

func Promptline(prompt string, line string, sizline uintptr) int {
	var n int = 0

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)

	_, err := out.WriteString(prompt)
	if err == io.EOF {
		syscall.Exit(1)
	}
	
	for {
		a, err := in.ReadString()
		if err == io.EOF {
			syscall.Exit(1)
		}
		n += a
	}
	head.Infile = "asdfrmasd"
	return 1
}
