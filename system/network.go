// This file contains functions for obtaining information about the currently
// active network, such as the connection name, speed, and total up/downloaded
// bytes since boot.
package system

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// parseFile returns a newline-delimited string slice of the file.
func parseFile(file string) []string {
	cached, _ := ioutil.ReadFile(file)
	return strings.Split(string(cached), "\n")
}

// connectionIsNotLoopback returns true if the connection is not a loopback
// address.
func connectionIsNotLoopback(connection *string) bool {
	if *connection == "lo" {
		return false
	} else {
		return true
	}
}

// connectionIsUp returns true if the connection status is 'up'.
func connectionIsUp(connection *string) bool {
	if parseFile("/sys/class/net/" + *connection + "/operstate")[0][0] == 'u' {
		return true
	} else {
		return false
	}
}

// connectionIsActive returns true if the connection is active and isn't a
// loopback address.
func connectionIsActive(connection string) bool {
	if connectionIsNotLoopback(&connection) && connectionIsUp(&connection) {
		return true
	} else {
		return false
	}
}

// Returns the currently active connection name.
func getCurrentNetwork() string {
	var activeConnection string
	networkConnections, _ := ioutil.ReadDir("/sys/class/net/")
	for _, connection := range networkConnections {
		if connectionIsActive(connection.Name()) {
			activeConnection = connection.Name()
		}
	}
	return activeConnection
}

// PadDigits pads digits with spaces so that the status bar always has the same
// number of characters.
func padDigits(number int) string {
	numberString := fmt.Sprintf("%d", number)
	switch len(numberString) {
	case 1:
		return "   " + numberString
	case 2:
		return "  " + numberString
	case 3:
		return " " + numberString
	default:
		return numberString
	}
}

// formatBytes formats the bytes into their respectiev scales.
func formatBytes(bytes int) string {
	switch {
	case bytes > 10737418240:
		return fmt.Sprintf("%sGiB", padDigits(bytes/1073741824))
	case bytes > 10485760:
		return fmt.Sprintf("%sMiB", padDigits(bytes/1048576))
	case bytes > 10240:
		return fmt.Sprintf("%sKiB", padDigits(bytes/1024))
	default:
		return fmt.Sprintf("%s  B", padDigits(bytes))
	}
}

// networkDir returns the transfer statistics directory name.
func networkDir() string {
	return fmt.Sprintf("/sys/class/net/%s/", getCurrentNetwork())
}

// fileAsInt returns the contents of the file as an integer variable.
func fileAsInt(file string) int {
	return strToInt(parseFile(file)[0])
}

// downloadInformation returns the amount of bytes downloaded since boot.
func downloadInformation() string {
	return formatBytes(fileAsInt(networkDir() + "statistics/rx_bytes"))
}

// uploadInformation returns the amount of bytes uploaded since boot.
func uploadInformation() string {
	return formatBytes(fileAsInt(networkDir() + "statistics/tx_bytes"))
}

// NetStats returns RX/TX transfer statistics since boot.
func NetStats(transferStat *string, done chan bool) {
	*transferStat = fmt.Sprintf("D:%s U:%s", downloadInformation(),
		uploadInformation())
	done <- true
}

// Speed returns the speed of the currently active connection in Mbps.
func NetworkSpeed() string {
	return "S: " + parseFile(networkDir() + "speed")[0] + " Mbps"
}

// Name returns the currently active network connection.
func NetworkName() string {
	return getCurrentNetwork()
}
