package aquapics

var ignoredKeys = []string{"error.html", "index.html", "favicon.ico"}

func IsBadKey(key string) bool {
	for _, v := range ignoredKeys {
		if key == v {
			return true
		}
	}
	return false
}
