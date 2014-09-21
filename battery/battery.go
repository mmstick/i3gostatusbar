package battery

import (
        "fmt"
        "io/ioutil"
        "strings"
)

var sprintf = fmt.Sprintf

// If a battery exists in the system, return true and the number of concurrent
// jobs to perform.
func Exists() (bool, int) {
        battery, _ := ioutil.ReadDir("/sys/class/power_supply")
        if len(battery) > 0 {
                return true, 6
        } else {
                return false, 5
        }
}

// Returns a newline-delimited string slice of the file
func parseFile(file string) []string {
        cached, _ := ioutil.ReadFile(file)
        return strings.Split(string(cached), "\n")
}

// Returns the status of the battery: [C]harging, [F]ull, or [D]ischarging.
func status() byte {
        return parseFile("/sys/class/power_supply/BAT1/status")[0][0]
}

// Returns the current battery life in percent.
func charge() string {
        return parseFile("/sys/class/power_supply/BAT1/capacity")[0]
}

// Returns a string indicating that the battery is currently charging.
func isCharging() string {
        return sprintf("BAT Charging: %s%%", charge())
}

// Returns a string indicating that the battery is currently discharging.
func isDischarging() string {
        return sprintf("BAT: %s%%", charge())
}

// Checks the battery status and returns information based on that information.
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
