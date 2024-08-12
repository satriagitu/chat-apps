package util

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	years := int(duration.Hours() / 8760)
	months := int(duration.Hours() / 730)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes())

	switch {
	case years > 0:
		return fmt.Sprintf("%d tahun yang lalu", years)
	case months > 0:
		return fmt.Sprintf("%d bulan yang lalu", months)
	case days > 0:
		return fmt.Sprintf("%d hari yang lalu", days)
	case hours > 0:
		return fmt.Sprintf("%d jam yang lalu", hours)
	case minutes > 0:
		return fmt.Sprintf("%d menit yang lalu", minutes)
	default:
		return "baru saja"
	}
}
