package utils

import (
	"fmt"
	"time"
)

func GenerateCouponCodeWithMillis() string {
	now := time.Now()

	// Ví dụ: CP250522141530123
	return fmt.Sprintf("CP%02d%02d%02d%02d%02d%02d%03d",
		now.Year()%100,
		int(now.Month()),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond()/1000000,
	)
}
