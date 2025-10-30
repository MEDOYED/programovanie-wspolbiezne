package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// FileResult przechowuje wynik zliczania slow dla jednego pliku
type FileResult struct {
	FileName  string
	WordCount int
	WorkerID  int
}

// WorkerPool zarzadza pula workerow do rownoleglego przetwarzania
type WorkerPool struct {
	numWorkers int
	jobs       chan string
	results    chan FileResult
	wg         sync.WaitGroup
	workerID   int32
}

// NewWorkerPool tworzy nowa pule workerow
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan string, 100),
		results:    make(chan FileResult, 100),
	}
}

// Start uruchamia workery
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker przetwarza pliki z kanalu jobs
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	workerID := int(atomic.AddInt32(&wp.workerID, 1))

	for filePath := range wp.jobs {
		wordCount := countWordsInFile(filePath)
		wp.results <- FileResult{
			FileName:  filepath.Base(filePath),
			WordCount: wordCount,
			WorkerID:  workerID,
		}
	}
}

// AddJob dodaje plik do kolejki przetwarzania
func (wp *WorkerPool) AddJob(filePath string) {
	wp.jobs <- filePath
}

// Close zamyka kanal jobs i czeka na zakonczenie wszystkich workerow
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// countWordsInFile liczy liczbe slow w pliku
func countWordsInFile(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Blad otwierania pliku %s: %v", filePath, err)
		return 0
	}
	defer file.Close()

	wordCount := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Blad czytania pliku %s: %v", filePath, err)
		return 0
	}

	return wordCount
}

// findTextFiles rekurencyjnie znajduje wszystkie pliki .txt w katalogu
func findTextFiles(rootDir string) ([]string, error) {
	var textFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".txt") {
			textFiles = append(textFiles, path)
		}
		return nil
	})

	return textFiles, err
}

func main() {
	// Sprawdzenie argumentow linii polecen
	if len(os.Args) < 2 {
		fmt.Println("Uzycie: go run main.go <sciezka_do_folderu>")
		os.Exit(1)
	}

	rootDir := os.Args[1]

	// Sprawdzenie czy katalog istnieje
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		log.Fatalf("Katalog nie istnieje: %s", rootDir)
	}

	fmt.Println("Szukanie plikow tekstowych...")
	textFiles, err := findTextFiles(rootDir)
	if err != nil {
		log.Fatalf("Blad wyszukiwania plikow: %v", err)
	}

	if len(textFiles) == 0 {
		fmt.Println("Nie znaleziono zadnych plikow .txt")
		return
	}

	fmt.Printf("Znaleziono %d plikow. Rozpoczynam przetwarzanie...\n", len(textFiles))

	// Tworzenie puli workerow
	numWorkers := 8
	pool := NewWorkerPool(numWorkers)
	pool.Start()

	// Dodawanie wszystkich plikow do kolejki
	go func() {
		for _, filePath := range textFiles {
			pool.AddJob(filePath)
		}
		pool.Close()
	}()

	// Zbieranie wynikow
	var results []FileResult
	totalWords := 0

	for result := range pool.results {
		results = append(results, result)
		totalWords += result.WordCount
	}

	// Tworzenie pliku logu
	logFile, err := os.Create("word_count_log.txt")
	if err != nil {
		log.Fatalf("Blad tworzenia pliku logu: %v", err)
	}
	defer logFile.Close()

	// Zapisywanie wynikow do konsoli i pliku logu
	fmt.Println("\nWyniki zliczania slow:")
	for _, result := range results {
		line := fmt.Sprintf("Worker-%d -> %s: %d slow", result.WorkerID, result.FileName, result.WordCount)
		fmt.Println(line)
		logFile.WriteString(line + "\n")
	}

	// Zapisywanie lacznej liczby slow
	totalLine := fmt.Sprintf("\nLaczna liczba slow: %d", totalWords)
	fmt.Println(totalLine)
	logFile.WriteString(totalLine + "\n")

	fmt.Println("\nLog zapisany do pliku: word_count_log.txt")
}
