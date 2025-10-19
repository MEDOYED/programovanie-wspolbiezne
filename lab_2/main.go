package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fork represents a fork (widelec)
type Fork struct {
	sync.Mutex
}

// Philosopher represents a philosopher
type Philosopher struct {
	id         int
	leftFork   *Fork
	rightFork  *Fork
	thinkCount int
	eatCount   int
	mu         sync.Mutex // for safe access to counters
}

// think simulates thinking
func (p *Philosopher) think() {
	p.mu.Lock()
	p.thinkCount++
	p.mu.Unlock()

	// Random delay for thinking (10-50ms)
	time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)
}

// eat simulates eating
func (p *Philosopher) eat() {
	p.mu.Lock()
	p.eatCount++
	count := p.eatCount
	p.mu.Unlock()

	if count%50 == 0 {
		fmt.Printf("Filozof %d je (cykl %d)\n", p.id, count)
	}

	// Random delay for eating (10-50ms)
	time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)
}

// dine is the main philosopher loop
func (p *Philosopher) dine(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Think
			p.think()

			// Deadlock avoidance strategy:
			// Even philosophers pick right fork first, odd philosophers pick left fork first
			// This breaks the circular wait condition and prevents deadlock
			firstFork, secondFork := p.leftFork, p.rightFork

			// Determine fork acquisition order
			// Even philosophers reverse the order
			if p.id%2 == 0 {
				firstFork, secondFork = p.rightFork, p.leftFork
			}

			// Acquire first fork
			firstFork.Lock()

			// Check if context is done
			select {
			case <-ctx.Done():
				firstFork.Unlock()
				return
			default:
			}

			// Acquire second fork
			secondFork.Lock()

			// Eat (both forks acquired)
			p.eat()

			// Release forks
			secondFork.Unlock()
			firstFork.Unlock()
		}
	}
}

// getStatistics returns philosopher statistics
func (p *Philosopher) getStatistics() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.thinkCount, p.eatCount
}

func main() {
	const numPhilosophers = 5
	const duration = 30 * time.Second

	fmt.Println("=== Problem ucztujacych filozofow ===")
	fmt.Printf("Liczba filozofow: %d\n", numPhilosophers)
	fmt.Printf("Czas symulacji: %v\n", duration)
	fmt.Println("Start symulacji...")
	fmt.Println()

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Create context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Create forks
	forks := make([]*Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &Fork{}
	}

	// Create philosophers
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhilosophers],
		}
	}

	// WaitGroup for waiting for all goroutines to finish
	var wg sync.WaitGroup

	// Start goroutines for each philosopher
	startTime := time.Now()
	for _, philosopher := range philosophers {
		wg.Add(1)
		go philosopher.dine(ctx, &wg)
	}

	// Wait for all philosophers to finish
	wg.Wait()
	elapsed := time.Since(startTime)

	// Print statistics
	fmt.Println()
	fmt.Println("=== Koniec symulacji ===")
	fmt.Printf("Rzeczywisty czas: %.2f s\n", elapsed.Seconds())
	fmt.Println()
	fmt.Println("Statystyki:")
	fmt.Println("-------------------------------------------")

	totalThinks := 0
	totalEats := 0

	for _, philosopher := range philosophers {
		thinks, eats := philosopher.getStatistics()
		totalThinks += thinks
		totalEats += eats
		fmt.Printf("F%d Myslal: %d, Jadl: %d\n", philosopher.id, thinks, eats)
	}

	fmt.Println("-------------------------------------------")
	fmt.Printf("Razem: Myslenie: %d, Jedzenie: %d\n", totalThinks, totalEats)
	fmt.Println()

	// Explanation of deadlock and starvation prevention mechanism
	fmt.Println("=== Mechanizm zapobiegania deadlock i starvation ===")
	fmt.Println()
	fmt.Println("Deadlock-free:")
	fmt.Println("  Filozofowie parzysti (0,2,4) najpierw biora prawy widelec,")
	fmt.Println("  a filozofowie nieparzysti (1,3) najpierw biora lewy widelec.")
	fmt.Println("  To przerywa cykliczne oczekiwanie i uniemozliwia deadlock.")
	fmt.Println()
	fmt.Println("Starvation-free:")
	fmt.Println("  Uzywamy sync.Mutex, ktory gwarantuje sprawiedliwosc (fairness).")
	fmt.Println("  Kazdy watek czekajacy na mutex ma szanse go zdobyc.")
	fmt.Println("  Brak priorytetow oznacza, ze zaden filozof nie jest glodzony.")
}
