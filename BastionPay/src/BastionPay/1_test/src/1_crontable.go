package main

import (
	"fmt"
	"os/exec"
)

func main()  {
	var (
		cmd *exec.Cmd
		out []byte
		err error
	)

	cmd = exec.Command("/bin/bash", "-c","/usr/local/bin/go run /Users/sywu/GoProject/test/src/sort.go")
	//cmd = exec.Command("/bin/bash", "-c","/usr/local/bin/go run /Users/sywu/GoProject/test/src/sort.go")

    if out, err = cmd.CombinedOutput(); err != nil {
		fmt.Println("eee",err)
		return
	}

	fmt.Println(string(out))

}
