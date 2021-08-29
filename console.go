package main

import (
	"fmt"
	"strings"
)

var NoConsole bool = false

func ReadLine() string {
	var s string
	if NoConsole {
		s = ""
	} else {
		fmt.Scanln(&s)
	}
	return s
}

func Prompt(msg string) string {
	Print(msg)
	return ReadLine()
}

func Print(msg string) {
	println(msg)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintError(err error) bool {
	if err != nil {
		Print(err.Error())
		return true
	}
	return false
}

func PromptOptions(msg string, options map[string]string) string {
	Print(msg)
	for o, m := range options {
		Print("\t[" + o + "] = " + m)
	}

	var r string = ""
	if NoConsole {
		// Select First key
		for o := range options {
			r = o
			break
		}
	} else {
		r = strings.TrimSpace(strings.ToLower(Prompt("Enter [value] : ")))
	}

	_, ok := options[r]
	if ok {
		return r
	} else {
		Print("Sorry, I didn't get that. Please enter the [option] you want ")
		return PromptOptions(msg, options)
	}

}
