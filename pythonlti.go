package main

import (
	"fmt"
	"os/exec"
)

const (
	PathToPythonRequest = "./python/transferGrade.py"
)

func sendRequest() {
	cmd := exec.Command(Python, PathToPythonRequest)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(outb))
	return
}
