package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

func GenerateTrackingNumber() string {
	// Lấy ngày hiện tại
	now := time.Now()
	// Format: YYYYMMDDHHMMSS (năm-tháng-ngày-giờ-phút-giây)
	dateTimeStr := now.Format("20060102150405")

	// Tạo chuỗi random 6 ký tự (chữ và số)
	randomStr := generateRandomString(6)

	// Kết hợp thành tracking number
	trackingNumber := fmt.Sprintf("TRK-%s-%s", dateTimeStr, randomStr)

	return trackingNumber
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	result.Grow(length)

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result.WriteByte(charset[num.Int64()])
	}

	return result.String()
}
