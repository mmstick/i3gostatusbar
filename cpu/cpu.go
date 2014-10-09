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

/* getHwmon searches for a hwmon entry that has the 'k10temp' name and returns
 * it. It currently only supports finding AMD CPU temperatures since I don't know
 * what Intel calls their temperature sensor. */
func getHwmon() string {
	hwmondir, _ := ioutil.ReadDir("/sys/class/hwmon/")
	for index := range hwmondir {
		dir := fmt.Sprintf("/sys/class/hwmon/hwmon%d/device/", index)
		name, err := ioutil.ReadFile(dir + "name")
		if err == nil {
			if string(name)[:3] == "k10" {
				return dir + "temp1_input"
			}
		}
	}
	return "err"	
}

/* hwmonCheck will return 0 if it could not find the correct hwmon entry, else it
 * will return an integer containing the value of the CPU temperature. */
func hwmonCheck(input string, err *error) int {
	var temp int
	if *err == nil {
		temp, _ = strconv.Atoi(input[:2])
	} else {
		temp = 0
	}
	return temp
}

// getTemperature returns the CPU temperature as indicated by 'hwmon'. 
func getTemperature() int {
	temp, err := ioutil.ReadFile(getHwmon())
	return hwmonCheck(string(temp), &err)
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
