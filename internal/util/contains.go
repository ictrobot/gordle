package util

func SliceContains[t comparable](slice []t, x t) bool {
	for _, v := range slice {
		if v == x {
			return true
		}
	}
	return false
}
