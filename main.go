package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// DataCenterMetrics - Serverin anlıq hardware göstəriciləri
type DataCenterMetrics struct {
	CPULoad     int
	Temperature float64
	EnergyUsage int // Vatt (W) ilə
}

func main() {
	fmt.Println("🚀 Go-Based Data Center Core & Security Monitor Başladılır...")

	// Kanallar (Channels): Goroutine-lər arasında təhlükəsiz məlumat ötürülməsi üçün "borular"
	metricsChan := make(chan DataCenterMetrics)
	alertChan := make(chan string)

	// 1. GOROUTINE: Hardware Monitor (Arxa planda fasiləsiz datanı simulyasiya edir)
	go func() {
		for {
			// Real-time sensor oxunmasını simulyasiya edirik
			metrics := DataCenterMetrics{
				CPULoad:     rand.Intn(40) + 40,        // 40% - 80% arası CPU
				Temperature: 45.0 + rand.Float64()*30, // 45°C - 75°C arası istilik
				EnergyUsage: rand.Intn(200) + 300,      // 300W - 500W arası enerji
			}

			// Datanı kanal vasitəsilə mərkəzə ötürürük
			metricsChan <- metrics

			// Hər 2 saniyədən bir oxuma həyata keçirilir
			time.Sleep(2 * time.Second)
		}
	}()

	// 2. GOROUTINE: Kiber-Təhlükəsizlik / Anomaliya Detektoru
	// Bu modul hardware monitorundan asılı olmadan, paralel şəkildə kiber riskləri skan edir
	go func() {
		for {
			time.Sleep(5 * time.Second) // Hər 5 saniyədən bir təhlükəsizlik yoxlanışı

			// Simulyasiya: Əgər CPU yükü qəfil 75%-i keçərsə, potensial DDoS həyəcanı ver
			currentHour := time.Now().Hour()
			if currentHour > 0 { // Davamlı aktiv olması üçün şərti simulyasiya
				// Təsadüfi olaraq kiber təhlükəsizlik insidenti simulyasiyası
				if rand.Float32() > 0.7 {
					alertChan <- "⚠️ TƏHLÜKƏ: Müəyyən olunmamış IP-lərdən yüksək trafik qeydə alındı! Potensial DDoS!"
				}
			}
		}
	}()

	// 3. GOROUTINE: Ağıllı Enerji Loqqeri (Telemetry Logging)
	// Bu modul sensor datalarını qəbul edir və fasiləsiz fayla yazır
	go func() {
		file, err := os.OpenFile("datacenter_energy_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("❌ Log faylı açıla bilmədi:", err)
			return
		}
		defer file.Close()

		for m := range metricsChan {
			logLine := fmt.Sprintf("[%s] CPU: %d%% | Temp: %.1f°C | Energy: %dW\n",
				time.Now().Format("15:04:05"), m.CPULoad, m.Temperature, m.EnergyUsage)

			_, err := file.WriteString(logLine)
			if err != nil {
				fmt.Println("❌ Loq yazıla bilmədi:", err)
			}
		}
	}()

	// MAIN LOOP: Əsas mühərrik kanalları dinləyir və ekrana çıxarır (Non-blocking multiplexing)
	for {
		select {
		case alert := <-alertChan:
			// Kiber-təhlükəsizlik modulundan gələn kritik siqnal dərhal ekrana qırmızı və ya diqqət çəkən tonda çıxır
			fmt.Println("\n=======================================================")
			fmt.Println(alert)
			fmt.Println("=======================================================\n")
		case <-time.After(1 * time.Second):
			// Hər saniyə mərkəzi idarəetmə panelinə sistemin stabil işlədiyini bildiririk
			fmt.Printf("Data Center Core Status: STABIL | Zaman: %s\n", time.Now().Format("15:04:05"))
		}
	}
}