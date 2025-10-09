package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DataBuffer represents the shared resource between goroutines
type DataBuffer struct {
	buffer       chan int
	statusAccess sync.RWMutex

	recentWriter    string
	recentConsumer  string
	pendingValue    *int
	operationCount  int
}

// Configuration constants
const (
	programDuration = 30 * time.Second
	studentID       = "121261"
	monitorInterval = 1000 * time.Millisecond
)

func main() {
	displayHeader()

	dataBuffer := initializeDataBuffer()

	programCtx, cancelProgram := context.WithTimeout(context.Background(), programDuration)
	defer cancelProgram()

	workersGroup := &sync.WaitGroup{}

	// Launch all workers
	workersGroup.Add(4)
	go firstDataProducer(programCtx, dataBuffer, workersGroup)
	go secondDataProducer(programCtx, dataBuffer, workersGroup)
	go dataAggregator(programCtx, dataBuffer, workersGroup)
	go systemMonitor(programCtx, dataBuffer, workersGroup)

	workersGroup.Wait()

	displayFooter()
}

func displayHeader() {
	fmt.Println("Starting concurrent program...")
	fmt.Printf("Duration: %v\n", programDuration)
}

func displayFooter() {
	fmt.Println("\nProgram finished")
}

func initializeDataBuffer() *DataBuffer {
	return &DataBuffer{
		buffer: make(chan int, 1),
	}
}

// Store writes a value to the buffer with proper synchronization
func (db *DataBuffer) Store(val int, producerID string) error {
	select {
	case db.buffer <- val:
		db.statusAccess.Lock()
		db.recentWriter = producerID
		db.pendingValue = &val
		db.operationCount++
		db.statusAccess.Unlock()

		fmt.Printf("[%s] wrote value: %d\n", producerID, val)
		return nil
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("buffer full")
	}
}

// Retrieve reads a value from the buffer with proper synchronization
func (db *DataBuffer) Retrieve(consumerID string) (int, error) {
	select {
	case val := <-db.buffer:
		db.statusAccess.Lock()
		db.recentConsumer = consumerID
		db.pendingValue = nil
		db.statusAccess.Unlock()

		return val, nil
	case <-time.After(50 * time.Millisecond):
		return 0, fmt.Errorf("buffer empty")
	}
}

// GetSnapshot returns the current state of the buffer
func (db *DataBuffer) GetSnapshot() (currentValue *int, lastProducer, lastConsumer string, operations int) {
	db.statusAccess.RLock()
	defer db.statusAccess.RUnlock()

	return db.pendingValue, db.recentWriter, db.recentConsumer, db.operationCount
}

// firstDataProducer generates numbers in range [21-37]
func firstDataProducer(ctx context.Context, buffer *DataBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	workerID := fmt.Sprintf("%s#Writer#1", studentID)
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Printf("[%s] started\n", workerID)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] stopped\n", workerID)
			return
		default:
			// Generate random number in specified range
			generatedNum := randomizer.Intn(17) + 21

			// Try to store the value
			buffer.Store(generatedNum, workerID)

			// Random sleep between operations
			pauseDuration := time.Duration(randomizer.Intn(800)+400) * time.Millisecond
			time.Sleep(pauseDuration)
		}
	}
}

// secondDataProducer generates numbers in range [1337-4200]
func secondDataProducer(ctx context.Context, buffer *DataBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	workerID := fmt.Sprintf("%s#Writer#2", studentID)
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano() + 12345))

	fmt.Printf("[%s] started\n", workerID)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] stopped\n", workerID)
			return
		default:
			// Generate random number in specified range
			generatedNum := randomizer.Intn(2864) + 1337

			// Try to store the value
			buffer.Store(generatedNum, workerID)

			// Random sleep between operations
			pauseDuration := time.Duration(randomizer.Intn(800)+400) * time.Millisecond
			time.Sleep(pauseDuration)
		}
	}
}

// dataAggregator consumes values and maintains a running sum
func dataAggregator(ctx context.Context, buffer *DataBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	workerID := fmt.Sprintf("%s#Reader#1", studentID)
	runningTotal := 0
	itemsProcessed := 0

	fmt.Printf("[%s] started\n", workerID)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] stopped. Total sum: %d\n", workerID, runningTotal)
			return
		default:
			// Attempt to retrieve value
			retrievedValue, err := buffer.Retrieve(workerID)

			if err == nil {
				runningTotal += retrievedValue
				itemsProcessed++

				fmt.Printf("[%s] read value: %d, current sum: %d\n",
					workerID, retrievedValue, runningTotal)
			}

			// Small delay between read attempts
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// systemMonitor periodically reports the system state
func systemMonitor(ctx context.Context, buffer *DataBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	workerID := fmt.Sprintf("%s#Monitor#1", studentID)
	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	fmt.Printf("[%s] started\n", workerID)

	reportCount := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] stopped\n", workerID)
			return
		case <-ticker.C:
			reportCount++
			generateStatusReport(workerID, buffer, reportCount)
		}
	}
}

func generateStatusReport(monitorID string, buffer *DataBuffer, reportNum int) {
	value, producer, consumer, operations := buffer.GetSnapshot()

	fmt.Printf("\n--- [%s] Status report ---\n", monitorID)

	if value != nil {
		fmt.Printf("Buffer value: %d\n", *value)
	} else {
		fmt.Printf("Buffer value: empty\n")
	}

	if producer != "" {
		fmt.Printf("Last writer: %s\n", producer)
	} else {
		fmt.Printf("Last writer: none\n")
	}

	if consumer != "" {
		fmt.Printf("Last reader: %s\n", consumer)
	} else {
		fmt.Printf("Last reader: none\n")
	}

	fmt.Printf("Operations: %d\n", operations)
	fmt.Println("All threads: active")
	fmt.Println()
}
