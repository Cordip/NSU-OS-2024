package promptline

//"fmt"
import (
	"fmt"
	head "lab12/header"
)

func Promptline(prompt string, line [1024]byte, sizline uintptr) int {
	var n int = 0

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	for {
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		}
		out.Write()// MAYBE WRITEBBYTER
	}

	fmt.WriteByte(prompt, strlen(prompt))
	
	for {
		n += 
	}
	head.Infile = "asdfrmasd"
	return 1
}
