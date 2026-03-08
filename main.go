package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
)

type MEMORYSTATUSEX struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// WMI Structs
type Win32_Processor struct {
	Name          string
	NumberOfCores uint32
	MaxClockSpeed uint32
}

type Win32_OperatingSystem struct {
	Caption     string
	Version     string
	BuildNumber string
}

type Win32_LogicalDisk struct {
	DeviceID   string
	FreeSpace  uint64
	Size       uint64
	FileSystem string
}

type Win32_NetworkAdapterConfiguration struct {
	Description          string
	IPAddress            []string
	DNSServerSearchOrder []string
	IPEnabled            bool
}

type AntiVirusProduct struct {
	DisplayName string
}

type FirewallProduct struct {
	DisplayName string
}

// Result struct
type SystemInfo struct {
	CPU struct {
		Name           string `json:"name"`
		Sockets        uint32 `json:"sockets"`
		CoresPerSocket uint32 `json:"cores_per_socket"`
		TotalCores     uint32 `json:"total_cores"`
		Speed          uint32 `json:"speed_mhz"`
	} `json:"cpu"`
	Memory struct {
		PhysicalTotalGB float64 `json:"physical_total_gb"`
		FreePhysicalGB  float64 `json:"free_physical_gb"`
		TotalVirtualGB  float64 `json:"total_virtual_gb"`
		FreeVirtualGB   float64 `json:"free_virtual_gb"`
		LoadPercentage  uint32  `json:"load_percentage"`
	} `json:"memory"`
	Disks []DiskInfo `json:"disks"`
	OS    struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Build   string `json:"build"`
	} `json:"os"`
	Antivirus []string `json:"antivirus"`
	Firewall  []string `json:"firewall"`
	Network   []NetInfo `json:"network"`
}

type DiskInfo struct {
	Drive      string  `json:"drive"`
	TotalGB    float64 `json:"total_gb"`
	FreeGB     float64 `json:"free_gb"`
	UsedGB     float64 `json:"used_gb"`
	FileSystem string  `json:"file_system"`
}

type NetInfo struct {
	Description string   `json:"description"`
	IPs         []string `json:"ips"`
	DNS         []string `json:"dns"`
}

