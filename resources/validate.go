package resources

// Validate whether or not the given body content is larger than 5 MB or not.
// If it's bigger than 5 MB it returns true.
// If it's not it returns false.
func ValidateSize(body []byte) bool {
	return len(body) > 1024*5
}

// This is a O(n) solution and I know this is not a very good algorithm.
// Will fix this later on.
func ValidateID(arr []string, id string) bool {
	for _, v := range arr {
		if v == id {
			return true
		}
	}
	return false
}
