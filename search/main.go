package main

import (
	"encoding/json"
	"github.com/polisgo2020/Akhmedov_Abdulla/invertedIndex"
	"github.com/polisgo2020/Akhmedov_Abdulla/readFiles"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var invertedIn invertedIndex.Index

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()

	var stopWords map[string]int
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		stopWords, err = readFiles.ReadStopWords(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
	}()

	var searchPhrase []string
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 3; i < len(os.Args); i++ {
			searchPhrase = append(searchPhrase, os.Args[i])
		}
	}()

	wg.Wait()
	invertedIndex.PrintSortedList(searchPhrase, stopWords, invertedIn)
}

