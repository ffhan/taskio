package main

import (
	"fmt"
	"io/ioutil"
	"taskio"
	"taskio/container"
)

func main() {
	bytes, err := ioutil.ReadFile("/usr/bin/bash")
	taskio.Must(err)

	run, err := container.Run(bytes)
	taskio.Must(err)

	fmt.Println(string(run))
}
