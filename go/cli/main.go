package main

import (
	"fmt"
	"os"

	"lib-post-interchange/libale"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("cli: Please provide a file path as first argument")
		return
	}
	inputFile := fmt.Sprintf("%s", os.Args[1])
	fmt.Printf("cli: Input file:", inputFile)
	ale, err := libale.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ALE object: %s", ale)

}
