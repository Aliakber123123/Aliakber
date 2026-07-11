package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Server - Enterprise infrastruktur vahidi
type Server struct {
	ID          string  `json:"id"`
	CPULoad     int     `json:"cpu_load"`
	Temperature float64 `json:"temperature"`
	PowerUsage  int     `json:"power_usage"`
	IsActive    bool    `json:"is_active"`
}

// TransactionTask - İcra olunmağı gözləyən irimiqyaslı istifadəçi sorğusu
type TransactionTask struct {
	ID        int
	Payload   string
	CreatedAt time.Time
}

// ClusterManager - 100,000 serverlik miqyası və FinOps qənaətini idarə edən nüvə
type ClusterManager struct {
	mu              sync.RWMutex
	Servers         map[string]*Server
	TotalRequests   int64
	DroppedRequests int64
	MoneySavedUSD   float64
	CostPerWattHour float64
	TaskQueue       chan TransactionTask
}

var manager = ClusterManager{
	Servers:         make(map[string]*Server),
	CostPerWattHour: 0.00015,
	TaskQueue:       make(chan TransactionTask, 5000),
}

func init() {
	// 10 magistral server klaster mərkəzini yaradırıq
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("srv-core-matrix-0%d", i)
		manager.Servers[id] = &Server{
			ID:          id,
			CPULoad:     20,
			Temperature: 45.0,
			PowerUsage:  150,
			IsActive:    true,
		}
	}
}

// BackgroundOptimizer - Hər saniyə FinOps analizini aparan mühərrik
func BackgroundOptimizer() {
	for {
		time.Sleep(1 * time.Second)
		manager.mu.Lock()

		var totalClusterLoad = 0
		var activeCount = 0

		for _, s := range manager.Servers {
			if s.IsActive {
				if s.CPULoad > 10 {
					s.CPULoad -= rand.Intn(10)
				}
				s.Temperature = 35.0 + (float64(s.CPULoad) * 0.4)
				s.PowerUsage = 150 + (s.CPULoad * 5)
				totalClusterLoad += s.CPULoad
				activeCount++
			} else {
				s.PowerUsage = 12
				s.Temperature = 22.5
				s.CPULoad = 0
			}
		}

		averageLoad := 0
		if activeCount > 0 {
			averageLoad = totalClusterLoad / activeCount
		}

		// STRATEGİYA: Resurs İsrafının Qarşısının Alınması (Cost Cutter)
		if averageLoad < 25 && activeCount > 2 {
			for _, s := range manager.Servers {
				if s.IsActive {
					s.IsActive = false
					fmt.Printf("📉 [FINOPS OPTIMIZER] %s söndürüldü. Klaster yükü: %d%%. Resurs israfı dayandırıldı.\n", s.ID, averageLoad)
					break
				}
			}
		}

		// STRATEGİYA: Proaktiv Auto-Scaling (SLA Qoruyucu)
		if (averageLoad > 65 || len(manager.TaskQueue) > 100) && activeCount < 10 {
			for _, s := range manager.Servers {
				if !s.IsActive {
					s.IsActive = true
					s.CPULoad = 40
					fmt.Printf("🚀 [PROACTIVE AUTO-SCALE] Quyruq dolur (%d tapşırıq). %s dərhal oyadıldı!\n", len(manager.TaskQueue), s.ID)
					break
				}
			}
		}

		// Maliyyə qənaət hesabı
		for _, s := range manager.Servers {
			if !s.IsActive {
				savedWatts := 350 - s.PowerUsage
				manager.MoneySavedUSD += (float64(savedWatts) / 1000.0) * (manager.CostPerWattHour / 3600.0) * 100000
			}
		}

		fmt.Printf("📊 [CLUSTERS] Aktiv: %d/10 | Quyruqda Gözləyən: %d | Cəmi Qənaət: $%.2f USD\n", activeCount, len(manager.TaskQueue), manager.MoneySavedUSD)
		manager.mu.Unlock()
	}
}

// CoreWorkerPool - Serverlərin daxilində tapşırıqları paralel icra edən işçi mühərriklər
func CoreWorkerPool() {
	for range manager.TaskQueue {
		manager.mu.Lock()
		var assignedServer *Server
		minLoad := 101

		for _, s := range manager.Servers {
			if s.IsActive && s.CPULoad < minLoad {
				minLoad = s.CPULoad
				assignedServer = s
			}
		}

		if assignedServer != nil {
			assignedServer.CPULoad += 5
			if assignedServer.CPULoad > 100 {
				assignedServer.CPULoad = 100
				atomic.AddInt64(&manager.DroppedRequests, 1)
			}
		} else {
			atomic.AddInt64(&manager.DroppedRequests, 1)
		}
		manager.mu.Unlock()

		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
	}
}

// GetClusterDashboard - Dashboard məlumatları
func GetClusterDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	response := map[string]interface{}{
		"infrastructure_status":      "OPERATIONAL",
		"simulated_global_scale":     "100,000 Enterprise Nodes",
		"total_processed_requests":   atomic.LoadInt64(&manager.TotalRequests),
		"sla_dropped_requests":       atomic.LoadInt64(&manager.DroppedRequests),
		"enterprise_money_saved_usd": fmt.Sprintf("$%.2f", manager.MoneySavedUSD),
		"queue_backlog_count":        len(manager.TaskQueue),
		"cluster_matrix":             manager.Servers,
	}
	json.NewEncoder(w).Encode(response)
}

// SimulateTrafficBurst - Trafik partlayışı simulyatoru
func SimulateTrafficBurst(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	go func() {
		for i := 1; i <= 1000; i++ {
			atomic.AddInt64(&manager.TotalRequests, 1)
			task := TransactionTask{
				ID:        int(atomic.LoadInt64(&manager.TotalRequests)),
				Payload:   "SECURE_TELEMETRY_DATA_PACKET",
				CreatedAt: time.Now(),
			}
			
			select {
			case manager.TaskQueue <- task:
			default:
				atomic.AddInt64(&manager.DroppedRequests, 1)
			}
		}
	}()

	w.Write([]byte(`{"message": "🚀 1000 Yüksək Sürətli Sorğu Klaster Quyruğuna Buraxıldı! Auto-scaling və Worker Pool işə düşür."}`))
}

func main() {
	fmt.Println("⚙️ Enterprise Scaled Cluster Core Manager Engine Start...")

	for i := 1; i <= 5; i++ {
		go CoreWorkerPool()
	}

	go BackgroundOptimizer()

	http.HandleFunc("/api/cluster/dashboard", GetClusterDashboard)
	http.HandleFunc("/api/cluster/simulate-burst", SimulateTrafficBurst)

	fmt.Println("🌐 Enterprise Gateway Ready: http://localhost:8080/api/cluster/dashboard")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("❌ Server Error:", err)
	}
}