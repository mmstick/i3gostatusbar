// This file belongs to the cpu group; gathers information about cpu frequencies.
package system

import "fmt"
import "strings"

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
func getFrequencyFormat(index int, lastCore int) string {
	if index == lastCore {
		return "%s"
	} else {
		return "%s "
	}
}

// getCoreFrequency returns the frequency of the current core.
func getCoreFrequency(cpuInfo []string, index int, format string) string {
	return fmt.Sprintf(format, parseFrequency(cpuInfo[(index*28)+7]))
}

// getFrequencyString returns all core frequencies as a string.
func getFrequencyString() string {
	cpuInfo := parseCPUInfo()
	numCPUs := parseCPUCount(cpuInfo)
	var cpuFrequencies string
	index := 0
	switch numCPUs % 4 {
	case 1:
		format := getFrequencyFormat(index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index, format)
		index++
	case 2:
		format := getFrequencyFormat(index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index, format)
		format = getFrequencyFormat(index+1, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+1, format)
		index += 2
	case 3:
		format := getFrequencyFormat(index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index, format)
		format = getFrequencyFormat(index+1, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+1, format)
		format = getFrequencyFormat(index+2, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+2, format)
		index += 3
	}
	for ; index != numCPUs; index += 4 {
		format := getFrequencyFormat(index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index, format)
		format = getFrequencyFormat(index+1, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+1, format)
		format = getFrequencyFormat(index+2, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+2, format)
		format = getFrequencyFormat(index+3, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, index+3, format)
	}
	return cpuFrequencies
}
