package internal

// bool2int converts a bool to an int. true is 1, false is 0.
func Bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Converts an int to a bool. 1 is true, anything else is false.
func Int2bool(a int) bool {
	return a == 1
}
