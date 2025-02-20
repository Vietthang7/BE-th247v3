package helpers

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"math"
	"math/big"
	"strconv"
)

const (
	maxDigitsOTP       = 6
	EmailVerifyTitle   = "Xác minh người dùng"
	EmailVerifyContent = "Bạn cần xác minh email? <br> Mã số của bạn là:<br> %s <br> Bạn có thể bỏ qua email này nếu người đăng nhập không phải là bạn."
)

type EmailSchema struct {
	Content   string
	Provider  string // // SMTP host (e.g., "smtp.gmail.com")
	Receivers string // người nhận
	Sender    string // người gửi
	Title     string // tiêu đề
	Password  string // Mật khẩu SMTP
	Port      string // Cổng SMTP (ví dụ: 587 cho Gmail)
}

func SendEmailOTP(emailInfo EmailSchema) (string, error) {
	code := generateOTP(maxDigitsOTP)
	emailInfo.Content = fmt.Sprintf(emailInfo.Content, code)
	//Tạo thư
	m := gomail.NewMessage()
	m.SetHeader("From", emailInfo.Sender)     // Email người gửi
	m.SetHeader("To", emailInfo.Receivers)    //Email người nhận
	m.SetHeader("Subject", emailInfo.Title)   //Tiêu đề
	m.SetBody("text/html", emailInfo.Content) // Nội dung của email
	// Cấu hình SMTP
	port, err := strconv.Atoi(emailInfo.Port)
	if err != nil {
		log.Fatal("Invalid port:", err)
	}
	d := gomail.NewDialer(emailInfo.Provider, port, emailInfo.Sender, emailInfo.Password)
	// gửi email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", err
	}
	log.Print("Email Sent successfully")
	return code, nil
}

func generateOTP(maxDigits uint32) string {
	bi, err := rand.Int(
		rand.Reader, // sử dụng cái này để bảo mật
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))), // tính giá trị tối đa của số ngẫu nhiên và chuyển sang kiểu số lớn trong Go
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}
