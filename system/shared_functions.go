// This file contains functions shared between all files.
package system

import "strconv"

// strToInt converts a string that we know is an integer to an integer
func strToInt(input string) int {
	output, _ := strconv.Atoi(input)
	return output
}
