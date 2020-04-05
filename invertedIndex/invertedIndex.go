package invertedIndex

import (
	"fmt"
	"github.com/caneroj1/stemmer"
	"github.com/polisgo2020/Akhmedov_Abdulla/readFiles"
	"math"
	"sort"
	"strings"
	"sync"
	"unicode"
)

type Index map[string]map[string][]int

var invertedIn Index

type safeIndex struct {
	invertedIndex Index
	mux           sync.Mutex
}

func (sin *safeIndex) addToken(token string, position int, stopWordsMap map[string]int, file string, wg *sync.WaitGroup) {
	token = strings.TrimFunc(token, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	token = stemmer.Stem(token)

	token = strings.ToLower(token)
	if _, ok := stopWordsMap[token]; !ok && len(token) != 0 {
		if sin.invertedIndex[token] == nil {
			sin.invertedIndex[token] = make(map[string][]int)
		}

		sin.mux.Lock()
		sin.invertedIndex[token][file] = append(sin.invertedIndex[token][file], position)
		sin.mux.Unlock()
	}
	wg.Done()
}

// GetInvertedIndex returns inverted index map that also stores position of each token in document
func GetInvertedIndex(flag bool, files []string, stopWordsFile string) (Index, error) {
	var (
		sin          = safeIndex{make(Index), sync.Mutex{}}
		filesMap     map[string]string
		err          error
		stopWordsMap = make(map[string]int)
		wg           sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		filesMap, err = readFiles.ReadFiles(flag, files)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if len(stopWordsFile) != 0 {
			stopWordsMap, err = readFiles.ReadStopWords(stopWordsFile)
		}
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return nil, err
	}

	for file, str := range filesMap {
		tokens := strings.Fields(str)
		for position, token := range tokens {
			wg.Add(1)
			sin.addToken(token, position, stopWordsMap, file, &wg)
		}
	}

	wg.Wait()
	// список всех файлов
	sin.invertedIndex[""] = make(map[string][]int)
	for file, _ := range filesMap {
		sin.invertedIndex[""][file] = append(sin.invertedIndex[""][file])
	}

	return sin.invertedIndex, nil
}

func PrintSortedList(searchPhrase []string, stopWords map[string]int, iIn Index) string {
	invertedIn = iIn
	var phrase []string

	var wg sync.WaitGroup
	for i := 0; i < len(searchPhrase); i++ {
		wg.Add(1)
		go func(j int) {
			if _, ok := stopWords[searchPhrase[j]]; ok {
				return
			}

			tmp := strings.ToLower(stemmer.Stem(searchPhrase[j]))
			if _, ok := invertedIn[tmp]; ok {
				phrase = append(phrase, tmp)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	searchPhrase = phrase
	var answer []float64
	answerMap := make(map[float64][]string)
	result := ""
	if len(searchPhrase) == 1 {
		if filesMap, ok := invertedIn[searchPhrase[0]]; ok {
			for file := range filesMap {
				tmp := float64(len(filesMap[file]))

				answer = append(answer, tmp)
				if answerMap[tmp] == nil {
					answerMap[tmp] = make([]string, 0, 0)
				}

				answerMap[tmp] = append(answerMap[tmp], file)
			}

			_ = sort.Reverse(sort.Float64Slice(answer))
			for _, v := range answer {
				files := answerMap[v]
				for _, file := range files {
					result = fmt.Sprintf("%s - %f\n", file, v)
				}
			}
		} else {
			fmt.Print("None of files contains this search-phrase")
		}
	} else {
		tmp := invertedIn[""]
		for file, _ := range tmp {
			res := getInfo(searchPhrase, file)
			answer = append(answer, res)

			if answerMap[res] == nil {
				answerMap[res] = make([]string, 0, 0)
			}

			answerMap[res] = append(answerMap[res], file)
		}

		sort.Float64s(answer)
		for _, v := range answer {
			files := answerMap[v]
			for _, file := range files {
				if v > 0.0000001 {
					result += fmt.Sprintf("%s - %f\n", file, v)
				}
			}
		}
	}

	return result
}

func getInfo(phrase []string, file string) float64 {
	distance, count := findMinWay(phrase, 0, file, 1)
	return float64(distance) / float64(count)
}

func findMinWay(phrase []string, index int, file string, count int) (int, int) {
	if index >= len(phrase) {
		return 0, count
	}

	curIndex := findFirstExistTokenIndex(phrase, index, file)
	nextIndex := findFirstExistTokenIndex(phrase, curIndex+1, file)
	if nextIndex == -1 || curIndex == -1 {
		return 0, count
	}

	curList := invertedIn[phrase[curIndex]][file]
	nextList := invertedIn[phrase[nextIndex]][file]

	min := math.MaxInt64
	var resCount int
	for _, v1 := range curList {
		var res int
		// TODO: кэшировать значения этого вызова в матрицу[index][count] -> сильно ускорит обход дерева
		res, resCount = findMinWay(phrase, nextIndex, file, count+1)
		for _, v2 := range nextList {
			delta := abs(v2 - v1)
			min = fMin(min, res+delta)
		}
	}

	return min, resCount
}

func fMin(a int, b int) int {
	if a <= b {
		return a
	}

	return b
}

func abs(a int) int {
	if a >= 0 {
		return a
	}

	return -a
}

func findFirstExistTokenIndex(phrase []string, index int, file string) int {
	for i := index; i < len(phrase); i++ {
		if _, ok := invertedIn[phrase[index]][file]; ok {
			return index
		}
	}

	return -1
}
