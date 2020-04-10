package main

import (
	"encoding/json"
	"flag"
	"github.com/polisgo2020/Akhmedov_Abdulla/invertedIndex"
	"github.com/polisgo2020/Akhmedov_Abdulla/readFiles"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	invertedIndexFile = "outputJSON.txt"
	stopWordsFile     = "stopWords.txt"
	invertedIndexMap  invertedIndex.Index
	stopWords         map[string]int
	wg                sync.WaitGroup
	errChannel        = make(chan error, 2)

	handler = http.NewServeMux()
	server  = http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
)

func init() {
	if len(os.Args) == 1 {
		log.Fatal("Not enough arguments: network interface is not provided")
	}

	invertedIndexFlag := flag.String("index", "", "Contains inverted index file path")
	stopWordsFlag := flag.String("sw", "", "Contains stop-words file path")
	flag.Parse()
	if *invertedIndexFlag != "" {
		invertedIndexFile = *invertedIndexFlag
	}
	if *stopWordsFlag != "" {
		stopWordsFile = *stopWordsFlag
	}

	server.Addr = os.Args[len(os.Args)-1]

	wg.Add(1)
	go func() {
		var mErr error
		stopWords, mErr = readFiles.ReadStopWords(stopWordsFile)
		if mErr != nil {
			errChannel <- mErr
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var mErr error
		file, mErr := ioutil.ReadFile(invertedIndexFile)
		if mErr != nil {
			errChannel <- mErr
		}
		mErr = json.Unmarshal(file, &invertedIndexMap)
		if mErr != nil {
			errChannel <- mErr
		}
		wg.Done()
	}()

	wg.Wait()
	close(errChannel)

	if err, ok := <- errChannel; ok {
		log.Fatal(err)
		return
	}
}

func main() {
	log.Printf("Server starts at %s", server.Addr)
	handler.HandleFunc("/search", logger(search))
	log.Fatal(server.ListenAndServe())
}

func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	var searchPhrase []string
	searchPhrase = strings.Fields(string(rBody))

	result := invertedIndex.PrintSortedList(searchPhrase, stopWords, invertedIndexMap)

	_, err = w.Write([]byte(result))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		log.Printf("Method [%s] connection from [%s]", r.Method, r.RemoteAddr)

		next.ServeHTTP(w, r)
	}
}
