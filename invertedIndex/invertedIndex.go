package invertedIndex

import (
	"fmt"
	"github.com/caneroj1/stemmer"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

type Index map[string]map[string][]int

func readFiles(flag bool, files []string) (map[string]string, error) {
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

// returns inverted index map that also stores position of each token in document
func GetInvertedIndex(flag bool, files []string, stopWordsFile string) (Index, error) {
	invertedIndex := make(map[string]map[string][]int)
	filesMap, err := readFiles(flag, files)
	if err != nil {
		return nil, err
	}

	stopWordsMap := make(map[string]int)
	if len(stopWordsFile) != 0 {
		stopWordsMap, err = ReadStopWords(stopWordsFile)
		if err != nil {
			return nil, err
		}
	}

	for file, str := range filesMap {
		tokens := strings.Fields(str)
		for position, token := range tokens {
			token = strings.TrimFunc(token, func(r rune) bool {
				return !unicode.IsLetter(r)
			})
			token = stemmer.Stem(token) // Насколько я понимаю эта либа сделана по этому алгоритму
			// https://tartarus.org/martin/PorterStemmer/def.txt

			token = strings.ToLower(token)
			if _, ok := stopWordsMap[token]; !ok && len(token) != 0 {
				if invertedIndex[token] == nil {
					invertedIndex[token] = make(map[string][]int)
				}

				invertedIndex[token][file] = append(invertedIndex[token][file], position)
			}
		}
	}

	// список всех файлов
	invertedIndex[""] = make(map[string][]int)
	for file, _ := range filesMap {
		invertedIndex[""][file] = append(invertedIndex[""][file])
	}

	return invertedIndex, nil
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
