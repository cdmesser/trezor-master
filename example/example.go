package main

import (
	"fmt"
	"os"

	"github.com/cdmesser/trezor-master"
)

func main() {
	pin, err := trezor.GetPIN("Enter PIN:")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Success! Scrambled PIN input:", pin)
}
