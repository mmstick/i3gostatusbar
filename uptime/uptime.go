package uptime

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

// Converts a string that we know is an integer to an integer
func strToInt(input string) int {
        output, _ := strconv.Atoi(input)
        return output
}

// Takes an uptime value and divides it by a scale (days/hours/minutes).
// After determining the amount of time in that scale, it subtracts that amount
// from uptime and returns the time in days/hours/minutes
func getTimeScale(time **int, scale int) int {
        var output int
        if **time > scale {
                output = **time / scale
                **time -= output * scale
        }
        return output
}

// Adds an extra zero in case the time only has one digit
func formatTime(time int) string {
        output := strconv.Itoa(time)
        if len(output) == 1 {
                output = "0" + output
        }
        return output
}

// Returns the time formatted in days and subtracts that from time.
func getDays(time *int) string {
        return formatTime(getTimeScale(&time, 86400))
}

// Returns the time formatted in hours and subtracts that from time.
func getHours(time *int) string {
        return formatTime(getTimeScale(&time, 3600))
}

// Returns the time formatted in minutes and subtracts that from time.
func getMinutes(time *int) string {
        return formatTime(getTimeScale(&time, 60))
}

// Returns the time formatted in seconds.
func getSeconds(time *int) string {
        return formatTime(*time)
}

// Takes the uptime integer and converts it into a human readable format.
// Ex: 01:21:18:57 for days:hours:seconds:minutes
func humanReadableTime(time int) string {
        return sprintf("%s:%s:%s:%s", getDays(&time), getHours(&time),
                getMinutes(&time), getSeconds(&time))
}

// Open /proc/uptime and return the value in integer format.
func parseUptime() int {
        parsedInfo := strings.SplitAfter(parseFile("/proc/uptime")[0], " ")
        return strToInt(parsedInfo[0][0 : len(parsedInfo[0])-4])
}

// Returns the current uptime in days:hours:minutes:seconds format
func Get(uptime *string, done chan bool) {
        *uptime = sprintf("Uptime: %s", humanReadableTime(parseUptime()))
        done <- true
}
