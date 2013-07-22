package dateparser

func findPeriod(input string) (pos int) {
	for i, char := range input {
		if char == '.' {
			return i
		}
	}
	
	return -1
}

func splitPeriod(input string) (left string, right string) {
	pos := findPeriod(input)
	if pos < 0 {
		return input, ""
	}
	
	return input[:pos], input[pos+1:]
}

// Return whether s has at least 2 periods in it.
func has2Periods(s []rune) (ok bool) {
	seenFirstPeriod := false

	for _, char := range s {
		if char == '.' {
			if seenFirstPeriod {
				return true
			} else {
				seenFirstPeriod = true
			}
		}
	}

	return false
}
