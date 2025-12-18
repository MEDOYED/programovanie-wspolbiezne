package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// VehicleType reprezentuje typ pojazdu
type VehicleType int

const (
	Car VehicleType = iota
	Truck
	Motorcycle
)

func (vt VehicleType) String() string {
	switch vt {
	case Car:
		return "Samochód"
	case Truck:
		return "Ciężarówka"
	case Motorcycle:
		return "Motocykl"
	default:
		return "Nieznany"
	}
}

// FuelType reprezentuje typ paliwa
type FuelType int

const (
	Gasoline95 FuelType = iota
	Gasoline98
	Diesel
	LPG
)

func (ft FuelType) String() string {
	switch ft {
	case Gasoline95:
		return "Benzyna 95"
	case Gasoline98:
		return "Benzyna 98"
	case Diesel:
		return "Diesel"
	case LPG:
		return "LPG"
	default:
		return "Nieznany"
	}
}

// Vehicle reprezentuje pojazd tankujący na stacji
type Vehicle struct {
	ID         int
	Type       VehicleType
	FuelType   FuelType
	FuelAmount float64 // litry
	ArrivalTime time.Time
}

// Pump reprezentuje dystrybutor paliwa
type Pump struct {
	ID           int
	FuelTypes    []FuelType
	IsOccupied   bool
	CurrentVehicle *Vehicle
	mutex        sync.Mutex
}

// Statistics przechowuje statystyki stacji
type Statistics struct {
	TotalVehicles      int
	ServedVehicles     int
	TotalFuelDispensed float64
	TotalRevenue       float64
	AverageWaitTime    time.Duration
	TotalWaitTime      time.Duration
	mutex              sync.RWMutex
}

// GasStation reprezentuje stację benzynową
type GasStation struct {
	Pumps         []*Pump
	Queue         chan *Vehicle
	Stats         *Statistics
	Running       bool
	mutex         sync.RWMutex
	wg            sync.WaitGroup
	pumpWg        sync.WaitGroup
}

// Ceny paliwa (za litr)
var fuelPrices = map[FuelType]float64{
	Gasoline95: 6.50,
	Gasoline98: 7.20,
	Diesel:     6.80,
	LPG:        3.50,
}

// NewGasStation tworzy nową stację benzynową
func NewGasStation(numPumps int) *GasStation {
	gs := &GasStation{
		Pumps:   make([]*Pump, numPumps),
		Queue:   make(chan *Vehicle, 50),
		Stats:   &Statistics{},
		Running: true,
	}

	// Inicjalizacja dystrybutorów
	for i := 0; i < numPumps; i++ {
		gs.Pumps[i] = &Pump{
			ID:        i + 1,
			FuelTypes: []FuelType{Gasoline95, Gasoline98, Diesel, LPG},
		}
	}

	return gs
}

// Start uruchamia stację benzynową
func (gs *GasStation) Start() {
	// Uruchom goroutines dla każdego dystrybutora
	for _, pump := range gs.Pumps {
		gs.pumpWg.Add(1)
		go gs.runPump(pump)
	}

	// Goroutine do monitorowania statystyk
	go gs.monitorStatistics()

	// Goroutine do wyświetlania interfejsu użytkownika
	go gs.displayUI()
}

// runPump obsługuje pojedynczy dystrybutor
func (gs *GasStation) runPump(pump *Pump) {
	defer gs.pumpWg.Done()

	for {
		gs.mutex.RLock()
		running := gs.Running
		gs.mutex.RUnlock()

		if !running {
			break
		}

		// Czekaj na pojazd z kolejki
		select {
		case vehicle := <-gs.Queue:
			gs.serveVehicle(pump, vehicle)
		case <-time.After(100 * time.Millisecond):
			// Timeout, aby móc sprawdzić status Running
		}
	}
}

// serveVehicle obsługuje pojazd na dystrybutorze
func (gs *GasStation) serveVehicle(pump *Pump, vehicle *Vehicle) {
	// Zajmij dystrybutor
	pump.mutex.Lock()
	pump.IsOccupied = true
	pump.CurrentVehicle = vehicle
	pump.mutex.Unlock()

	waitTime := time.Since(vehicle.ArrivalTime)

	// Symulacja tankowania (różny czas w zależności od ilości paliwa)
	refuelingTime := time.Duration(vehicle.FuelAmount*100) * time.Millisecond
	time.Sleep(refuelingTime)

	// Oblicz koszt
	cost := vehicle.FuelAmount * fuelPrices[vehicle.FuelType]

	// Aktualizuj statystyki
	gs.Stats.mutex.Lock()
	gs.Stats.ServedVehicles++
	gs.Stats.TotalFuelDispensed += vehicle.FuelAmount
	gs.Stats.TotalRevenue += cost
	gs.Stats.TotalWaitTime += waitTime
	gs.Stats.AverageWaitTime = gs.Stats.TotalWaitTime / time.Duration(gs.Stats.ServedVehicles)
	gs.Stats.mutex.Unlock()

	// Zwolnij dystrybutor
	pump.mutex.Lock()
	pump.IsOccupied = false
	pump.CurrentVehicle = nil
	pump.mutex.Unlock()
}

// AddVehicle dodaje pojazd do kolejki
func (gs *GasStation) AddVehicle(vehicle *Vehicle) {
	vehicle.ArrivalTime = time.Now()

	gs.Stats.mutex.Lock()
	gs.Stats.TotalVehicles++
	gs.Stats.mutex.Unlock()

	gs.Queue <- vehicle
}

// monitorStatistics monitoruje i loguje statystyki
func (gs *GasStation) monitorStatistics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		gs.mutex.RLock()
		running := gs.Running
		gs.mutex.RUnlock()

		if !running {
			break
		}

		<-ticker.C
		// Statystyki są wyświetlane w displayUI
	}
}

