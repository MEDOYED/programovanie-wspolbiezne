# Problem ucztujących filozofów - Rozwiązanie z Arbiterem

## Opis zadania

To zadanie jest analizą porównawczą dwóch różnych rozwiązań klasycznego **problemu ucztujących filozofów** (Dining Philosophers Problem). Implementacja wykorzystuje **wzorzec Arbiter/Waiter** i jest porównywana z wcześniejszą implementacją używającą **asymetrycznego pobierania zasobów** (lab_2).

## Czym jest problem filozofów?

### Opis problemu
- **5 filozofów** siedzi przy okrągłym stole
- Między każdymi dwoma filozofami leży **jeden widelec** (razem 5 widelców)
- Aby zjeść, filozof potrzebuje **dwóch widelców** - lewego i prawego
- Filozofowie mogą:
  - **Myśleć** - nie potrzebują widelców
  - **Jeść** - muszą mieć oba widelce

### Wyzwania synchronizacji

#### 1. Deadlock (zakleszczenie)
Jeśli każdy filozof jednocześnie podniesie swój lewy widelec i będzie czekał na prawy - nastąpi **deadlock**. Wszyscy czekają w nieskończoność.

#### 2. Starvation (zagłodzenie)
Jeśli niektórzy filozofowie dominują w dostępie do widelców, inni mogą **nigdy nie zjeść**.

## Nasze rozwiązanie: Wzorzec Arbiter/Waiter

### Jak działa?

W tym rozwiązaniu wprowadzamy **Kelnera (Waiter/Arbiter)** - centralną jednostkę kontrolującą dostęp do jadalni:

1. **Kelner ma ograniczoną liczbę "pozwoleń"** - maksymalnie **N-1** (czyli 4 z 5 filozofów)
2. Filozof **musi poprosić kelnera o pozwolenie** przed próbą wzięcia widelców
3. Jeśli kelner nie ma wolnych pozwoleń, filozof czeka
4. Po zjedzeniu filozof **zwraca pozwolenie** kelnerowi

### Dlaczego to działa?

#### Brak deadlocka
Jeśli tylko **4 filozofów** może być w jadalni jednocześnie:
- Co najmniej **jeden filozof** nie trzyma żadnego widelca
- Ten filozof nie blokuje żadnego zasobu
- **Cykliczne oczekiwanie jest niemożliwe** → brak deadlocka

#### Brak starvation
- Semafory w Go działają w kolejności **FIFO** (First In, First Out)
- Każdy filozof czekający na pozwolenie dostanie je **po kolei**
- Nikt nie jest pomijany → brak starvation

### Schemat działania

```
                    ┌──────────────┐
                    │   WAITER     │
                    │  (Arbiter)   │
                    │ Max 4/5 slots│
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         ┌────▼───┐   ┌───▼────┐   ┌──▼─────┐
         │ Phil 0 │   │ Phil 1 │   │ Phil 2 │ ...
         └────┬───┘   └────┬───┘   └────┬───┘
              │            │            │
    Request   │            │            │   Request
    Permission│            │            │   Permission
              ▼            ▼            ▼
         ┌────────────────────────────────┐
         │      Forks (Resources)         │
         └────────────────────────────────┘
```

## Jak uruchomić?

```bash
cd Filozofowie-ADA
go run main.go
```

Program automatycznie zatrzyma się po **dokładnie 30 sekundach**.

## Co zobaczysz?

### Podczas działania:
```
=== Problem ucztujacych filozofow - Rozwiazanie z Arbiterem ===
Liczba filozofow: 5
Czas symulacji: 30s
Strategia: Waiter/Arbiter (maksymalnie 4 filozofow jednoczesnie)
Start symulacji...

Filozof 3 je (cykl 50)
Filozof 1 je (cykl 50)
Filozof 4 je (cykl 50)
...
```

### Po zakończeniu:
```
=== Koniec symulacji ===
Rzeczywisty czas: 30.07 s

Statystyki:
-------------------------------------------
F0 Myslal: 352, Jadl: 351
F1 Myslal: 344, Jadl: 343
F2 Myslal: 361, Jadl: 360
F3 Myslal: 358, Jadl: 357
F4 Myslal: 366, Jadl: 365
-------------------------------------------
Razem: Myslenie: 1781, Jedzenie: 1776
```

**Obserwacja**: Liczby są zbliżone - dowód na brak starvation!

## Szczegóły techniczne

### Struktury danych

```go
// Waiter controls access using semaphore pattern
type Waiter struct {
    semaphore chan struct{}  // Buffered channel (size = N-1)
}

type Fork struct {
    sync.Mutex
    id int
}

type Philosopher struct {
    id         int
    leftFork   *Fork
    rightFork  *Fork
    waiter     *Waiter      // Reference to arbiter
    thinkCount int
    eatCount   int
    mu         sync.Mutex
}
```

### Algorytm działania filozofa

