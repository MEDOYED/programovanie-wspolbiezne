package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// TaskResult stores the result of a single goroutine task
type TaskResult struct {
	TaskID      int
	ExecutionMS int64
	Message     string
}

// executeTask simulates a lightweight task with a small delay
func executeTask(taskID int, results chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	// Simulate some work with a small delay
	time.Sleep(1 * time.Millisecond)

	// Simulate computation
	result := fmt.Sprintf("Task %d completed", taskID)

	executionTime := time.Since(startTime).Milliseconds()

	results <- TaskResult{
		TaskID:      taskID,
		ExecutionMS: executionTime,
		Message:     result,
	}
}

// runGoroutineTest creates and executes N goroutines
func runGoroutineTest(numTasks int) (time.Duration, []TaskResult) {
	fmt.Printf("\n=== Starting Goroutine Test with %d tasks ===\n", numTasks)

	results := make(chan TaskResult, numTasks)
	var wg sync.WaitGroup

	// Measure memory before
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	startTime := time.Now()

	// Launch all goroutines
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go executeTask(i, results, &wg)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)

	// Measure memory after
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)

	// Collect results
	var taskResults []TaskResult
	for result := range results {
		taskResults = append(taskResults, result)
	}

	// Print statistics
	fmt.Printf("Total execution time: %v\n", totalTime)
	fmt.Printf("Tasks completed: %d\n", len(taskResults))
	fmt.Printf("Average time per task: %.2f ms\n", float64(totalTime.Milliseconds())/float64(numTasks))
	fmt.Printf("Memory allocated: %.2f MB\n", float64(memStatsAfter.Alloc-memStatsBefore.Alloc)/1024/1024)
	fmt.Printf("Number of OS threads used: %d\n", runtime.NumGoroutine())

	return totalTime, taskResults
}

// runHeavyThreadTest simulates behavior closer to OS threads
// by limiting parallelism with GOMAXPROCS
func runHeavyThreadTest(numTasks int, maxThreads int) (time.Duration, []TaskResult) {
	fmt.Printf("\n=== Starting Heavy Thread Test with %d tasks (limited to %d threads) ===\n", numTasks, maxThreads)

	// Save current GOMAXPROCS value
	oldMaxProcs := runtime.GOMAXPROCS(maxThreads)
	defer runtime.GOMAXPROCS(oldMaxProcs)

	results := make(chan TaskResult, numTasks)
	var wg sync.WaitGroup

	// Measure memory before
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	startTime := time.Now()

	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, maxThreads)

	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(taskID int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			startTime := time.Now()

			// Simulate some work with a small delay
			time.Sleep(1 * time.Millisecond)

			// Simulate computation
			result := fmt.Sprintf("Task %d completed", taskID)

			executionTime := time.Since(startTime).Milliseconds()

			results <- TaskResult{
				TaskID:      taskID,
				ExecutionMS: executionTime,
				Message:     result,
			}
		}(i)
	}

	// Wait for all tasks to complete
	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)

	// Measure memory after
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)

	// Collect results
	var taskResults []TaskResult
	for result := range results {
		taskResults = append(taskResults, result)
	}

	// Print statistics
	fmt.Printf("Total execution time: %v\n", totalTime)
	fmt.Printf("Tasks completed: %d\n", len(taskResults))
	fmt.Printf("Average time per task: %.2f ms\n", float64(totalTime.Milliseconds())/float64(numTasks))
	fmt.Printf("Memory allocated: %.2f MB\n", float64(memStatsAfter.Alloc-memStatsBefore.Alloc)/1024/1024)

	return totalTime, taskResults
}

// printComparison displays the comparison between both approaches
func printComparison(goroutineTime, heavyThreadTime time.Duration, numTasks int) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("                         COMPARISON RESULTS")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Number of tasks: %d\n\n", numTasks)

	fmt.Printf("Goroutines (lightweight):\n")
	fmt.Printf("  Total time: %v\n", goroutineTime)
	fmt.Printf("  Performance: %.2f tasks/second\n\n", float64(numTasks)/goroutineTime.Seconds())

	fmt.Printf("Limited threads (heavyweight):\n")
	fmt.Printf("  Total time: %v\n", heavyThreadTime)
	fmt.Printf("  Performance: %.2f tasks/second\n\n", float64(numTasks)/heavyThreadTime.Seconds())

	speedup := float64(heavyThreadTime) / float64(goroutineTime)
	fmt.Printf("Speedup: %.2fx faster with goroutines\n", speedup)

	fmt.Println(strings.Repeat("=", 70))
}

func main() {
	fmt.Println("Go Goroutines vs Traditional Threading Model")
	fmt.Println("Demonstration of lightweight concurrency (similar to Project Loom)")
	fmt.Println(strings.Repeat("=", 70))

	// Configuration
	numTasks := 10000
	numHeavyThreads := 50 // Simulate limited OS threads

	// Show system info
	fmt.Printf("System info:\n")
	fmt.Printf("  CPU cores: %d\n", runtime.NumCPU())
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	// Run lightweight goroutine test (analogous to Java Virtual Threads)
	goroutineTime, _ := runGoroutineTest(numTasks)

	// Allow garbage collection between tests
	runtime.GC()
	time.Sleep(1 * time.Second)

	// Run limited thread test (analogous to traditional Java threads)
	heavyThreadTime, _ := runHeavyThreadTest(numTasks, numHeavyThreads)

	// Print comparison
	printComparison(goroutineTime, heavyThreadTime, numTasks)

	fmt.Println("\nConclusion:")
	fmt.Println("Goroutines in Go provide lightweight concurrency similar to Java's Virtual Threads")
	fmt.Println("They are much more efficient than traditional OS threads for I/O-bound tasks")
}
