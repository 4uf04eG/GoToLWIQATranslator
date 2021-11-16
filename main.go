package main

import (
	"./translator"
	"fmt"
	"io/ioutil"
)

func main() {
	code, err := ioutil.ReadFile("test5.notgo")
	if err != nil {
		fmt.Println("Could not open file")
		return
	}

	genCode := translator.Translate(string(code))

	fmt.Println(genCode)
}
