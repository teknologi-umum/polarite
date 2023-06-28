package resources

// TruncateString will safely truncate a string with a given length.
// Yes, the c++ was an intended joke.
func TruncateString(s string, len int) string {
	if len <= 0 {
		return ""
	}

	var t string
	var c int
	for _, v := range s {
		t += string(v)
		c++
		if c >= len {
			break
		}
	}

	return t
}
