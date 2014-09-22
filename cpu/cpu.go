package cpu

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var sprintf = fmt.Sprintf

// Returns a newline-delimited string slice of the file
func parseFile() []string {
	cached, _ := ioutil.ReadFile("/proc/cpuinfo")
	return strings.Split(string(cached), "\n")
}

// Returns the number of CPU cores in the system
func parseCPUCount(count string) int {
	cores, _ := strconv.Atoi(count[11:])
	return cores
}

// Parses CPU frequency information
func parseFrequency(frequency string) string {
	frequency = frequency[11 : len(frequency)-4]
	if len(frequency) < 4 {
		frequency = " " + frequency
	}
	return frequency + "MHz"
}

// Obtains the printf format for the current core frequency
func getFrequencyFormat(index *int, lastCore int) string {
	if *index == lastCore {
		return "%s"
	} else {
		return "%s "
	}
}

// Returns the frequency of the current core
func getCoreFrequency(cpuInfo []string, index *int, format *string) string {
	return sprintf(*format, parseFrequency(cpuInfo[*index*27+7]))
}

// Returns a string containing the frequencies of each CPU core
func Frequencies(cpufreqs *string, done chan bool) {
	cpuInfo := parseFile()
	numCPUs := parseCPUCount(cpuInfo[len(cpuInfo)-17]) + 1
	var cpuFrequencies string
	for index := 0; index < numCPUs; index++ {
		format := getFrequencyFormat(&index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, &index, &format)
	}
	*cpufreqs = "Cores:" + cpuFrequencies
	done <- true
}

// Returns the CPU Model
func Model() string {
	return parseFile()[4][13:]
}
