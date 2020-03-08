package main

import (
	"flag"
	"fmt"
	"github.com/polisgo2020/Akhmedov_Abdulla/inverted_index"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// os.Args[0] = path to executable file so len(os.Args) equals 1 when no arguments provided
	if len(os.Args) == 1 {
		log.Fatal("There are no arguments!")
	}

	// if false directory provided
	sin := flag.Bool("s", false, "True if either single file or sequence of files provided")
	flag.Parse()

	var directories []string
	if *sin {
		directories = os.Args[2:] // os.Args[1] = sin
	} else {
		directories = os.Args[1:]
	}

	tmp := inverted_index.GetInvertedIndex(*sin, directories)
	var str string
	for token, documents := range tmp {
		str += fmt.Sprintf("\"%s\": ", token)
		i := 0
		for document, positions := range documents {
			str += fmt.Sprintf("%s{", document)
			for i, position := range positions {
				if i == len(positions)-1 {
					str += fmt.Sprintf("%d", position)
				} else {
					str += fmt.Sprintf("%d, ", position)
				}
			}
			if i == len(documents)-1 {
				str += fmt.Sprintf("} \n")
			} else {
				str += fmt.Sprintf("} | ")
			}
			i++
		}
	}

	const PERMISSION = 0444 // read only
	if err := ioutil.WriteFile("output.txt", []byte(str), PERMISSION); err != nil {
		fmt.Println(err)
		return
	}
}
