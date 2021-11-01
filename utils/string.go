package utils

func StrCoalesce(values ...string) string {
	for _, el := range values {
		if len(el) > 0 {
			return el
		}
	}

	return ""
}
