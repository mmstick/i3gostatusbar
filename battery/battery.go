package battery

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var sprintf = fmt.Sprintf

// Exists checks if a battery exists in the current system, and if so returns
// true as well as the number of jobs to process concurrently.
func Exists() (bool, int) {
	battery, _ := ioutil.ReadDir("/sys/class/power_supply")
	if len(battery) > 0 {
		return true, 7
	} else {
		return false, 6
	}
}

// parseFile returns a newline-delimited string of the file contents.
func parseFile(file string) string {
	cached, _ := ioutil.ReadFile(file)
	return strings.TrimSuffix(string(cached), "\n")
}

// status returns the battery satus: [C]harging, [F]ull, or [D]ischarging.
func status() byte {
	return parseFile("/sys/class/power_supply/BAT1/status")[0]
}

// charge returns the current battery life in percent.
func charge() string {
	return parseFile("/sys/class/power_supply/BAT1/capacity")
}

// isCharging returns a string indicating that the battery is currently
// charging.
func isCharging() string {
	return sprintf("BAT Charging: %s%%", charge())
}

// isDischarging returns a string indicating that the battery is currently
// discharging.
func isDischarging() string {
	return sprintf("BAT: %s%%", charge())
}

// Information checks the status of the battery and returns information
// regarding that status.
func Information(batteryStat *string, done chan bool) {
	switch status() {
	case 'C':
		*batteryStat = isCharging()
	case 'F':
		*batteryStat = "BAT Full"
	default:
		*batteryStat = isDischarging()
	}
	done <- true
}
