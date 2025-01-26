package consts

const (
	// Lỗi chung
	InvalidReqInput = 1000 // Lỗi dữ liệu đầu vào
	GetFailed       = 1002 // Lấy dữ liệu lỗi
	UpdateFailed    = 1004 // Cập nhật dữ liệu lỗi
	//public data
	ERROR_INTERNAL_SERVER_ERROR = 6104 // Lỗi hệ thống
	ERROR_EXPIRED_TIME          = 6105 // Dữ liệu không phù hợp
	// Auth
	LoginFailed        = 2005 // Đăng nhập thất bại
	UserIsInactive     = 2006 // Tài khoản không hoạt động
	SendOTPFailed      = 2007 // Lỗi gửi OTP
	EmailIsNotVerified = 2008 // Tài khoản chưa xác thực email
)
