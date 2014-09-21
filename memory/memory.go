package memory

import (
        "fmt"
        "io/ioutil"
        "strconv"
        "strings"
)

var sprintf = fmt.Sprintf

// Returns a newline-delimited string slice of the file
func parseFile(file string) []string {
        cached, _ := ioutil.ReadFile(file)
        return strings.Split(string(cached), "\n")
}

// Returns the amount of memory installed in the system
func Installed() int {
        mem := parseFile("/proc/meminfo")[0]
        totalMem, _ := strconv.Atoi(mem[17 : len(mem)-3])
        return totalMem / 1024
}

// Returns the memory currently available as an int
func memAvailable() int {
        mem := parseFile("/proc/meminfo")[2]
        memAvailable, _ := strconv.Atoi(mem[18 : len(mem)-3])
        return memAvailable / 1024
}

// Returns a string indicating memory usage out of total available memory
func Statistics(memTotal *int, memStat *string, done chan bool) {
        *memStat = sprintf("RAM: %d/%dMB", *memTotal-memAvailable(), *memTotal)
        done <- true
}
