package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var inputFilePath *string = flag.String("f", "", "Json input file")
var outputFilePath *string = flag.String("o", "", "Go output file")
var outputPackageName *string = flag.String("p", "", "Go package name for the output file")

func main() {
	flag.Parse()
	if *inputFilePath == "" {
		flag.PrintDefaults()
		log.Fatal("Missing input file (-f).")
	}
	if *outputFilePath == "" {
		flag.PrintDefaults()
		log.Fatal("Missing output file (-o).")
	}
	if *outputPackageName == "" {
		flag.PrintDefaults()
		log.Fatal("Missing output package name (-p).")
	}

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		log.Fatalf("Could not open file: %v\n", err)
	}
	defer inputFile.Close()

	idList := make([]string, 0)
	decoder := json.NewDecoder(inputFile)
	err = decoder.Decode(&idList)
	if err != nil {
		log.Fatalf("Could not decode JSON: %v\n", err)
	}

	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		log.Fatalf("Could not create output file: %v\n", err)
	}
	defer outputFile.Close()

	var writeErr error
	write := func(s string) {
		if writeErr != nil {
			return
		}
		_, writeErr = outputFile.WriteString(s)
	}

	write("package " + *outputPackageName + "\n\n")
	write("const (\n")

	for _, line := range idList {
		id := strings.Replace(line, ".", "_", -1)
		write(fmt.Sprintf("\t%s = \"%s\"\n", id, line))
	}

	write(")\n")
}
