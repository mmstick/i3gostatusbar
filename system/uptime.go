// This file contains functions related to processing system uptime information
package system

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// getTimeScale tkes an uptime value and divides it by a scale
// (days/hours/minutes). After determining the amount of time in that scale,
// it subtracts that amount from uptime and returns the time in
// days/hours/minutes.
func getTimeScale(time **int, scale int) int {
	var output int
	if **time > scale {
		output = **time / scale
		**time -= output * scale
	}
	return output
}

// padTime an extra zero in case the time only has one digit
func padTime(time int) string {
	output := strconv.Itoa(time)
	if len(output) == 1 {
		output = "0" + output
	}
	return output
}

// getDays returns the time formatted in days and subtracts that from time.
func getDays(time *int) string {
	return padTime(getTimeScale(&time, 86400))
}

// getHours returns the time formatted in hours and subtracts that from time.
func getHours(time *int) string {
	return padTime(getTimeScale(&time, 3600))
}

// getMinutes returns the time formatted in minutes and subtracts that from time.
func getMinutes(time *int) string {
	return padTime(getTimeScale(&time, 60))
}

// getSeconds returns the time formatted in seconds.
func getSeconds(time *int) string {
	return padTime(*time)
}

// humanReadableTime takes the uptime integer and converts it into a human
// readable format. Ex: 01:21:18:57 for days:hours:seconds:minutes
func humanReadableTime(time int) string {
	return fmt.Sprintf("%s:%s:%s:%s", getDays(&time), getHours(&time),
		getMinutes(&time), getSeconds(&time))
}

// parseUptime opens /proc/uptime and return the value in integer format.
func parseUptime() int {
	cached, _ := ioutil.ReadFile("/proc/uptime")
	return strToInt(strings.Split(string(cached), ".")[0])
}

// Uptime returns the current uptime in days:hours:minutes:seconds format
func Uptime(uptime *string, done chan bool) {
	*uptime = fmt.Sprintf("Uptime: %s", humanReadableTime(parseUptime()))
	done <- true
}
