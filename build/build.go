package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/polisgo2020/Akhmedov_Abdulla/invertedIndex"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// os.Args[0] = path to executable file so len(os.Args) equals 1 when no arguments provided
	if len(os.Args) == 1 {
		log.Fatal("There are no arguments!")
	}

	//if false directory provided
	sin := flag.Bool("s", false, "True if either single file or sequence of files provided")
	stopWords := flag.String("sw", "", "True if stop-words file provided")
	flag.Parse()

	var directories []string
	if *sin && len(*stopWords) == 0 {
		directories = os.Args[2:] // os.Args[1] = sin
	} else if *sin {
		directories = os.Args[3:]
	} else if len(*stopWords) != 0 {
		directories = os.Args[2:]
	} else {
		directories = os.Args[1:]
	}
	invertedIn, err := invertedIndex.GetInvertedIndex(*sin, directories, *stopWords)
	jsonInverted, err := json.Marshal(invertedIn)
	if err != nil {
		log.Fatal(err, "Could not Marshall!")
	}

	const PERMISSION = 0444 // read only
	if err := ioutil.WriteFile("outputJSON.txt", jsonInverted, PERMISSION); err != nil {
		fmt.Println(err)
		return
	}
}