```
1. LOOP (dopóki nie timeout):
2.   Myśl (10-50ms)
3.   Poproś kelnera o pozwolenie
4.   CZEKAJ na pozwolenie od kelnera
5.   Weź lewy widelec (Lock)
6.   Weź prawy widelec (Lock)
7.   Jedz (10-50ms)
8.   Oddaj prawy widelec (Unlock)
9.   Oddaj lewy widelec (Unlock)
10.  Zwróć pozwolenie kelnerowi
11. END LOOP
```

### Synchronizacja
- **Buffered channel** - implementacja semafora dla kelnera
- **sync.Mutex** - ochrona widelców (zasobów współdzielonych)
- **sync.Mutex** - ochrona liczników statystyk
- **sync.WaitGroup** - oczekiwanie na zakończenie wszystkich filozofów
- **context.WithTimeout** - automatyczne zatrzymanie po 30s

## Analiza porównawcza z lab_2

### Podejście 1: Asymetryczne pobieranie zasobów (lab_2)

**Strategia:**
- Filozofowie **parzyści** (0, 2, 4): prawy widelec → lewy widelec
- Filozofowie **nieparzyści** (1, 3): lewy widelec → prawy widelec

**Mechanizm zapobiegania deadlock:**
- Przerywa **cykliczne oczekiwanie** (circular wait)
- Jeden filozof zawsze pobiera w odwrotnej kolejności
- Nie może powstać sytuacja, gdzie wszyscy czekają w kółko

**Zalety:**
- ✅ Bardzo proste w implementacji (jedna linijka warunku)
- ✅ Brak dodatkowych struktur danych
- ✅ Maksymalna równoległość - wszyscy mogą próbować jeść
- ✅ Zero overhead - brak centralnej kontroli

**Wady:**
- ❌ Trudniejsze do zrozumienia (dlaczego różna kolejność?)
- ❌ Wymaga "sztuczki" z asymetrią
- ❌ Trudne do skalowania na inne liczby filozofów

### Podejście 2: Arbiter/Waiter (Filozofowie-ADA)

**Strategia:**
- **Kelner** (Arbiter) kontroluje dostęp do jadalni
- Maksymalnie **N-1 filozofów** może jednocześnie próbować jeść
- Filozofowie pobierają widelce w tej samej kolejności (lewy → prawy)

**Mechanizm zapobiegania deadlock:**
- **Ogranicza liczbę konkurentów** o zasoby
- Zawsze jest przynajmniej jeden filozof poza jadalnią
- Niemożliwe jest, żeby wszyscy trzymali jeden widelec

**Zalety:**
- ✅ Intuicyjne i łatwe do zrozumienia
- ✅ Centralna kontrola - łatwe debugowanie
- ✅ Łatwe do modyfikacji (zmiana limitu)
- ✅ Uniwersalne - działa dla dowolnej liczby filozofów

**Wady:**
- ❌ Wymaga dodatkowej struktury (Waiter)
- ❌ Jeden filozof zawsze czeka (mniejsza równoległość)
- ❌ Minimal overhead z powodu centralnego arbitra

### Porównanie wydajności

| Cecha | Asymetryczne (lab_2) | Arbiter (Filozofowie-ADA) |
|-------|----------------------|---------------------------|
| **Struktura kodu** | Prostsza | Bardziej złożona |
| **Zrozumiałość** | Wymaga wyjaśnienia | Intuicyjna |
| **Równoległość** | 5/5 filozofów | 4/5 filozofów |
| **Overhead** | Brak | Minimalny (semafore) |
| **Średnia liczba jedzeń** | ~345-390 | ~343-365 |
| **Skalowalność** | Trudna | Łatwa |
| **Zastosowanie w praktyce** | Akademickie | Produkcyjne |

### Wyniki z testów (30 sekund)

**Lab_2 (Asymetryczne):**
```
F0 Myslal: 347, Jadl: 346
F1 Myslal: 348, Jadl: 348
F2 Myslal: 367, Jadl: 367
F3 Myslal: 367, Jadl: 366
F4 Myslal: 390, Jadl: 390
Razem: Jedzenie: 1817
```

**Filozofowie-ADA (Arbiter):**
```
F0 Myslal: 352, Jadl: 351
F1 Myslal: 344, Jadl: 343
F2 Myslal: 361, Jadl: 360
F3 Myslal: 358, Jadl: 357
F4 Myslal: 366, Jadl: 365
Razem: Jedzenie: 1776
```

**Obserwacja:**
- Asymetryczne: **1817 jedzeń** (średnio 363.4 na filozofa)
- Arbiter: **1776 jedzeń** (średnio 355.2 na filozofa)
- Różnica: **~2.3%** - Arbiter jest minimalnie wolniejszy z powodu ograniczenia do 4 filozofów

### Które rozwiązanie jest lepsze?

**Akademickie / Teoria:**
- **Asymetryczne** - minimalistyczne, zero overhead, maksymalna wydajność

**Produkcja / Praktyka:**
- **Arbiter** - łatwiejsze do zrozumienia, debugowania i modyfikacji

**Wniosek:**
W praktycznych zastosowaniach (np. zarządzanie pulą połączeń do bazy danych), **wzorzec Arbiter** jest preferowany z powodu:
- Czytelności kodu
- Łatwości utrzymania
- Możliwości dynamicznej konfiguracji limitu

