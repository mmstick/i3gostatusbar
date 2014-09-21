package main

import (
	"./battery"
	"./cpu"
	"fmt"
	"./memory"
	"./network"
	"runtime"
	"./system"
	"time"
	"./uptime"
)

var refreshRate = 1 * time.Second

// Returns the printf format for formatting the status bar. If a battery exists,
// this will return an additional string element.
func statusBarFormat(batteryExists bool) string {
	if batteryExists {
		return "%s@%s | %s | %s | %s %s | %s | %s | %s | %s\n"
	} else {
		return "%s@%s | %s | %s | %s %s | %s | %s | %s\n"
	}
}

// Returns system information that does not need to be dynamically updated.
func getStaticSystemInformation() *system.Info {
	info := system.Info{
		Kernel: system.KernelVersion(),
		Model: cpu.Model(),
		MemTotal: memory.Installed(),
		Host: system.Host(),
		User: system.Username(),
	}
	return &info
}

// Adds dynamic system information to the systemInfo struct and checks for
// the existence of a battery.
func getDynamicSystemInformation(info system.Info) (*system.Info, bool) {
	batteryExists, jobs := battery.Exists()
	synchronize := make(chan bool, jobs)
	if batteryExists {
		go battery.Information(&info.BatteryInfo, synchronize)
	}
	go uptime.Get(&info.Uptime, synchronize)
	go cpu.Frequencies(&info.Cpufreqs, synchronize)
	go memory.Statistics(&info.MemTotal, &info.MemStat, synchronize)
	go network.Statistics(&info.TransferStat, synchronize)
	go system.CurrentTime(&info.Date, synchronize)
	for jobCount := 0; jobCount < jobs; jobCount++ {
		<-synchronize
	}
	return &info, batteryExists
}

// Refreshes the status bar
func refreshBar(info *system.Info, batteryExists bool) {
	if batteryExists {
		fmt.Printf(statusBarFormat(true), info.User, info.Host, info.Kernel,
			info.Uptime, info.Model, info.Cpufreqs, info.MemStat,
			info.TransferStat, info.BatteryInfo, info.Date)
	} else {
		fmt.Printf(statusBarFormat(false), info.User, info.Host, info.Kernel,
			info.Uptime, info.Model, info.Cpufreqs, info.MemStat,
			info.TransferStat, info.Date)
	}
}

func main() {
	runtime.GOMAXPROCS(6)
	system := getStaticSystemInformation()
	for {
		go refreshBar(getDynamicSystemInformation(*system))
		time.Sleep(refreshRate)
	}
}
