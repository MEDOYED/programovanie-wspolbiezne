# Laboratorium 2 - Problem ucztujących filozofów

## Co robi ten program?

To klasyczna symulacja problemu **ucztujących filozofów** (Dining Philosophers Problem) w języku Go. Program pokazuje jak rozwiązać problem synchronizacji współdzielonych zasobów bez deadlocka i starvation. Symulacja działa dokładnie 30 sekund.

## Czym jest problem filozofów?

### Opis problemu
- **5 filozofów** siedzi przy okrągłym stole
- Między każdymi dwoma filozofami leży **jeden widelec** (razem 5 widelców)
- Aby zjeść, filozof potrzebuje **dwóch widelców** - lewego i prawego
- Filozofowie mogą:
  - **Myśleć** - nie potrzebują widelców
  - **Jeść** - muszą mieć oba widelce

### Problem deadlock (zakleszczenie)
Jeśli każdy filozof jednocześnie podniesie swój lewy widelec i będzie czekał na prawy - nastąpi **deadlock**. Żaden nie będzie mógł jeść, bo wszyscy czekają w nieskończoność.

### Problem starvation (zagłodzenie)
Jeśli niektórzy filozofowie ciągle dominują w dostępie do widelców, inni mogą **nigdy nie zjeść**.

## Jak nasze rozwiązanie działa?

### Strategia unikania deadlocka
Program używa prostej ale skutecznej strategii:

- **Filozofowie parzyści** (0, 2, 4) najpierw biorą **prawy** widelec, potem lewy
- **Filozofowie nieparzyści** (1, 3) najpierw biorą **lewy** widelec, potem prawy

To przerywa **cykliczne oczekiwanie** i gwarantuje brak deadlocka!

### Unikanie starvation
- Używamy `sync.Mutex` z Go, który gwarantuje **sprawiedliwość** (fairness)
- Każdy wątek czekający na mutex ma szansę go zdobyć
- Brak priorytetów = żaden filozof nie jest głodzony

## Jak uruchomić?

```bash
go run main.go
```

Program automatycznie zatrzyma się po **dokładnie 30 sekundach**.

## Co zobaczysz?

### Podczas działania:
```
=== Problem ucztujacych filozofow ===
Liczba filozofow: 5
Czas symulacji: 30s
Start symulacji...

Filozof 4 je (cykl 50)
Filozof 0 je (cykl 50)
Filozof 1 je (cykl 50)
...
```

Co 50 cykli jedzenia, filozof wypisuje komunikat.

### Po zakończeniu:
```
=== Koniec symulacji ===
Rzeczywisty czas: 30.07 s

Statystyki:
-------------------------------------------
F0 Myslal: 347, Jadl: 346
F1 Myslal: 348, Jadl: 348
F2 Myslal: 367, Jadl: 367
F3 Myslal: 367, Jadl: 366
F4 Myslal: 390, Jadl: 390
-------------------------------------------
Razem: Myslenie: 1819, Jedzenie: 1817
```

**Ważne**: Wszystkie liczby są zbliżone - to dowód braku starvation!

## Ważne szczegóły techniczne

### Struktury danych

```go
type Fork struct {
    sync.Mutex  // Widelec chroniony mutexem
}

type Philosopher struct {
    id         int        // Numer filozofa (0-4)
    leftFork   *Fork      // Lewy widelec
    rightFork  *Fork      // Prawy widelec
    thinkCount int        // Ile razy myslal
    eatCount   int        // Ile razy jadl
    mu         sync.Mutex // Mutex dla licznikow
}
```

### Cykl życia filozofa

1. **Myśli** (10-50ms) - zwiększa licznik `thinkCount`
2. **Bierze pierwszy widelec** (według strategii)
3. **Sprawdza timeout** - czy symulacja się nie skończyła
4. **Bierze drugi widelec**
5. **Je** (10-50ms) - zwiększa licznik `eatCount`
6. **Oddaje oba widelce**
7. Powrót do kroku 1

