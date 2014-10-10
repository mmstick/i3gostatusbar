// This file contains functions for obtaining information about memory.
package system

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	fields   = strings.Fields
	split   = strings.Split
)

// parseMem opens meminfo and splits a specific line from that file, returning
// the field as a value representing MiB.
func parseMem(line uint) uint {
	cached, _ := ioutil.ReadFile("/proc/meminfo")
	memory := strToInt(fields(split(string(cached), "\n")[line])[1])
	return uint(memory / 1024)
}

// TotalMem returns the amount of memory installed in the system.
func TotalMem() uint {
	return parseMem(0)
}

// memAvailable returns the memory currently available as an int.
func memAvailable() uint {
	return parseMem(2)
}

// memUsed returns the memory used by subtracting available from the total.
func memUsed(memTotal *uint) uint {
	return *memTotal - memAvailable()
}

// MemStats returns a string indicating memory usage out of total available
// memory.
func MemStats(memTotal *uint, memStat *string, done chan bool) {
	*memStat = fmt.Sprintf("RAM: %d/%dMiB", memUsed(memTotal), *memTotal)
	done <- true
}
