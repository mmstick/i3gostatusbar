package memory

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	sprintf = fmt.Sprintf
	fields  = strings.Fields
	split   = strings.Split
	toInt   = strconv.Atoi
)

// Opens meminfo and splits a specific line from that file, returning the field as
// a value representing MiB.
func parseMem(line uint) uint {
	cached, _ := ioutil.ReadFile("/proc/meminfo")
	memory, _ := toInt(fields(split(string(cached), "\n")[line])[1])
	return uint(memory / 1024)
}

// Returns the amount of memory installed in the system
func Installed() uint {
	return parseMem(0)
}

// Returns the memory currently available as an int
func memAvailable() uint {
	return parseMem(2)
}

// Returns a string indicating memory usage out of total available memory
func Statistics(memTotal *uint, memStat *string, done chan bool) {
	*memStat = sprintf("RAM: %d/%dMB", *memTotal-memAvailable(), *memTotal)
	done <- true
}
