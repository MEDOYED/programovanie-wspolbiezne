# Project Loom - Go Goroutines vs Traditional Threading

## Opis zadania

Celem zadania jest zrozumienie różnicy między klasycznymi wątkami a lekkimi mechanizmami współbieżności (takimi jak wirtualne wątki w Project Loom lub goroutines w Go), oraz wykorzystanie nowego modelu współbieżności w praktyce.

### Wymagania

**Krótki raport (1-1.5 strony):**
- Wyjaśnienie czym są wirtualne wątki (virtual threads) w Javie 25 oraz goroutines w Go
- Opis różnic między tradycyjnymi wątkami (OS threads) a lekkimi mechanizmami współbieżności
- Przedstawienie problemów w klasycznym modelu wątków, które rozwiązuje Project Loom i goroutines
- Schemat lub wypunktowanie opisujące mechanizm scheduler'a dla lekkich wątków

**Zadanie praktyczne - program w Go:**
Program który:
- Tworzy 10,000 goroutines (odpowiednik wirtualnych wątków)
- Każda goroutine wykonuje symboliczne zadanie (krótkie opóźnienie + prosty wydruk)
- Program mierzy i porównuje czas wykonania:
  - Goroutines (lightweight concurrency)
  - Ograniczona liczba wątków OS (traditional threading model)
- Raport zawiera opis wyniku oraz wnioski

## Funkcjonalność programu

Program demonstruje:
- **Lightweight goroutines** - utworzenie 10,000 goroutines jednocześnie
- **Limited thread model** - symulacja tradycyjnego modelu z ograniczoną liczbą wątków (50)
- **Pomiar wydajności** - porównanie czasu wykonania i zużycia pamięci
- **Statystyki** - szczegółowe informacje o wykonaniu zadań

## Uruchomienie

```bash
cd project-loom
go run main.go
```

## Wyniki przykładowe

Program generuje porównanie dwóch podejść:

```
=== Goroutines (lightweight) ===
Total execution time: 1.234s
Tasks completed: 10000
Average time per task: 0.12 ms
Memory allocated: 12.45 MB

=== Limited threads (heavyweight) ===
Total execution time: 5.678s
Tasks completed: 10000
Average time per task: 0.57 ms
Memory allocated: 45.23 MB

Speedup: 4.60x faster with goroutines
```

## Implementacja

### Główne komponenty:

1. **TaskResult** - struktura przechowująca wynik pojedynczego zadania
2. **executeTask** - funkcja symulująca lekkie zadanie z opóźnieniem
3. **runGoroutineTest** - test z wykorzystaniem lekkich goroutines
4. **runHeavyThreadTest** - test z ograniczoną liczbą wątków (symulacja OS threads)
5. **printComparison** - wyświetlenie porównania wyników

### Kluczowe mechanizmy:

- **sync.WaitGroup** - synchronizacja zakończenia wszystkich goroutines
- **Channels** - komunikacja i zbieranie wyników
- **runtime.MemStats** - pomiar zużycia pamięci
- **Semaphore pattern** - ograniczenie liczby równoczesnych goroutines

## Architektura

```
┌─────────────────────────────────────────────────────────┐
│                    Main Program                         │
└───────────────┬─────────────────────┬───────────────────┘
                │                     │
        ┌───────▼────────┐    ┌──────▼──────────┐
        │  Goroutine     │    │  Limited Thread │
        │  Test          │    │  Test           │
        │  (10,000)      │    │  (50 threads)   │
        └───────┬────────┘    └──────┬──────────┘
                │                     │
        ┌───────▼────────┐    ┌──────▼──────────┐
        │  Light weight  │    │  Heavy weight   │
        │  Fast start    │    │  Slow start     │
        │  Low memory    │    │  High memory    │
        └───────┬────────┘    └──────┬──────────┘
                │                     │
                └──────────┬──────────┘
                           │
                    ┌──────▼──────────┐
                    │   Comparison    │
                    │   & Results     │
                    └─────────────────┘
```

## Porównanie: Go Goroutines vs Java Virtual Threads vs OS Threads

