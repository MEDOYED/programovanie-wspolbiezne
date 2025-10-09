# Laboratorium 1 - Programowanie współbieżne

## Co robi ten program?

To prosty program w języku Go, który pokazuje jak działają wątki (goroutines). Program działa przez 30 sekund i pokazuje jak różne wątki mogą współdzielić dane.

## Jak to działa?

Program ma 4 wątki, które pracują jednocześnie:

### 1. Pierwszy Producent (Writer #1)
- Losuje liczby od 21 do 37
- Zapisuje je do wspólnego bufora
- Czeka losowy czas (400-1200ms) między operacjami

### 2. Drugi Producent (Writer #2)
- Losuje liczby od 1337 do 4200
- Zapisuje je do wspólnego bufora
- Czeka losowy czas (400-1200ms) między operacjami

### 3. Konsument (Reader #1)
- Odczytuje liczby z bufora
- Dodaje je do sumy
- Wypisuje aktualną sumę
- Czeka 150ms między próbami odczytu

### 4. Monitor
- Co 1 sekundę wypisuje raport:
  - Jaka wartość jest teraz w buforze
  - Który wątek ostatnio zapisał
  - Który wątek ostatnio odczytał
  - Ile operacji zostało wykonanych

## Jak uruchomić?

```bash
go run main.go
```

Program automatycznie zatrzyma się po 30 sekundach.

## Co zobaczysz?

W konsoli zobaczysz:
- Komunikaty od producentów, którzy zapisują liczby
- Komunikaty od konsumenta, który czyta liczby i pokazuje sumę
- Co sekundę raport ze stanu bufora

## Ważne szczegóły techniczne

### Synchronizacja
Program używa:
- **Kanałów (channels)** - do bezpiecznego przekazywania danych między wątkami
- **Mutex (RWMutex)** - do ochrony informacji o stanie bufora
- **Context** - do kontrolowania czasu działania programu

### Nazewnictwo wątków
Każdy wątek ma nazwę w formacie: `121261#Writer#1`
- `121261` - numer indeksu studenta
- `Writer/Reader/Monitor` - typ wątku
- `1` - numer wątku

## Odpowiedzi na pytania teoretyczne

### Różnica między PID a TID
- **PID (Process ID)** - to numer całego programu w systemie
- **TID (Thread ID)** - to numer konkretnego wątku w programie

Jeden program (PID) może mieć wiele wątków (TID).

### Wątek użytkownika vs wątek demoniczny
- **Wątek użytkownika** - program czeka aż się skończy przed zamknięciem
- **Wątek demoniczny** - zostaje zabity gdy główny program się kończy

Wątki demoniczne są przydatne do zadań w tle, które nie są krytyczne.

### Czy nazywać wątki?
**Tak, to dobra praktyka!** Nazwy pomagają:
- Znaleźć błędy w programie
- Zobaczyć co się dzieje w logach
- Debugować program

### Priorytet wątku
Priorytet jest **mało ważny** w większości programów, bo:
- System sam decyduje jak przydzielać czas
- W Go nie możemy bezpośrednio ustawiać priorytetów
- Lepiej skupić się na dobrej synchronizacji

### Czym jest monitor?
Monitor to mechanizm synchronizacji, który:
- Chroni współdzielone dane
- Pozwala tylko jednemu wątkowi na dostęp naraz
- W Go to jest **Mutex** lub kanały (channels)

## Struktura kodu

```
lab_1/
├── main.go      - cały kod programu
├── lab1.pdf     - treść zadania
└── README.md    - ten plik
```

## Technologie

- **Język**: Go (Golang)
- **Współbieżność**: Goroutines + Channels
- **Synchronizacja**: sync.RWMutex, context.Context
