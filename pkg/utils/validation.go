package utils

import "regexp"

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func IsValidPhone(phone string) bool {
	regex := `^(\+?[1-9]\d{6,14}|0\d{8,14})$`
	re := regexp.MustCompile(regex)
	return re.MatchString(phone)
}
