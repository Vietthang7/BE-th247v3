package utils

import "net/mail"

func IsVerifiedEmail(status *bool) bool {
	return status != nil && *status
}
func IsActiveData(active *bool) bool {
	return active != nil && *active
}

// Index trả về chỉ số xuất hiện đầu tiên của v trong s,
// hoặc -1 nếu không có.

func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[S ~[]E, E comparable](s S, v E) bool {
	return Index(s, v) >= 0
}

// CHeck email is valid
func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
