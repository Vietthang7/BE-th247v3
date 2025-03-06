package utils

import (
	"fmt"
	"gorm.io/datatypes"
	"math/rand"
	"net/mail"
	"regexp"
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
func ContainSpecialCharacter(v string) bool {
	pattern := `[!@#~$%^&*()_+|<>?:"\[\]{}\\\/;'’‘]`
	re := regexp.MustCompile(pattern)
	return re.FindString(v) != ""
}
func MixedDateAndTime(startTime *time.Time, gormTime *datatypes.Time) *time.Time {
	if gormTime == nil || startTime == nil {
		return nil
	}
	loc, err := time.LoadLocation("Local") // lấy thông tin múi giờ của hệ thống
	if err != nil {
		return nil
	}
	duration := time.Duration(*gormTime)
	hour := duration / time.Hour
	minutes := (duration % time.Hour) / time.Minute
	seconds := (duration % time.Minute) / time.Second
	nanoseconds := duration % time.Second
	newTime := startTime.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minutes) + time.Second*time.Duration(seconds) + nanoseconds)
	newTime = newTime.In(loc) //Chuyển đổi newTime về múi giờ cục bộ.
	return &newTime
}
func UniqueSliceElements[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}
func ConvertStringToTime(v string) (*time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	time, err := time.Parse(layout, v)
	time = *GetTimeLocation(time)
	if err != nil {
		return nil, err
	}
	return &time, nil
}
func GetTimeLocation(value time.Time) *time.Time {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil
	}
	value = value.In(loc)
	return &value
}
func IsDateInRange(start1, start2, end2 time.Time) bool {
	start1 = start1.Truncate(24 * time.Hour) // chỉ giữ lại phần ngày
	start2 = start2.Truncate(24 * time.Hour) // chỉ giữ lại phần ngày
	end2 = end2.Truncate(24 * time.Hour)     // chỉ giữ lại phần ngày
	return !start1.Before(start2) && !start1.After(end2)
}
