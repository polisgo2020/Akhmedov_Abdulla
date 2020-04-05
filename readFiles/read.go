package readFiles

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type safeRead struct {
	filesMap map[string]string
	mux      sync.Mutex
}

func (sr *safeRead) addFile(filePath string, index *int) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Print(err, "Could not read file!")
		return err
	}

	_, fileName := filepath.Split(filePath)
	sr.mux.Lock()
	sr.filesMap[fileName] = string(data)
	sr.mux.Unlock()
	*index++

	return nil
}

func ReadFiles(flag bool, files []string) (map[string]string, error) {
	sr := safeRead{make(map[string]string), sync.Mutex{}}
	var index int

	index = 0
	wg := sync.WaitGroup{}
	if flag {
		for _, v := range files {
			data, err := ioutil.ReadFile(v)
			if err != nil {
				log.Print(err, "Could not read file!")
				return nil, err
			}

			_, fileName := filepath.Split(v)
			sr.filesMap[fmt.Sprint(index)+"_"+fileName] = string(data)
			index++
		}
	} else {
		for _, v := range files {
			dir, err := ioutil.ReadDir(v)
			if err != nil {
				log.Print(err, "Could not read directory!")
				return nil, err
			}

			for _, file := range dir {
				wg.Add(1)
				err = sr.addFile(filepath.Join(v, file.Name()), &index)
				if err != nil {
					return nil, err
				}
				wg.Done()
			}
		}
	}

	wg.Wait()
	return sr.filesMap, nil
}

func ReadStopWords(file string) (map[string]int, error) {
	m := make(map[string]int)

	str, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	words := strings.Fields(string(str))
	for _, word := range words {
		m[word] = 0
	}

	return m, nil
}
