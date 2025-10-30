# Laboratorium 3 - Parallel Word Counter

## Opis
Program zlicza laczna liczbe slow znajdujacych sie we wszystkich plikach tekstowych (.txt) w podanym folderze, korzystajac z rownoleglego przetwarzania (goroutines w Go).

## Funkcjonalnosc
- Rekurencyjne przechodzenie przez folder
- Zliczanie slow w plikach .txt
- Rownolegla obrobka z uzyciem puli 8 workerow
- Generowanie logu z wynikami dla kazdego pliku
- Wyswietlanie lacznej liczby slow

## Uruchomienie

```bash
go run main.go <sciezka_do_folderu>
```

Przyklad:
```bash
go run main.go texts
```

## Wyniki
Program generuje plik `word_count_log.txt` z logami w formacie:
```
Worker-1 -> file1.txt: 3526 slow
Worker-2 -> file2.txt: 4132 slow
...
Laczna liczba slow: 274824047
```

## Implementacja
Program wykorzystuje:
- **WorkerPool** - pula 8 goroutines do rownoleglego przetwarzania
- **Channels** - do komunikacji miedzy workerami (jobs, results)
- **sync.WaitGroup** - do synchronizacji zakonczenia workerow
- **atomic.AddInt32** - do bezpiecznego przydzielania ID workerom

## Architektura
Program implementuje wzorzec Worker Pool, ktory jest odpowiednikiem ForkJoinPool z Javy:
- Glowny watek tworzy pule 8 workerow (goroutines)
- Workery czekaja na zadania w kanale `jobs`
- Kazdy worker pobiera nazwe pliku, liczy slowa i wysyla wynik do kanalu `results`
- Glowny watek zbiera wszystkie wyniki i zapisuje je do pliku logu

## Wynik testowy
Na folderze z 50000 plikami tekstowymi:
- Laczna liczba slow: **274,824,047**
- Czas wykonania: ~30 sekund
- Plik logu: 1.8 MB
- Rozklad pracy: rownomierny pomiedzy 8 workerow

## Porownanie z ForkJoinPool (Java)
| Koncepcja Java | Odpowiednik Go |
|----------------|----------------|
| ForkJoinPool | WorkerPool + goroutines |
| RecursiveTask | worker function |
| fork() | go keyword |
| join() | channel receive |
| ExecutorService | WorkerPool struct |

## Uwagi techniczne
- Program wykorzystuje `bufio.Scanner` z `ScanWords` do efektywnego zliczania slow
- Kanal `jobs` jest buforowany (100 elementow) dla lepszej wydajnosci
- Uzycie `atomic.AddInt32` zapewnia bezpieczne przydzielanie ID bez mutex
- `sync.WaitGroup` zapewnia synchronizacje - program czeka az wszystkie workery zakoncza prace
