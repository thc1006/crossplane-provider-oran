package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// DeviceState 儲存單一設備的模擬狀態
type DeviceState struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Bandwidth  string `json:"bandwidth"`
	LaserPower string `json:"laserPower"`
}

// ConfigRequest 是來自 Adapter 的設定請求
type ConfigRequest struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Bandwidth  string `json:"bandwidth"`
	LaserPower string `json:"laserPower"`
}

// ConfigResponse 是模擬器回覆給 Adapter 的回應
type ConfigResponse struct {
	Status            string `json:"status"`
	ObservedBandwidth string `json:"observedBandwidth"`
	LastUpdated       string `json:"lastUpdated"`
}

var (
	// 使用 map 來模擬儲存多個設備的狀態
	deviceStore = make(map[string]DeviceState)
	// 使用 mutex 來確保併發安全
	storeMutex = &sync.Mutex{}
)

func configureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "僅支援 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var req ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	deviceStore[req.Hostname] = DeviceState{
		Hostname:   req.Hostname,
		Port:       req.Port,
		Bandwidth:  req.Bandwidth,
		LaserPower: req.LaserPower,
	}

	log.Printf("設定設備: %+v\n", req)

	resp := ConfigResponse{
		Status:            "configured",
		ObservedBandwidth: req.Bandwidth, // 模擬回讀的頻寬與設定值相同
		LastUpdated:       time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deconfigureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "僅支援 DELETE 方法", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Hostname string `json:"hostname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	if _, ok := deviceStore[req.Hostname]; ok {
		delete(deviceStore, req.Hostname)
		log.Printf("移除設備設定: %s\n", req.Hostname)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "deconfigured", "hostname": "%s"}`, req.Hostname)
	} else {
		log.Printf("嘗試移除不存在的設備: %s\n", req.Hostname)
		http.Error(w, "設備不存在", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/configure", configureHandler)
	http.HandleFunc("/deconfigure", deconfigureHandler)

	log.Println("硬體模擬器啟動於 :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
