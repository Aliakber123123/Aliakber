package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	// 🛠️ REAL SİSTEM DATA KİTABXANASI
	"github.com/shirou/gopsutil/v3/cpu"
)

type AIServerNode struct {
	Name              string
	CurrentLoad       float64 // Artıq bu göstərici həqiqi CPU yükü olacaq!
	Temperature       float64
	EnergyConsumption float64
	mu                sync.Mutex
}

// Artıq təsadüfi rəqəmlər yoxdur, birbaşa kompüterin CPU yükünü oxuyuruq!
func (n *AIServerNode) UpdateRealMetrics() {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Kompüterin anlıq 1 saniyəlik CPU istifadəsini faizlə ölçürük
	percent, err := cpu.Percent(time.Second, false)
	if err == nil && len(percent) > 0 {
		n.CurrentLoad = percent[0] // Həqiqi CPU faizi (məs: %12.5)
	}

	// Real data əsasında temperatur və enerji simulyasiyasını daha məntiqli qururuq
	n.Temperature = 40.0 + (n.CurrentLoad * 0.3)
	n.EnergyConsumption = 150.0 + (n.CurrentLoad * 4.5)
}

type AIDataCenterLoadBalancer struct {
	Nodes                 []*AIServerNode
	IPRequestHistory      map[string][]time.Time
	mu                    sync.Mutex
	RateLimitWindow       time.Duration
	MaxRequestsPerWindow  int
}

func (lb *AIDataCenterLoadBalancer) IsCyberAttack(ipAddress string) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	hazirkiZaman := time.Now()
	if lb.IPRequestHistory[ipAddress] == nil {
		lb.IPRequestHistory[ipAddress] = []time.Time{}
	}

	var validTimes []time.Time
	for _, t := range lb.IPRequestHistory[ipAddress] {
		if hazirkiZaman.Sub(t) < lb.RateLimitWindow {
			validTimes = append(validTimes, t)
		}
	}
	lb.IPRequestHistory[ipAddress] = validTimes

	if len(lb.IPRequestHistory[ipAddress]) >= lb.MaxRequestsPerWindow {
		return true
	}

	lb.IPRequestHistory[ipAddress] = append(lb.IPRequestHistory[ipAddress], hazirkiZaman)
	return false
}

func (lb *AIDataCenterLoadBalancer) DistributeAITask(taskName string, taskWeight float64) string {
	fmt.Printf("\n[ŞƏBƏKƏ SORĞUSU] %s (Ağırlıq: %.1f)\n", taskName, taskWeight)
	
	var bestNode *AIServerNode
	minScore := 999999.0

	for _, node := range lb.Nodes {
		node.UpdateRealMetrics()
		
		if node.Temperature > 85.0 || node.CurrentLoad > 95.0 {
			continue
		}

		// Real yük və real temperatura əsasən ən yaxşı serveri seçirik
		score := (node.CurrentLoad * 0.6) + (node.Temperature * 0.4)
		if score < minScore {
			minScore = score
			bestNode = node
		}
	}

	if bestNode != nil {
		msg := fmt.Sprintf("YÖNLƏNDİRİLDİ: %s -> %s (Anlıq CPU Yükü: %.1f%%)", taskName, bestNode.Name, bestNode.CurrentLoad)
		fmt.Printf("  🟢 %s\n", msg)
		return "SUCCESS: Assigned to " + bestNode.Name
	}

	fmt.Println("  🔴 BÖHRAN: Bütün real serverlər maksimum yükdədir!")
	return "ERROR: All nodes overloaded"
}

func (lb *AIDataCenterLoadBalancer) MonitorStatus() {
	for {
		fmt.Println("\n==================================================")
		fmt.Println("REAL SİSTEM MONITORU (Kompüterinin CPU Göstəriciləri)")
		fmt.Println("==================================================")
		for _, node := range lb.Nodes {
			node.UpdateRealMetrics()
			status := "STABİL"
			if node.CurrentLoad > 80 {
				status = "YÜKSƏK YÜK"
			}
			fmt.Printf("[%s] Real CPU: %.1f%% | Hesablanan Temp: %.1f°C -> [%s]\n", node.Name, node.CurrentLoad, node.Temperature, status)
		}
		fmt.Println("==================================================")
		time.Sleep(3 * time.Second)
	}
}

func main() {
	balancer := &AIDataCenterLoadBalancer{
		Nodes: []*AIServerNode{
			{Name: "AI-Real-Core-Cluster-01"},
		},
		IPRequestHistory:     make(map[string][]time.Time),
		RateLimitWindow:      5 * time.Second,
		MaxRequestsPerWindow: 2,
	}

	go balancer.MonitorStatus()

	listener, err := net.Listen("tcp", "127.0.0.1:5005")
	if err != nil {
		fmt.Println("Server başladılarkən xəta:", err)
		return
	}
	defer listener.Close()

	fmt.Println("\n[MƏLUMAT] Real Sistem Göstəricili Go Server 5005 portunda aktivdir...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		ipAddress := strings.Split(remoteAddr, ":")[0]

		if balancer.IsCyberAttack(ipAddress) {
			fmt.Printf("  ❌ TƏHLÜKƏ: %s adresindən DDoS təsbit edildi! Bloklandı.\n", ipAddress)
			conn.Write([]byte("ERROR: 429 Too Many Requests\n"))
			conn.Close()
			continue
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err == nil && n > 0 {
			data := strings.TrimSpace(string(buf[:n]))
			parts := strings.Split(data, ",")
			if len(parts) == 2 {
				var weight float64
				fmt.Sscanf(parts[1], "%f", &weight)
				response := balancer.DistributeAITask(parts[0], weight)
				conn.Write([]byte(response + "\n"))
			}
		}
		conn.Close()
	}
}