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
	id int
}

// Waiter controls access to the dining room (Arbiter pattern)
// This prevents deadlock by limiting concurrent diners
type Waiter struct {
	semaphore chan struct{}
}

// NewWaiter creates a new waiter with capacity for maxDiners
func NewWaiter(maxDiners int) *Waiter {
	return &Waiter{
		semaphore: make(chan struct{}, maxDiners),
	}
}

// RequestPermission asks waiter for permission to dine
func (w *Waiter) RequestPermission(ctx context.Context) bool {
	select {
	case w.semaphore <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

// ReleasePermission returns permission to the waiter
func (w *Waiter) ReleasePermission() {
	<-w.semaphore
}

// Philosopher represents a philosopher
type Philosopher struct {
	id         int
	leftFork   *Fork
	rightFork  *Fork
	waiter     *Waiter
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

// dine is the main philosopher loop with Waiter (Arbiter) pattern
func (p *Philosopher) dine(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Think
			p.think()

			// Request permission from waiter to enter dining room
			// Waiter allows max N-1 philosophers to dine simultaneously
			// This prevents deadlock - there's always one philosopher not competing for forks
			if !p.waiter.RequestPermission(ctx) {
				return
			}

			// Now we have permission, acquire forks
			// Since waiter limits concurrent diners to N-1, deadlock is impossible
			p.leftFork.Lock()

			// Check if context is done before acquiring second fork
			select {
			case <-ctx.Done():
				p.leftFork.Unlock()
				p.waiter.ReleasePermission()
				return
			default:
			}

			p.rightFork.Lock()

			// Eat (both forks acquired)
			p.eat()

			// Release forks
			p.rightFork.Unlock()
			p.leftFork.Unlock()

			// Release permission back to waiter
			p.waiter.ReleasePermission()
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

	fmt.Println("=== Problem ucztujacych filozofow - Rozwiazanie z Arbiterem ===")
	fmt.Printf("Liczba filozofow: %d\n", numPhilosophers)
	fmt.Printf("Czas symulacji: %v\n", duration)
	fmt.Printf("Strategia: Waiter/Arbiter (maksymalnie %d filozofow jednoczesnie)\n", numPhilosophers-1)
	fmt.Println("Start symulacji...")
	fmt.Println()

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Create context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Create waiter that allows max N-1 philosophers to dine simultaneously
	// This is the key to preventing deadlock in the Arbiter solution
	waiter := NewWaiter(numPhilosophers - 1)

	// Create forks
	forks := make([]*Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &Fork{id: i}
	}

	// Create philosophers
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhilosophers],
			waiter:    waiter,
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
	fmt.Println("Rozwiazanie z Arbiterem (Waiter Pattern):")
	fmt.Println("  Kelner (Arbiter) kontroluje dostep do jadalni.")
	fmt.Printf("  Maksymalnie N-1 filozofow (%d z %d) moze jednoczesnie próbować jesc.\n", numPhilosophers-1, numPhilosophers)
	fmt.Println()
	fmt.Println("Deadlock-free:")
	fmt.Println("  Poniewaz maksymalnie 4 filozofow moze byc w jadalni,")
	fmt.Println("  zawsze jest przynajmniej jeden filozof ktory nie konkuruje o widelce.")
	fmt.Println("  To uniemozliwia sytuacje, w ktorej wszyscy trzymaja jeden widelec")
	fmt.Println("  i czekaja na drugi - czyli deadlock jest niemozliwy.")
	fmt.Println()
	fmt.Println("Starvation-free:")
	fmt.Println("  Semafor (buffered channel) w Go dziala w kolejnosci FIFO.")
	fmt.Println("  Kazdy filozof czekajacy na pozwolenie od kelnera dostanie je")
	fmt.Println("  po kolei - nikt nie jest pomijany.")
	fmt.Println()
	fmt.Println("Zalety tego podejscia:")
	fmt.Println("  + Proste w implementacji i zrozumieniu")
	fmt.Println("  + Centralna kontrola dostepu do zasobow")
	fmt.Println("  + Latwe do modyfikacji (mozna zmienic limit)")
	fmt.Println()
	fmt.Println("Wady tego podejscia:")
	fmt.Println("  - Jeden filozof zawsze czeka (mniejsza rownolegliosc)")
	fmt.Println("  - Potrzebna dodatkowa struktura (Waiter)")
}
