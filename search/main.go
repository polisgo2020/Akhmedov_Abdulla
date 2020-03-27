package main

import (
	"encoding/json"
	"github.com/polisgo2020/Akhmedov_Abdulla/invertedIndex"
	"github.com/polisgo2020/Akhmedov_Abdulla/readFiles"
	"io/ioutil"
	"log"
	"os"
)

var invertedIn invertedIndex.Index

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments")
	}

	file, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(file, &invertedIn)
	if err != nil {
		log.Println(err)
		return
	}

	stopWords, err := readFiles.ReadStopWords(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	var searchPhrase []string
	for i := 3; i < len(os.Args); i++ {
		searchPhrase = append(searchPhrase, os.Args[i])
	}

	invertedIndex.PrintSortedList(searchPhrase, stopWords, invertedIn)
}