### Synchronizacja
Program używa:
- **Goroutines** - każdy filozof to osobna goroutine
- **sync.Mutex** - dla widelców (zasobów współdzielonych)
- **sync.Mutex** - dla liczników statystyk (bezpieczny dostęp)
- **sync.WaitGroup** - do oczekiwania na zakończenie wszystkich filozofów
- **context.WithTimeout** - do automatycznego zatrzymania po 30s

## Odpowiedzi na pytania teoretyczne

### Dlaczego nasze rozwiązanie jest deadlock-free?

**Warunek deadlocka**: Cykliczne oczekiwanie (circular wait)
- F0 czeka na F1
- F1 czeka na F2
- ...
- F4 czeka na F0

**Nasze rozwiązanie**: Różna kolejność pobierania zasobów
- Parzyści: prawy → lewy
- Nieparzyści: lewy → prawy

To **przerywa cykl** i deadlock jest niemożliwy!

### Dlaczego nasze rozwiązanie jest starvation-free?

1. **Go's Mutex jest sprawiedliwy** - używa kolejki FIFO dla oczekujących goroutine
2. **Wszyscy mają równy priorytet** - żaden filozof nie ma przewagi
3. **Statystyki to potwierdzają** - liczby cykli są zbliżone (346-390)

### Czym się różni od Javy?

W Javie (z PDF) używano:
- `ReentrantLock` zamiast `sync.Mutex`
- `CountDownLatch` do synchronizowanego startu
- `Thread` zamiast goroutines

W Go jest prościej:
- `sync.Mutex` - lżejszy i wystarczający
- `context.WithTimeout` - eleganckie zarządzanie czasem
- Goroutines - lżejsze niż wątki systemowe

## Struktura kodu

```
lab_2/
├── main.go      - cały kod programu
├── lab2.pdf     - treść zadania
└── README.md    - ten plik
```

## Funkcje

### `think()`
Symuluje myślenie filozofa. Zwiększa licznik i czeka losowy czas.

### `eat()`
Symuluje jedzenie. Zwiększa licznik i wypisuje komunikat co 50 cykli.

### `dine(ctx, wg)`
Główny cykl życia filozofa. Pętla: myśl → weź widelce → jedz → oddaj widelce.

### `getStatistics()`
Thread-safe odczyt statystyk filozofa.

### `main()`
- Inicjalizacja 5 widelców i 5 filozofów
- Start 5 goroutines
- Oczekiwanie 30 sekund
- Wydruk statystyk i wyjaśnień

## Wymagania zadania

✅ **Liczba filozofów**: 5
✅ **Czas symulacji**: dokładnie 30s
✅ **Brak deadlocka**: program się nie zawiesza
✅ **Brak starvation**: wszyscy filozofowie jedzą podobną ilość razy
✅ **Statystyki**: liczba myślenia i jedzenia dla każdego filozofa
✅ **Uzasadnienie**: wyjaśnienie mechanizmów w kodzie i wydruku

## Technologie

- **Język**: Go (Golang)
- **Współbieżność**: Goroutines
- **Synchronizacja**: sync.Mutex, sync.WaitGroup, context.Context
- **Strategia**: Asymetryczne pobieranie zasobów

## Możliwe rozszerzenia

Jeśli chcesz pobawić się kodem:

1. **Zmień liczbę filozofów**:
   ```go
   const numPhilosophers = 10  // Zamiast 5
   ```

2. **Zmień czas symulacji**:
   ```go
   const duration = 60 * time.Second  // 1 minuta
   ```

3. **Dodaj więcej statystyk**:
   - Średni czas czekania na widelce
   - Maksymalny czas głodzenia
   - Histogram aktywności

4. **Wypróbuj inne strategie**:
   - Arbitr (kelner) kontrolujący dostęp
   - Semafory (max 4 filozofów jednocześnie)
   - Timeout przy próbie zdobycia widelca

## Podsumowanie

Ten program pokazuje jak:
- Unikać deadlocka w problemach synchronizacji
- Gwarantować brak starvation
- Używać goroutines i mutexów w Go
- Mierzyć i raportować wydajność współbieżnych systemów

**Kluczowa lekcja**: Czasem proste rozwiązanie (różna kolejność) jest lepsze niż skomplikowane algorytmy!