// displayUI wyświetla interfejs użytkownika
func (gs *GasStation) displayUI() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		gs.mutex.RLock()
		running := gs.Running
		gs.mutex.RUnlock()

		if !running {
			break
		}

		<-ticker.C
		clearScreen()

		fmt.Println(" ")
		fmt.Println("                    SYMULACJA STACJI BENZYNOWEJ                            ")
		fmt.Println(" ")
		fmt.Println()

		// Wyświetl status dystrybutorów
		fmt.Println(" ")
		fmt.Println("STATUS DYSTRYBUTORÓW")
		fmt.Println(" ")
	

		for _, pump := range gs.Pumps {
			pump.mutex.Lock()
			if pump.IsOccupied && pump.CurrentVehicle != nil {
				fmt.Printf("  Dystrybutor %d: [ZAJĘTY]   Pojazd #%d (%s, %s, %.1fL)\n",
					pump.ID,
					pump.CurrentVehicle.ID,
					pump.CurrentVehicle.Type,
					pump.CurrentVehicle.FuelType,
					pump.CurrentVehicle.FuelAmount)
			} else {
				fmt.Printf("  Dystrybutor %d: [WOLNY]\n", pump.ID)
			}
			pump.mutex.Unlock()
		}
		fmt.Println()

		// Wyświetl statystyki
		gs.Stats.mutex.RLock()
		fmt.Println("STATYSTYKI")
		fmt.Println(" ")
		fmt.Printf("  Pojazdy w kolejce:        %d\n", len(gs.Queue))
		fmt.Printf("  Pojazdy łącznie:          %d\n", gs.Stats.TotalVehicles)
		fmt.Printf("  Obsłużone pojazdy:        %d\n", gs.Stats.ServedVehicles)
		fmt.Printf("  Zużyte paliwo:            %.2f L\n", gs.Stats.TotalFuelDispensed)
		fmt.Printf("  Przychód:                 %.2f PLN\n", gs.Stats.TotalRevenue)
		if gs.Stats.ServedVehicles > 0 {
			fmt.Printf("  Średni czas oczekiwania:  %v\n", gs.Stats.AverageWaitTime.Round(time.Millisecond))
		} else {
			fmt.Printf("  Średni czas oczekiwania:  N/A\n")
		}
		gs.Stats.mutex.RUnlock()
		fmt.Println()
		fmt.Println("Naciśnij Ctrl+C aby zakończyć symulację...")
	}
}

// Stop zatrzymuje stację benzynową
func (gs *GasStation) Stop() {
	gs.mutex.Lock()
	gs.Running = false
	gs.mutex.Unlock()

	// Poczekaj na zakończenie wszystkich dystrybutorów
	gs.pumpWg.Wait()

	// Zamknij kolejkę
	close(gs.Queue)
}

// generateRandomVehicle generuje losowy pojazd
func generateRandomVehicle(id int) *Vehicle {
	vehicleTypes := []VehicleType{Car, Truck, Motorcycle}
	fuelTypes := []FuelType{Gasoline95, Gasoline98, Diesel, LPG}

	vType := vehicleTypes[rand.Intn(len(vehicleTypes))]
	fType := fuelTypes[rand.Intn(len(fuelTypes))]

	var fuelAmount float64
	switch vType {
	case Car:
		fuelAmount = 20 + rand.Float64()*40 // 20-60 litrów
	case Truck:
		fuelAmount = 50 + rand.Float64()*150 // 50-200 litrów
	case Motorcycle:
		fuelAmount = 5 + rand.Float64()*15 // 5-20 litrów
	}

	return &Vehicle{
		ID:         id,
		Type:       vType,
		FuelType:   fType,
		FuelAmount: fuelAmount,
	}
}

// clearScreen czyści ekran konsoli
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Utwórz stację z 4 dystrybutorami
	station := NewGasStation(4)

	// Uruchom stację
	station.Start()

	// Goroutine generująca pojazdy
	go func() {
		vehicleID := 1
		for {
			station.mutex.RLock()
			running := station.Running
			station.mutex.RUnlock()

			if !running {
				break
			}

			// Generuj nowy pojazd co 1-3 sekundy
			time.Sleep(time.Duration(1000+rand.Intn(2000)) * time.Millisecond)

			vehicle := generateRandomVehicle(vehicleID)
			station.AddVehicle(vehicle)
			vehicleID++
		}
	}()

	// Czekaj na przerwanie (Ctrl+C)
	fmt.Println("Stacja benzynowa uruchomiona. Naciśnij Ctrl+C aby zakończyć.")
	time.Sleep(60 * time.Second) // Symulacja trwa 60 sekund

	// Zatrzymaj stację
	station.Stop()

	// Wyświetl ostateczne statystyki
	clearScreen()
	fmt.Println("\nPODSUMOWANIE SYMULACJI\n")

	station.Stats.mutex.RLock()
	fmt.Printf("Łączna liczba pojazdów:       %d\n", station.Stats.TotalVehicles)
	fmt.Printf("Obsłużone pojazdy:            %d\n", station.Stats.ServedVehicles)
	fmt.Printf("Pojazdy w kolejce:            %d\n", len(station.Queue))
	fmt.Printf("Łączne zużycie paliwa:        %.2f L\n", station.Stats.TotalFuelDispensed)
	fmt.Printf("Łączny przychód:              %.2f PLN\n", station.Stats.TotalRevenue)
	if station.Stats.ServedVehicles > 0 {
		fmt.Printf("Średni czas oczekiwania:      %v\n", station.Stats.AverageWaitTime.Round(time.Millisecond))
	}
	station.Stats.mutex.RUnlock()

	fmt.Println("\nSymulacja zakończona.")
}
