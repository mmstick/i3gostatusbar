package system

import (
	"os"
	"os/user"
	"syscall"
	"time"
)

// Contains system information
type Info struct {
	Kernel      string
	Model       string
	Host        string
	User        string
	Uptime      string
	Cpufreqs    string
	MemStat     string
	NetStat     string
	NetName     string
	NetSpeed    string
	BatteryInfo string
	Date        string
	MemTotal    uint
}

// This is used to get a string from the kernel utsname
func utsnameToString(unameArray [65]int8) string {
	var byteString [65]byte
	var indexLength int
	for ; unameArray[indexLength] != 0; indexLength++ {
		byteString[indexLength] = uint8(unameArray[indexLength])
	}
	return string(byteString[0:indexLength])
}

// Returns kernel version information
func KernelVersion() string {
	var utsname syscall.Utsname
	_ = syscall.Uname(&utsname)
	return utsnameToString(utsname.Release)
}

// Returns the current time
func CurrentTime(date *string, done chan bool) {
	*date = time.Now().Format(time.RFC1123)
	done <- true
}

// Returns the hostname
func Host() string {
	host, _ := os.Hostname()
	return host
}

// Returns the user's name
func Username() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}
