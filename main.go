package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type MaxVowelsSubset struct {
	Text        string
	VowelsCount int32
}

var (
	chars          = "abcdefghijklmnopqrstuvwxyz"
	vowels         = "aeiou"
	str            string
	numOfWorkers   int
	wg             sync.WaitGroup
	mu             sync.Mutex
	maxVowelSubset MaxVowelsSubset
)

func init() {
	str = randomCharacters(500000)
	numOfWorkers = runtime.NumCPU()
	maxVowelSubset = MaxVowelsSubset{
		Text:        "Not found!",
		VowelsCount: 0,
	}
}

func main() {
	ch := make(chan string, 10)
	launchWorkers(ch, numOfWorkers)
	wg.Add(numOfWorkers)
	start := time.Now()

	findSubsets(str, 5, ch)
	close(ch)

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Process took", end)
	fmt.Println("Maximum Vowel Subset:", maxVowelSubset.Text)
}

func launchWorkers(queueChannel <-chan string, numOfWorkers int) {
	for i := 0; i < numOfWorkers; i++ {
		go readChannel(queueChannel, i)
	}
}

func readChannel(queueChannel <-chan string, workerId int) {
	for msg := range queueChannel {
		fmt.Printf("Worker %d got the subset: %s\n", workerId, msg)
		countVowels(msg)
	}
	wg.Done()
}

// findSubset returns the subset with specific length that contains more numbers of vowels
func findSubsets(str string, length int, ch chan<- string) {
	for i := 0; i < len(str)-length; i++ {
		j := i + length + 1
		ch <- str[i:j]
	}
}

func countVowels(str string) {
	var count int32
	for i := 0; i < len(str); i++ {
		if isVowel(str[i]) {
			count++
		}
	}
	mu.Lock()
	if count > maxVowelSubset.VowelsCount {
		maxVowelSubset.VowelsCount = count
		maxVowelSubset.Text = str
	}
	mu.Unlock()
}

func isVowel(ch byte) bool {
	for i := 0; i < 5; i++ {
		if ch == vowels[i] {
			return true
		}
	}
	return false
}

func randomCharacters(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = chars[rand.Intn(26)]
	}
	return string(b)
}
