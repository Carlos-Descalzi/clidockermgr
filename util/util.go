package util

import "fmt"

func Max(v1 int, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func Min(v1 int, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func FormatMemory(amount int) string {
	if amount < 1024 {
		return fmt.Sprintf("%d B", amount)
	}
	var v = float32(amount) / 1024.0
	if v < 1024 {
		return fmt.Sprintf("%.2f KB", v)
	}
	v /= 1024.0
	if v < 1024 {
		return fmt.Sprintf("%.2f MB", v)
	}
	v /= 1024.0
	if v < 1024 {
		return fmt.Sprintf("%.2f GB", v)
	}
	v /= 1024.0
	return fmt.Sprintf("%.2f TB", v)
}
