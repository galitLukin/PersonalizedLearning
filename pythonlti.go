package main

import (
	"fmt"
	"os/exec"
)

const (
	PathToPythonRequest = "./python/transferGrade.py"
)

func sendRequest(url, sourcedId string) {
	cmd := exec.Command(Python, PathToPythonRequest, url, sourcedId)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(outb))
	return
}
