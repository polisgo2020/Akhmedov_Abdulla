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
	mux sync.Mutex
}

// GetInvertedIndex returns inverted index map that also stores position of each token in document
func GetInvertedIndex(flag bool, files []string, stopWordsFile string) (Index, error) {
	var sin = safeIndex{make(Index), sync.Mutex{}}

	filesMap, err := readFiles.ReadFiles(flag, files)
	if err != nil {
		return nil, err
	}

	stopWordsMap := make(map[string]int)
	if len(stopWordsFile) != 0 {
		stopWordsMap, err = readFiles.ReadStopWords(stopWordsFile)
		if err != nil {
			return nil, err
		}
	}

	var wg sync.WaitGroup
	for file, str := range filesMap {
		tokens := strings.Fields(str)
		for position, token := range tokens {
			wg.Add(1)
			go func(goToken string, goPosition int) {
				sin.mux.Lock()
				defer sin.mux.Unlock()
				defer wg.Done()
				goToken = strings.TrimFunc(goToken, func(r rune) bool {
					return !unicode.IsLetter(r)
				})
				goToken = stemmer.Stem(goToken) // Насколько я понимаю эта либа сделана по этому алгоритму
				// https://tartarus.org/martin/PorterStemmer/def.txt

				goToken = strings.ToLower(goToken)
				if _, ok := stopWordsMap[goToken]; !ok && len(goToken) != 0 {
					if sin.invertedIndex[goToken] == nil {
						sin.invertedIndex[goToken] = make(map[string][]int)
					}

					sin.invertedIndex[goToken][fie] = append(sin.invertedIndex[goToken][file], goPosition)
				}
			}(token, position)
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

func PrintSortedList(searchPhrase []string, stopWords map[string]int, iIn Index) {
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
					fmt.Printf("%s - %f\n", file, v)
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
					fmt.Printf("%s - %f\n", file, v)
				}
			}
		}
	}
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
