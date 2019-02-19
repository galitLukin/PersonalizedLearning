package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	Python             = "python3"
	PathToPythonScript = "./python/transferGrade.py"
)

func sendRequest() {
	cmd := exec.Command(Python, PathToPythonScript)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(outb))
	err = json.Unmarshal(outb, &q)
	if err != nil {
		fmt.Println(err)
	}
	return
}
