// Package system contains the Info struct that stores all system information
// as well as a few functions for gathering system information.
package system

import (
	"os"
	"os/user"
	"syscall"
	"time"
)

// Info contains system information
type Info struct {
	Kernel      string
	Model       string
	Host        string
	User        string
	Uptime      string
	Cpufreqs    string
	Cputemp     int
	Memory      string
	NetStat     string
	NetName     string
	NetSpeed    string
	Battery     string
	Date        string
	MemTotal    uint
}

// utsnameToString is used to get a string from the kernel utsname.
func utsnameToString(unameArray [65]int8) string {
	var byteString [65]byte
	var indexLength int
	for ; unameArray[indexLength] != 0; indexLength++ {
		byteString[indexLength] = uint8(unameArray[indexLength])
	}
	return string(byteString[0:indexLength])
}

// KernelVersion returns kernel version information
func KernelVersion() string {
	var utsname syscall.Utsname
	_ = syscall.Uname(&utsname)
	return utsnameToString(utsname.Release)
}

// CurrentTime returns the current time
func CurrentTime(date *string, done chan bool) {
	*date = time.Now().Format(time.RFC1123)
	done <- true
}

// Host returns the hostname
func Host() string {
	host, _ := os.Hostname()
	return host
}

// Username returns the user's name
func Username() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}
