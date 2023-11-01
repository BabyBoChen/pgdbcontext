package utils

func HasKey(dict map[string]interface{}, key string) bool {
	hasKey := false
	for k := range dict {
		if k == key {
			hasKey = true
			break
		}
	}
	return hasKey
}
