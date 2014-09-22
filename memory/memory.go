package memory

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var sprintf = fmt.Sprintf
var split = strings.Split
var toInt = strconv.Atoi

// Opens meminfo, splits it into a newline-delimited slice, cuts a single line
// from that slice, splits that line into a space-delimited slice, gets a
// single element from that slie and returns the value in MiB.
func parseMem(line, element uint) uint {
	cached, _ := ioutil.ReadFile("/proc/meminfo")
	mem, _ := toInt(split(split(string(cached), "\n")[line], " ")[element])
	return uint(mem / 1024)
}

// Returns the amount of memory installed in the system
func Installed() uint {
	return parseMem(0, 8)
}

// Returns the memory currently available as an int
func memAvailable() uint {
	return parseMem(2, 4)
}

// Returns a string indicating memory usage out of total available memory
func Statistics(memTotal *uint, memStat *string, done chan bool) {
	*memStat = sprintf("RAM: %d/%dMB", *memTotal-memAvailable(), *memTotal)
	done <- true
}
