package promptline

import (
	"fmt"
	"syscall"
)

func Promptline(prompt string, line []byte, sizline int64) int64 {

	var len int64 = 0
	tmpLine := make([]byte, 1024)
	hihi := []byte(prompt)
	_, err := syscall.Write(1, hihi)
	if err != nil {
		panic("Write Problem")
	}
	for {
		n, err := syscall.Read(0, tmpLine)
		line = append(line[:len], tmpLine[:n]...)
		fmt.Println(line)
		fmt.Println(err)
		if err != nil {
			panic("Read Problem")
		}
		b := int64(n)
		len += b
		if line[len-2] == '\\' && line[len-1] == '\n' {
			line[len-1] = ' '
			line[len-2] = ' '
			continue
		}

		return (len)
	}
}
