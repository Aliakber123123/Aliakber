package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// Global mühitdə anlıq metrixləri saxlamaq və thread-safe oxumaq üçün struktur
type SafeMetrics struct {
	mu          sync.RWMutex
	CPULoad     int
	Temperature float64
	EnergyUsage int
}

var currentMetrics SafeMetrics

// RateLimiter - Hər IP üçün daxili sorğu sayğacı
type RateLimiter struct {
	mu           sync.Mutex
	visitorCount map[string]int
	lastReset    time.Time
}

var limiter = RateLimiter{
	visitorCount: make(map[string]int),
	lastReset:    time.Now(),
}

// SecurityMiddleware - Bütün sorğuları yoxlayan kiber-müdafiə divarı
func SecurityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter.mu.Lock()
		
		// Hər 5 saniyədən bir sayğacları sıfırlayırıq ki, normal istifadəçilər blokda qalmasın
		if time.Since(limiter.lastReset) > 5*time.Second {
			limiter.visitorCount = make(map[string]int)
			limiter.lastReset = time.Now()
		}

		// Sorğu göndərən istifadəçinin IP-sini təyin edirik
		visitorIP := r.RemoteAddr

		// Həmin IP-dən gələn sorğu sayını 1 vahid artırırıq
		limiter.visitorCount[visitorIP]++

		// LIMIT: Əgər eyni IP 5 saniyə daxilində 5-dən çox sorğu göndərərsə, blokla!
		if limiter.visitorCount[visitorIP] > 5 {
			limiter.mu.Unlock()
			
			// Kiber-təhlükəsizlik mərkəzinə dərhal xəbərdarlıq basırıq
			fmt.Printf("🛡️ [WAF ALERT] DDoS/Brute-Force təhlükəsi bloklandı! IP: %s | Sorğu sayı: %d\n", visitorIP, limiter.visitorCount[visitorIP])
			
			// Hücum edənə HTTP 429 Too Many Requests statusu qaytarırıq
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "DDoS/Flood Attack Detected. Your IP has been throttled by WAF."}`))
			return
		}
		limiter.mu.Unlock()

		// Əgər hər şey qaydasındadırsa, sorğunun keçməsinə icazə ver
		next(w, r)
	}
}

func main() {
	fmt.Println("🚀 API Gateway & Real-Time WAF Security Monitor Aktivləşdirilir...")

	// 1. GOROUTINE: Real OS CPU Monitor
	go func() {
		for {
			percentages, err := cpu.Percent(time.Second, false)
			cpuVal := 0
			if err == nil && len(percentages) > 0 {
				cpuVal = int(percentages[0])
			}

			currentMetrics.mu.Lock()
			currentMetrics.CPULoad = cpuVal
			currentMetrics.Temperature = 40.0 + (float64(cpuVal) * 0.4)
			currentMetrics.EnergyUsage = 250 + (cpuVal * 3)
			currentMetrics.mu.Unlock()

			// Fayla loq yazılması
			file, err := os.OpenFile("datacenter_energy_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				logLine := fmt.Sprintf("[%s] CPU: %d%% | Temp: %.1f°C | Power: %dW\n",
					time.Now().Format("15:04:05"), cpuVal, currentMetrics.Temperature, currentMetrics.EnergyUsage)
				file.WriteString(logLine)
				file.Close()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// HTTP ENDPOINT: Artıq xüsusi SecurityMiddleware qoruması altındadır!
	http.HandleFunc("/api/metrics", SecurityMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		currentMetrics.mu.RLock()
		data := map[string]interface{}{
			"status":      "active",
			"cpu_load":    fmt.Sprintf("%d%%", currentMetrics.CPULoad),
			"temperature": fmt.Sprintf("%.1f°C", currentMetrics.Temperature),
			"power_w":     currentMetrics.EnergyUsage,
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		}
		currentMetrics.mu.RUnlock()

		json.NewEncoder(w).Encode(data)
	}))

	fmt.Println("🌐 Protected API Server hazır: http://localhost:8080/api/metrics")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("❌ Server başladılarkən xəta yarandı:", err)
	}
}