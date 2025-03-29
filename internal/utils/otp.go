package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())  // Đảm bảo randomness mỗi lần gọi
	otp := rand.Intn(900000) + 100000 // Sinh số từ 100000 đến 999999

	return fmt.Sprintf("%06d", otp) // Đảm bảo luôn có 6 chữ số
}
