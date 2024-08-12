package util_test

import (
	"chat-apps/internal/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeAgo(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		pastTime time.Time
		expected string
	}{
		{"just now", now.Add(-1 * time.Second), "baru saja"},
		{"minutes ago", now.Add(-2 * time.Minute), "2 menit yang lalu"},
		{"hours ago", now.Add(-3 * time.Hour), "3 jam yang lalu"},
		{"days ago", now.Add(-4 * 24 * time.Hour), "4 hari yang lalu"},
		{"months ago", now.Add(-5 * 31 * 24 * time.Hour), "5 bulan yang lalu"},
		{"years ago", now.Add(-6 * 365 * 24 * time.Hour), "6 tahun yang lalu"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.TimeAgo(tt.pastTime)
			assert.Equal(t, tt.expected, got)
		})
	}
}

/*
func TestTimeAgoSpecific(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	fixedTimeSecond := time.Date(2024, 8, 12, 22, 51, 0, 0, location)
	fixedTimeMinute := time.Date(2024, 8, 12, 22, 45, 0, 0, location)
	fixedTimeHour := time.Date(2024, 8, 12, 21, 0, 0, 0, location)
	fixedTimeDay := time.Date(2024, 8, 11, 21, 0, 0, 0, location)
	fixedTimeMonth := time.Date(2024, 7, 12, 15, 0, 0, 0, location)
	fixedTimeYear := time.Date(2023, 8, 12, 15, 0, 0, 0, location)

	now := time.Now()
	tests := []struct {
		name     string
		pastTime time.Time
		expected string
	}{
		{"just now", fixedTimeSecond, "baru saja"},
		{"minutes ago", fixedTimeMinute, "6 menit yang lalu"},
		{"hours ago", fixedTimeHour, "1 jam yang lalu"},
		{"days ago", fixedTimeDay, "1 hari yang lalu"},
		{"months ago", fixedTimeMonth, "1 bulan yang lalu"},
		{"years ago", fixedTimeYear, "1 tahun yang lalu"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.TimeAgo(tt.pastTime)
			assert.Equal(t, tt.expected, got)
		})
	}
}
*/
