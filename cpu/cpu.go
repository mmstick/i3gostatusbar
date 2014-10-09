// Packge cpu contains functions for obtaining information about the cetral
// processing unit.
package cpu

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var sprintf = fmt.Sprintf

// parseCPUInfo returns a newline-delimited string slice of '/proc/cpuinfo'.
func parseCPUInfo() []string {
	cached, _ := ioutil.ReadFile("/proc/cpuinfo")
	return strings.Split(string(cached), "\n")
}

// parseCPUCount returns the number of CPU cores in the system.
func parseCPUCount(cpuinfo []string) int {
	cores, _ := strconv.Atoi(strings.Fields(cpuinfo[len(cpuinfo) - 18])[3])
	return cores + 1
}

/* parseFrequency obtains the CPU frequency. It first uses Fields to return only
 * the field on that line that contains the CPU frequency. Then, it uses Replace
 * to replace the periods (for three digit frequencies) with a space. Finally,
 * it returns the frequency and adds 'MHz' at the end. */
func parseFrequency(frequency string) string {
	frequency = strings.Fields(frequency)[3][:4]
	frequency = strings.Replace(frequency, ".", " ", -1)
	return frequency + " MHz"
}

// getFrequencyFormat returns the printf format for the current core frequency.
func getFrequencyFormat(index *int, lastCore int) string {
	if *index == lastCore {
		return "%s"
	} else {
		return "%s "
	}
}

// getCoreFrequency returns the frequency of the current core.
func getCoreFrequency(cpuInfo []string, index *int, format *string) string {
	return sprintf(*format, parseFrequency(cpuInfo[*index*28+7]))
}

// Frequencies returns a string containing the frequencies of each CPU core.
func Frequencies(cpufreqs *string, done chan bool) {
	cpuInfo := parseCPUInfo()
	numCPUs := parseCPUCount(cpuInfo)
	var cpuFrequencies string
	for index := 0; index < numCPUs; index++ {
		format := getFrequencyFormat(&index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, &index, &format)
	}
	*cpufreqs = "Cores:" + cpuFrequencies
	done <- true
}

// getTemperature returns the CPU temperature as indicated by 'hwmon'.
func getTemperature() int {
	input, err := ioutil.ReadFile("/sys/class/hwmon/hwmon0/temp1_input")
	var temp int
	if err == nil {
		temp, _ = strconv.Atoi(string(input)[:2])
	} else {
		temp = 0
	}
	return temp
}

// CPUTemp sets the temperature of the CPU.
func CPUTemp(cputemp *int, done chan bool) {
	*cputemp = getTemperature()
	done <- true
}

// Model returns the CPU Model
func Model() string {
	modelinfo := strings.Fields(parseCPUInfo()[4])[3:]
	vendor := modelinfo[0]
	model := modelinfo[1]
	return vendor + " " + model
}