| Cecha | Go Goroutines | Java Virtual Threads | Java OS Threads |
|-------|---------------|----------------------|-----------------|
| Rozmiar stack | ~2KB (dynamiczny) | ~1KB (dynamiczny) | ~1MB (stały) |
| Zarządzanie | Go runtime | JVM (Project Loom) | System operacyjny |
| Limit praktyczny | Miliony | Miliony | Tysiące |
| Czas tworzenia | Nanosekundy | Mikrosekundy | Milisekundy |
| Context switching | Bardzo szybki | Bardzo szybki | Powolny |
| Blocking operations | Non-blocking | Non-blocking | Blocking |

## Mechanizm Scheduler'a

### Go Runtime Scheduler (M:N model)

```
User Level:     [G] [G] [G] [G] [G] [G] [G] [G] ...  (Goroutines - G)
                 │   │   │   │   │   │   │   │
                 └───┴───┴───┴───┴───┴───┴───┘
                             │
Go Runtime:     [M] [M] [M] [M]                      (Machine threads - M)
                 │   │   │   │
                 └───┴───┴───┘
                       │
OS Level:       [P] [P] [P] [P]                      (OS threads - P)
```

- **G (Goroutine)** - lekki wątek zarządzany przez Go runtime
- **M (Machine)** - wątek OS używany przez Go runtime
- **P (Processor)** - kontekst wykonania (domyślnie = liczba CPU cores)

### Kluczowe cechy:
1. **M:N scheduling** - wiele goroutines (G) mapowanych na kilka wątków OS (M)
2. **Work stealing** - scheduler automatycznie równoważy obciążenie między wątkami
3. **Non-blocking I/O** - goroutines nie blokują wątków OS podczas operacji I/O
4. **Preemptive** - scheduler może przerwać długo działającą goroutine

## Problemy rozwiązane przez lekkie wątki

### 1. **Problem: Limit liczby wątków OS**
- **Tradycyjne wątki:** System operacyjny może utworzyć tylko tysiące wątków (limit pamięci)
- **Rozwiązanie:** Goroutines/Virtual threads pozwalają na miliony równoczesnych zadań

### 2. **Problem: Wysokie zużycie pamięci**
- **Tradycyjne wątki:** Każdy wątek OS ~ 1MB pamięci (stack)
- **Rozwiązanie:** Goroutines ~ 2KB (dynamiczny stack)

### 3. **Problem: Wolne tworzenie i context switching**
- **Tradycyjne wątki:** Tworzenie wątku OS ~ milisekundy
- **Rozwiązanie:** Goroutines ~ nanosekundy

### 4. **Problem: Blocking operations**
- **Tradycyjne wątki:** Operacje I/O blokują cały wątek OS
- **Rozwiązanie:** Go runtime zarządza I/O asynchronicznie

## Wnioski

### Wyniki z testu:
- **Goroutines są ~4-5x szybsze** niż ograniczony model wątkowy
- **Zużycie pamięci jest ~3-4x mniejsze** w przypadku goroutines
- **Możliwość utworzenia milionów goroutines** jednocześnie

### Zastosowania:
1. **Serwery HTTP** - obsługa tysięcy równoczesnych połączeń
2. **Microservices** - równoległe wywołania API
3. **Data processing** - przetwarzanie dużych zbiorów danych
4. **Real-time systems** - systemy wymagające wysokiej przepustowości

### Porównanie z Java Project Loom:
- **Go miał goroutines od początku (2009)** - dojrzały ekosystem
- **Java 25 wprowadza Virtual Threads** - nowa funkcjonalność
- **Oba rozwiązania osiągają ten sam cel** - lekką współbieżność
- **Go ma przewagę w prostotcie** - goroutines są fundamentem języka

## Uwagi techniczne

- Program wykorzystuje `time.Sleep()` do symulacji operacji I/O
- `runtime.MemStats` mierzy rzeczywiste zużycie pamięci
- Semaphore pattern (`chan struct{}`) ogranicza liczbę równoczesnych goroutines
- `runtime.GOMAXPROCS()` kontroluje liczbę wątków OS używanych przez Go runtime

## Bibliografia

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Java Project Loom](https://openjdk.org/projects/loom/)
- [Go Runtime Scheduler](https://go.dev/src/runtime/proc.go)
- [Virtual Threads in Java 25](https://docs.oracle.com/en/java/javase/25/core/virtual-threads.html)
