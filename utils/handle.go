package utils

import (
	"fmt"
	"math/rand"
	"net/mail"
	"strconv"
	"time"
	"unicode/utf8"
)

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
func IsValidStrLen(v string, max int) bool {
	if v == "" {
		return false
	}
	return utf8.RuneCountInString(v) <= max
}

func getLastTwoDigitsOfCurrentYear() int {
	currentYear := time.Now().Year()
	return currentYear % 100
}
func generateRandomFiveDigitNumber() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	fmt.Println(r)
	return r.Intn(90000) + 10000
}

func GenerateRandomCodeFormatByKey(key string) string {
	return key + strconv.Itoa(getLastTwoDigitsOfCurrentYear()) + strconv.Itoa(generateRandomFiveDigitNumber())
	//strconv.Itoa(...) chuyển số nguyên thành chuỗi để có thể ghép nối.
}
