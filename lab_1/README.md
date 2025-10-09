# Lab 1 - Concurrent Programming

## How to run

```bash
nix develop
go run main.go
```

Program runs for 30 seconds and then stops all threads automatically.

---

## Questions

### What's the difference between PID and TID?

PID (Process ID) is the ID of the whole program. TID (Thread ID) is the ID of a single thread inside that program. One program can have many threads, so one PID can have many TIDs.

### What's the difference between a user thread and a daemon thread?

User threads keep the program alive - the program waits for them to finish before closing. Daemon threads don't keep the program alive - if only daemon threads are left, the program can close and kill them. In Go we don't have daemon threads, we use context to stop goroutines instead.

### Is it bad to not name your threads?

Yes, it's bad practice. Named threads make debugging easier because you can see which thread is doing what in the logs. Without names everything looks the same and it's hard to find problems.

### How important is thread priority?

Not very important in most cases. The operating system decides which thread gets to run anyway, regardless of the priority you set. It can be useful in real-time systems or when one thread needs to respond quickly, but using it wrong can cause problems. Go doesn't let you set priority - it manages goroutines automatically.

### What is a monitor in multithreading?

A monitor is a way to synchronize threads. It combines two things: a lock (so only one thread can access shared data at a time) and condition variables (so threads can wait and wake each other up). In Java you use `synchronized` with `wait()` and `notify()`. In Go you use `sync.Mutex` with `sync.Cond`. In this program the DataBuffer structure works as a monitor.

---

## Implementation

The program uses the producer-consumer pattern with 4 goroutines:

- **T1 (121261#Writer#1)**: Writes random numbers from [21-37]
- **T2 (121261#Writer#2)**: Writes random numbers from [1337-4200]
- **T3 (121261#Reader#1)**: Reads values and adds them up
- **T4 (121261#Monitor#1)**: Shows system status every second

### Synchronization:

- Channel - passes data between threads
- `sync.RWMutex` - protects shared status information
- `context.Context` - controls when goroutines should stop (after 30 seconds)
- `sync.WaitGroup` - waits for all goroutines to finish

The program runs for exactly 30 seconds, then all goroutines stop cleanly.
