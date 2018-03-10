package main

func DefaultString(val, def string) string {
	if val == "" {
		return def
	}

	return val
}

func StripEmptyStrings(slice []string) []string {
	out := []string{}
	if slice == nil {
		return out
	}

	for _, s := range slice {
		if s != "" {
			out = append(out, s)
		}
	}

	return out
}
