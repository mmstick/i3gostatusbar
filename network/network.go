package network

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

// Returns true if the connection is not a loopback address.
func connectionIsNotLoopback(connection *string) bool {
	if *connection == "lo" {
		return false
	} else {
		return true
	}
}

// Returns true if the connection is up
func connectionIsUp(connection *string) bool {
	if parseFile("/sys/class/net/" + *connection + "/operstate")[0][0] == 'u' {
		return true
	} else {
		return false
	}
}

// Returns true if the connection is active and isn't a loopback address.
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

// Pads digits with spaces so that the status bar always has the same number of
// characters.
func padDigits(number int) string {
	numberString := sprintf("%d", number)
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

// Formats the bytes into their respectiev scales
func formatBytes(bytes int) string {
	switch {
	case bytes > 10737418240:
		return sprintf("%sGiB", padDigits(bytes/1073741824))
	case bytes > 10485760:
		return sprintf("%sMiB", padDigits(bytes/1048576))
	case bytes > 10240:
		return sprintf("%sKiB", padDigits(bytes/1024))
	default:
		return sprintf("%s  B", padDigits(bytes))
	}
}

// Returns the transfer statistics directory name.
func networkDir() string {
	return sprintf("/sys/class/net/%s/", getCurrentNetwork())
}

// Returns the contents of the file as an integer variable
func fileAsInt(file string) int {
	return strToInt(parseFile(file)[0])
}

// Returns the amount of bytes downloaded since boot.
func downloadInformation() string {
	return formatBytes(fileAsInt(networkDir() + "statistics/rx_bytes"))
}

// Returns the amount of bytes uploaded since boot.
func uploadInformation() string {
	return formatBytes(fileAsInt(networkDir() + "statistics/tx_bytes"))
}

// Returns RX/TX transfer statistics
func Statistics(transferStat *string, done chan bool) {
	*transferStat = sprintf("D:%s U:%s", downloadInformation(),
		uploadInformation())
	done <- true
}

// Returns the speed of the currently active connection in Mbps
func Speed() string {
	return "S: " + parseFile(networkDir() + "speed")[0] + " Mbps"
}

// Returns the currently active network connection.
func Name() string {
	return getCurrentNetwork()
}