func main() {
	var info SystemInfo
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Ensure output directory exists
	outputDir := "./output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Fetch Memory using syscall
	wg.Add(1)
	go func() {
		defer wg.Done()
		var mem MEMORYSTATUSEX
		mem.Length = uint32(unsafe.Sizeof(mem))
		ret, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)))
		if ret != 0 {
			mu.Lock()
			info.Memory.PhysicalTotalGB = float64(mem.TotalPhys) / (1024 * 1024 * 1024)
			info.Memory.FreePhysicalGB = float64(mem.AvailPhys) / (1024 * 1024 * 1024)
			info.Memory.TotalVirtualGB = float64(mem.TotalPageFile) / (1024 * 1024 * 1024)
			info.Memory.FreeVirtualGB = float64(mem.AvailPageFile) / (1024 * 1024 * 1024)
			info.Memory.LoadPercentage = mem.MemoryLoad
			mu.Unlock()
		}
	}()

	// Fetch CPU
	wg.Add(1)
	go func() {
		defer wg.Done()
		var dst []Win32_Processor
		if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err == nil && len(dst) > 0 {
			mu.Lock()
			info.CPU.Name = dst[0].Name
			info.CPU.Sockets = uint32(len(dst))
			info.CPU.CoresPerSocket = dst[0].NumberOfCores
			info.CPU.TotalCores = 0
			for _, p := range dst {
				info.CPU.TotalCores += p.NumberOfCores
			}
			info.CPU.Speed = dst[0].MaxClockSpeed
			mu.Unlock()
		}
	}()

	// Fetch OS
	wg.Add(1)
	go func() {
		defer wg.Done()
		var dst []Win32_OperatingSystem
		if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err == nil && len(dst) > 0 {
			mu.Lock()
			info.OS.Name = dst[0].Caption
			info.OS.Version = dst[0].Version
			info.OS.Build = dst[0].BuildNumber
			mu.Unlock()
		}
	}()

	// Fetch Disks
	wg.Add(1)
	go func() {
		defer wg.Done()
		var dst []Win32_LogicalDisk
		q := "SELECT DeviceID, FreeSpace, Size, FileSystem FROM Win32_LogicalDisk WHERE DriveType = 3"
		if err := wmi.Query(q, &dst); err == nil {
			mu.Lock()
			for _, d := range dst {
				info.Disks = append(info.Disks, DiskInfo{
					Drive:      d.DeviceID,
					TotalGB:    float64(d.Size) / (1024 * 1024 * 1024),
					FreeGB:     float64(d.FreeSpace) / (1024 * 1024 * 1024),
					UsedGB:     float64(d.Size-d.FreeSpace) / (1024 * 1024 * 1024),
					FileSystem: d.FileSystem,
				})
			}
			mu.Unlock()
		}
	}()

	// Fetch Network
	wg.Add(1)
	go func() {
		defer wg.Done()
		var dst []Win32_NetworkAdapterConfiguration
		q := "SELECT Description, IPAddress, DNSServerSearchOrder, IPEnabled FROM Win32_NetworkAdapterConfiguration WHERE IPEnabled = TRUE"
		if err := wmi.Query(q, &dst); err == nil {
			mu.Lock()
			for _, n := range dst {
				info.Network = append(info.Network, NetInfo{
					Description: n.Description,
					IPs:         n.IPAddress,
					DNS:         n.DNSServerSearchOrder,
				})
			}
			mu.Unlock()
		}
	}()

	// Fetch AV/Firewall
	wg.Add(1)
	go func() {
		defer wg.Done()
		var avDst []AntiVirusProduct
		_ = wmi.QueryNamespace("SELECT DisplayName FROM AntiVirusProduct", &avDst, "root\\SecurityCenter2")
		mu.Lock()
		for _, av := range avDst {
			info.Antivirus = append(info.Antivirus, av.DisplayName)
		}
		mu.Unlock()

		var fwDst []FirewallProduct
		_ = wmi.QueryNamespace("SELECT DisplayName FROM FirewallProduct", &fwDst, "root\\SecurityCenter2")
		mu.Lock()
		for _, fw := range fwDst {
			info.Firewall = append(info.Firewall, fw.DisplayName)
		}
		mu.Unlock()
	}()

	wg.Wait()

	var buf bytes.Buffer
	buf.WriteString("=== System Information ===\n")
	buf.WriteString(fmt.Sprintf("OS: %s (Version: %s, Build: %s)\n", info.OS.Name, info.OS.Version, info.OS.Build))
	buf.WriteString(fmt.Sprintf("CPU: %s (%d Sockets, %d Cores/Socket, %d Total Cores, %d MHz)\n", 
		info.CPU.Name, info.CPU.Sockets, info.CPU.CoresPerSocket, info.CPU.TotalCores, info.CPU.Speed))
	buf.WriteString(fmt.Sprintf("RAM: Total: %.2f GB, Free: %.2f GB, Load: %d%%\n", info.Memory.PhysicalTotalGB, info.Memory.FreePhysicalGB, info.Memory.LoadPercentage))
	buf.WriteString(fmt.Sprintf("Virtual RAM: Total: %.2f GB, Free: %.2f GB\n", info.Memory.TotalVirtualGB, info.Memory.FreeVirtualGB))

	buf.WriteString("\nDisks:\n")
	for _, d := range info.Disks {
		buf.WriteString(fmt.Sprintf("- %s [%s] Total: %.2f GB, Free: %.2f GB, Used: %.2f GB\n", d.Drive, d.FileSystem, d.TotalGB, d.FreeGB, d.UsedGB))
	}

	buf.WriteString("\nAntivirus:\n")
	if len(info.Antivirus) == 0 {
		buf.WriteString("- None found\n")
	} else {
		for _, av := range info.Antivirus {
			buf.WriteString(fmt.Sprintf("- %s\n", av))
		}
	}

	buf.WriteString("\nFirewall:\n")
	if len(info.Firewall) == 0 {
		buf.WriteString("- None found\n")
	} else {
		for _, fw := range info.Firewall {
			buf.WriteString(fmt.Sprintf("- %s\n", fw))
		}
	}

	buf.WriteString("\nNetwork:\n")
	for _, n := range info.Network {
		buf.WriteString(fmt.Sprintf("- %s\n  IPs: %v\n  DNS: %v\n", n.Description, n.IPs, n.DNS))
	}

	fmt.Print(buf.String())

	jsonData, _ := json.MarshalIndent(info, "", "  ")
	_ = os.WriteFile(filepath.Join(outputDir, "system_info.json"), jsonData, 0644)
	_ = os.WriteFile(filepath.Join(outputDir, "system_info.txt"), buf.Bytes(), 0644)

	fmt.Printf("\nData saved to %s/system_info.json and %s/system_info.txt\n", outputDir, outputDir)
}
