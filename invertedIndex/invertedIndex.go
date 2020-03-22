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

				m[strings.Join([]string{fmt.Sprintf("%d", i), "_", file.Name()}, "")] = string(data)
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
		tmp, err := readFiles(true, []string{stopWordsFile})

		for _, str := range tmp {
			words := strings.Fields(str)
			for _, word := range words {
				stopWordsMap[word] = 0
			}
		}
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
			if _, ok := stopWordsMap[token]; !ok {
				if invertedIndex[token] == nil {
					invertedIndex[token] = make(map[string][]int)
				}

				invertedIndex[token][file] = append(invertedIndex[token][file], position)
			}
		}
	}

	return invertedIndex, nil
}
