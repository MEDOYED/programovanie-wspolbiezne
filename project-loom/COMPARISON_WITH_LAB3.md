# Porównanie: Lab 3 vs Project Loom

## Lab 3 - Worker Pool Pattern
**Cel:** Równoległe przetwarzanie plików tekstowych z wykorzystaniem puli workerów

**Koncepcja:**
- Worker Pool Pattern (8 workerów)
- Przetwarzanie 50,000 plików tekstowych
- Zliczanie słów w plikach
- Praktyczna aplikacja do przetwarzania danych

**Mechanizmy:**
- Goroutines jako workery
- Channels do komunikacji (jobs, results)
- sync.WaitGroup do synchronizacji
- Buforowane kanały dla wydajności

**Wynik:**
- 274,824,047 słów przetworzone
- Czas: ~30 sekund
- 8 workerów równomiernie rozłożyło pracę

---

## Project Loom - Lightweight Concurrency Demonstration
**Cel:** Demonstracja różnicy między lekkimi goroutines a tradycyjnym modelem wątkowym

**Koncepcja:**
- Porównanie goroutines vs ograniczony model wątkowy
- 10,000 równoczesnych zadań
- Pomiar wydajności i zużycia zasobów
- Teoretyczna demonstracja koncepcji Project Loom

**Mechanizmy:**
- Goroutines (lightweight) - miliony możliwych
- Limited threads (50) - symulacja OS threads
- Semaphore pattern do ograniczenia współbieżności
- runtime.MemStats do pomiaru pamięci

**Wynik:**
- 20.81x speedup dla goroutines
- Goroutines: 11.5ms, 871,315 tasks/sec
- Limited threads: 238ms, 41,874 tasks/sec
- Drastyczna różnica w wydajności

---

## Kluczowe różnice

| Aspekt | Lab 3 | Project Loom |
|--------|-------|--------------|
| **Cel** | Praktyczna aplikacja | Teoretyczna demonstracja |
| **Zadanie** | Zliczanie słów w plikach | Symulacja lekkich zadań |
| **Skala** | 50,000 plików | 10,000 goroutines |
| **Pattern** | Worker Pool (fixed size) | Unlimited goroutines vs Limited |
| **Fokus** | Efektywne przetwarzanie | Porównanie modeli współbieżności |
| **I/O** | Rzeczywiste operacje I/O | Symulowane opóźnienia |
| **Analogia** | ForkJoinPool (Java) | Virtual Threads (Project Loom) |

---

## Co demonstruje każdy projekt?

### Lab 3 pokazuje:
✓ Praktyczne użycie Worker Pool Pattern
✓ Równoległe przetwarzanie dużych zbiorów danych
✓ Efektywne zarządzanie ograniczoną liczbą workerów
✓ Rzeczywiste operacje I/O (czytanie plików)
✓ Logging i raportowanie wyników

### Project Loom pokazuje:
✓ Fundamentalną różnicę między lekkimi a ciężkimi wątkami
✓ Skalowalność goroutines (miliony vs tysiące)
✓ Wydajność lekkich mechanizmów współbieżności
✓ Analogię między Go goroutines a Java Virtual Threads
✓ Teoretyczne podstawy nowoczesnej współbieżności

---

## Wnioski

**Lab 3** to **praktyczna aplikacja** pokazująca, jak efektywnie wykorzystać Worker Pool Pattern do przetwarzania dużych zbiorów danych. To wzorzec używany w produkcyjnych aplikacjach.

**Project Loom** to **demonstracja koncepcyjna** pokazująca, dlaczego lekkie wątki (goroutines/virtual threads) są rewolucją w programowaniu współbieżnym. To fundament, na którym buduje się nowoczesne aplikacje high-throughput.

Oba projekty są komplementarne:
- Lab 3 pokazuje "co można zrobić"
- Project Loom pokazuje "dlaczego to działa tak dobrze"

---

## Zastosowania w praktyce

**Worker Pool (Lab 3):**
- Batch processing systemów
- ETL pipelines
- Web scrapers
- File processors
- Data transformation

**Lightweight threads (Project Loom):**
- Web servers (miliony połączeń)
- Microservices (równoległe API calls)
- Real-time systems
- Streaming applications
- Reactive programming

---

## Przyszłość współbieżności

Go miał goroutines od 2009 roku → **16 lat przewagi**
Java 25 wprowadza Virtual Threads → **nadrabia zaległości**

**Wniosek:** Model lekkich wątków to przyszłość programowania współbieżnego. Go był pionierem, Java dołącza z Project Loom.
