package utils

import "time"

// format from time to string
func FormatBirthDate(birthDate *time.Time) *string {
	if birthDate == nil {
		return nil
	}
	formattedDate := birthDate.Format("2006-01-02")
	return &formattedDate
}
