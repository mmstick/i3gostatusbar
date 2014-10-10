// This file contains functions for gathering battery information.
package system

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// BatteryExists checks if a battery exists in the current system, and if so returns
// true as well as the number of jobs to process concurrently.
func BatteryExists() (bool, int) {
	battery, _ := ioutil.ReadDir("/sys/class/power_supply")
	if len(battery) > 0 {
		return true, 7
	} else {
		return false, 6
	}
}

// parseLine returns a single line string of the file contents.
func parseLine(file string) string {
	cached, _ := ioutil.ReadFile(file)
	return strings.TrimSuffix(string(cached), "\n")
}

// batteryStatus returns the battery satus: [C]harging, [F]ull, or [D]ischarging.
func batteryStatus() byte {
	return parseLine("/sys/class/power_supply/BAT1/status")[0]
}

// batteryCharge returns the current battery life in percent.
func batteryCharge() string {
	return parseLine("/sys/class/power_supply/BAT1/capacity")
}

// isCharging returns a string indicating that the battery is currently
// charging.
func isCharging() string {
	return fmt.Sprintf("BAT Charging: %s%%", batteryCharge())
}

// isDischarging returns a string indicating that the battery is currently
// discharging.
func isDischarging() string {
	return fmt.Sprintf("BAT: %s%%", batteryCharge())
}

// BatteryInfo checks the status of the battery and returns information
// regarding that status.
func BatteryInfo(batteryStat *string, done chan bool) {
	switch batteryStatus() {
	case 'C':
		*batteryStat = isCharging()
	case 'F':
		*batteryStat = "BAT Full"
	default:
		*batteryStat = isDischarging()
	}
	done <- true
}
