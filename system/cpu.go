// This is the main cpu file -- it gathers all CPU-related information.
package system

import (
	"io/ioutil"
	"strings"
)

// parseCPUInfo returns a newline-delimited string slice of '/proc/cpuinfo'.
func parseCPUInfo() []string {
	cached, _ := ioutil.ReadFile("/proc/cpuinfo")
	return strings.Split(string(cached), "\n")
}

// parseCPUCount returns the number of CPU cores in the system.
func parseCPUCount(cpuinfo []string) int {
	return strToInt(strings.Fields(cpuinfo[len(cpuinfo) - 18])[3]) + 1
}

// Model returns the CPU Model
func CPUModel() string {
	modelinfo := strings.Fields(parseCPUInfo()[4])[3:]
	return modelinfo[0] + " " + modelinfo[1]
}

// CPUTemp sets the temperature of the CPU.
func CPUTemp(cputemp *int, done chan bool) {
	*cputemp = getTemperature()
	done <- true
}

// Frequencies sets '*cpufreq' with a string containing all core frequencies.
func CPUFrequencies(cpufreqs *string, done chan bool) {
	*cpufreqs = "Cores:" + getFrequencyString()
	done <- true
}
