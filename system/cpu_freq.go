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
func getFrequencyFormat(index *int, lastCore int) string {
	if *index == lastCore {
		return "%s"
	} else {
		return "%s "
	}
}

// getCoreFrequency returns the frequency of the current core.
func getCoreFrequency(cpuInfo []string, index *int, format *string) string {
	return fmt.Sprintf(*format, parseFrequency(cpuInfo[*index*28+7]))
}

// getFrequencyString returns all core frequencies as a string.
func getFrequencyString() string {
	cpuInfo := parseCPUInfo()
	numCPUs := parseCPUCount(cpuInfo)
	var cpuFrequencies string
	for index := 0; index < numCPUs; index++ {
		format := getFrequencyFormat(&index, numCPUs-1)
		cpuFrequencies += getCoreFrequency(cpuInfo, &index, &format)
	}
	return cpuFrequencies
}
