package utils

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	sb := strings.Builder{} // để tối ưu hiệu suất khi nối chuỗi.
	sb.Grow(n)              //giúp cấp phát trước bộ nhớ để chứa n ký tự, tránh phân bổ lại nhiều lần.
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))]) // rand.Intn(len(charset)) chọn chỉ số ngẫu nhiên trong charset.
		//	sb.WriteByte(...) ghi từng ký tự vào strings.Builder.
	}
	return sb.String()
}
func GenerateUniqueUsername() string {
	randomStr := RandomString(12)
	id := uuid.New()
	shortUUID := id.String()[:8]
	uniqueUsername := fmt.Sprintf("%s-%s", randomStr, shortUUID)
	return uniqueUsername
}