## Inne rozwiązania problemu filozofów

### 1. Resource Hierarchy (Numerowanie zasobów)
- Widelce mają numery
- Filozof zawsze bierze najpierw widelec o **niższym numerze**
- Prosty pomysł Dijkstry

### 2. Chandy/Misra (Tokeny)
- Widelce reprezentowane jako **tokeny** (messages)
- Filozofowie wymieniają się tokenami
- Używane w systemach rozproszonych

### 3. Timeout
- Filozof czeka określony czas na drugi widelec
- Jeśli nie dostanie, **oddaje pierwszy** i próbuje ponownie
- Może prowadzić do livelock

## Zastosowania w praktyce

Problem filozofów to **metafora** dla wielu rzeczywistych problemów synchronizacji:

### 1. Pule połączeń do bazy danych
- **Filozofowie** → wątki aplikacji
- **Widelce** → połączenia DB
- **Rozwiązanie** → Connection pool manager (Arbiter!)

### 2. Zarządzanie zasobami w chmurze
- **Filozofowie** → mikroserwisy
- **Widelce** → zasoby CPU/pamięć
- **Rozwiązanie** → Resource scheduler (Kubernetes)

### 3. Systemy plików rozproszonych
- **Filozofowie** → procesy
- **Widelce** → locki na pliki
- **Rozwiązanie** → Distributed lock manager

## Wymagania zadania

✅ **Implementacja w Go** - język zgodny z wymaganiami
✅ **Rozwiązuje problem filozofów** - wzorzec Arbiter/Waiter
✅ **Brak deadlocka** - udowodniono teoretycznie i praktycznie
✅ **Brak starvation** - statystyki pokazują sprawiedliwość
✅ **Czas symulacji: 30s** - dokładnie według specyfikacji
✅ **Statystyki** - szczegółowe dane dla każdego filozofa
✅ **Analiza porównawcza** - szczegółowe porównanie z lab_2
✅ **Dokumentacja po polsku** - kompletny opis rozwiązania

## Technologie

- **Język**: Go (Golang) 1.21+
- **Współbieżność**: Goroutines
- **Synchronizacja**:
  - `sync.Mutex` - ochrona zasobów
  - `chan struct{}` - semafor (buffered channel)
  - `sync.WaitGroup` - synchronizacja zakończenia
  - `context.Context` - kontrola czasu życia

## Struktura projektu

```
Filozofowie-ADA/
├── main.go      - implementacja rozwiązania z Arbiterem
├── go.mod       - definicja modułu Go
└── README.md    - ten plik (dokumentacja i analiza)
```

## Funkcje w kodzie

### `Waiter` (Kelner/Arbiter)
- `NewWaiter(maxDiners)` - tworzy kelnera z limitem
- `RequestPermission(ctx)` - prosi o pozwolenie (blokuje jeśli brak)
- `ReleasePermission()` - zwraca pozwolenie

### `Philosopher` (Filozof)
- `think()` - symuluje myślenie, zwiększa licznik
- `eat()` - symuluje jedzenie, wypisuje co 50 cykli
- `dine(ctx, wg)` - główny cykl życia filozofa
- `getStatistics()` - thread-safe odczyt statystyk

### `main()`
- Inicjalizacja kelnera, widelców i filozofów
- Start goroutines dla każdego filozofa
- Oczekiwanie 30 sekund
- Wydruk statystyk i wyjaśnień

## Wnioski

### Teoretyczne
1. **Wzorzec Arbiter** to klasyczne, eleganckie rozwiązanie problemu filozofów
2. **Ograniczenie liczby konkurentów** (N-1) gwarantuje brak deadlocka
3. **Semafory FIFO** w Go zapewniają sprawiedliwość i brak starvation
4. **Różne strategie** (asymetryczne vs arbiter) mają różne trade-offy

### Praktyczne
1. W realnych systemach **wzorzec Arbiter** jest preferowany (czytelność)
2. **Overhead jest minimalny** (~2-3% spadek wydajności)
3. **Go's goroutines + channels** są doskonałe do problemów synchronizacji
4. **Dokumentacja i czytelność** kodu są ważniejsze niż minimalna optymalizacja

### Porównanie rozwiązań
- **Asymetryczne** (lab_2): szybsze, ale mniej intuicyjne
- **Arbiter** (Filozofowie-ADA): wolniejsze o ~2%, ale znacznie bardziej zrozumiałe

**Rekomendacja**: W produkcyjnym kodzie używaj **wzorca Arbiter** - łatwiej go zrozumieć, debugować i modyfikować.

## Bibliografia

- E. W. Dijkstra - "Hierarchical ordering of sequential processes" (1971)
- K. M. Chandy, J. Misra - "The Drinking Philosophers Problem" (1984)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Dining Philosophers Problem - Wikipedia](https://en.wikipedia.org/wiki/Dining_philosophers_problem)

---

**Autor**: Maksym
**Data**: 2025-12-17
**Kurs**: Programowanie Współbieżne
**Zadanie**: Filozofowie ADA - analiza porównawcza rozwiązań
