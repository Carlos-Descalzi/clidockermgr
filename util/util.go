package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
)

func Max(v1, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func MaxU8(v1, v2 uint8) uint8 {
	if v1 > v2 {
		return v1
	}
	return v2
}

func Min(v1, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

const KB = 1024
const MB = KB * 1024
const GB = MB * 1024
const TB = GB * 1024

func FormatMemory(amount uint64) string {
	if amount < KB {
		return fmt.Sprintf("%d B", amount)
	}
	if amount < MB {
		return fmt.Sprintf("%.2f KB", float32(amount)/KB)
	}
	if amount < GB {
		return fmt.Sprintf("%.2f MB", float32(amount)/MB)
	}
	if amount < TB {
		return fmt.Sprintf("%.2f GB", float32(amount)/GB)
	}
	return fmt.Sprintf("%.2f TB", float32(amount)/TB)
}

func ParseStats(statsBody []byte) *types.Stats {
	var stats types.Stats
	json.Unmarshal(statsBody, &stats)
	return &stats
}

func ParseStatsBody(body io.ReadCloser) *types.Stats {
	var buf = bufio.NewReader(body)
	var result, _ = buf.ReadBytes(byte(0))
	return ParseStats(result)
}
