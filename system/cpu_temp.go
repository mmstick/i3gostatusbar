package system

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

/* getHwmon searches for a hwmon entry that has the 'k10temp' name and returns
 * it. It currently only supports finding AMD CPU temperatures since I don't know
 * what Intel calls their temperature sensor. */
func getHwmon() string {
	hwmondir, _ := ioutil.ReadDir("/sys/class/hwmon/")
	for index := range hwmondir {
		dir := fmt.Sprintf("/sys/class/hwmon/hwmon%d/device/", index)
		name, err := ioutil.ReadFile(dir + "name")
		if err == nil {
			if string(name)[:3] == "k10" {
				return dir + "temp1_input"
			}
		}
	}
	return "err"	
}

/* hwmonCheck will return 0 if it could not find the correct hwmon entry, else it
 * will return an integer containing the value of the CPU temperature. */
func hwmonCheck(input string, err *error) int {
	var temp int
	if *err == nil {
		temp, _ = strconv.Atoi(input[:2])
	} else {
		temp = 0
	}
	return temp
}

// getTemperature returns the CPU temperature as indicated by 'hwmon'. 
func getTemperature() int {
	temp, err := ioutil.ReadFile(getHwmon())
	return hwmonCheck(string(temp), &err)
}
