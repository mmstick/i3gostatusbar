package main

import (
	"github.com/mmstick/i3gostatusbar/battery"
	"github.com/mmstick/i3gostatusbar/cpu"
	"github.com/mmstick/i3gostatusbar/memory"
	"github.com/mmstick/i3gostatusbar/network"
	"github.com/mmstick/i3gostatusbar/system"
	"github.com/mmstick/i3gostatusbar/uptime"
	"fmt"
	"time"
)

var refreshRate = 1 * time.Second

// statusBarFormat returns the printf format for formatting the status bar.
// If a battery exists, this will return an additional string element.
func statusBarFormat(batteryExists bool) string {
	if batteryExists {
		return "%s@%s | %s | %s | %s %s | %s | Temp: %sC | %s %s %s | %s | %s\n"
	} else {
		return "%s@%s | %s | %s | %s %s | %s | Temp: %sC | %s %s %s | %s\n"
	}
}

// getStaticSystemInformation collects system information that does not need
// to be dynamically updated.
func getStaticSystemInformation() *system.Info {
	info := system.Info{
		Kernel:   system.KernelVersion(),
		Model:    cpu.Model(),
		MemTotal: memory.Installed(),
		Host:     system.Host(),
		User:     system.Username(),
		NetName:  network.Name(),
		NetSpeed: network.Speed(),
	}
	return &info
}

// getDynamicSystemInformation adds dynamic system information to the
// systemInfo struct and checks for the existence of a battery.
func getDynamicSystemInformation(info system.Info) (*system.Info, bool) {
	batteryExists, jobs := battery.Exists()
	synchronize := make(chan bool, jobs)
	if batteryExists {
		go battery.Information(&info.Battery, synchronize)
	}
	go uptime.Get(&info.Uptime, synchronize)
	go cpu.Frequencies(&info.Cpufreqs, synchronize)
	go cpu.CPUTemp(&info.Cputemp, synchronize)
	go memory.Statistics(&info.MemTotal, &info.Memory, synchronize)
	go network.Statistics(&info.NetStat, synchronize)
	go system.CurrentTime(&info.Date, synchronize)
	for jobCount := 0; jobCount < jobs; jobCount++ {
		<-synchronize
	}
	return &info, batteryExists
}

// refreshBar refreshes the status bar
func refreshBar(info *system.Info, batteryExists bool) {
	if batteryExists {
		fmt.Printf(statusBarFormat(true), info.User, info.Host,
			info.Kernel, info.Uptime, info.Model, info.Cpufreqs,
			info.Cputemp, info.Memory, info.NetName, info.NetSpeed,
			info.NetStat, info.Battery, info.Date)
	} else {
		fmt.Printf(statusBarFormat(false), info.User, info.Host,
			info.Kernel, info.Uptime, info.Model, info.Cpufreqs,
			info.Cputemp, info.Memory, info.NetName, info.NetSpeed,
			info.NetStat, info.Date)
	}
}

func main() {
	system := getStaticSystemInformation()
	for {
		go refreshBar(getDynamicSystemInformation(*system))
		time.Sleep(refreshRate)
	}
}
