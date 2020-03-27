package readFiles

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func ReadFiles(flag bool, files []string) (map[string]string, error) {
	m := make(map[string]string)
	i := 0

	if flag {
		for _, v := range files {
			data, err := ioutil.ReadFile(v)
			if err != nil {
				log.Print(err, "Could not read file!")
				return nil, err
			}

			_, fileName := filepath.Split(v)
			m[fmt.Sprint(i)+"_"+fileName] = string(data)
			i++
		}
	} else {
		for _, v := range files {
			dir, err := ioutil.ReadDir(v)
			if err != nil {
				log.Print(err, "Could not read directory!")
				return nil, err
			}

			for _, file := range dir {
				data, err := ioutil.ReadFile(filepath.Join(v, file.Name()))
				if err != nil {
					log.Print(err, "Could not read file!")
					return nil, err
				}

				m[file.Name()] = string(data)
				i++
			}
		}
	}

	return m, nil
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

