package consts

const (
	// Lỗi chung
	CreateFailed    = 1001 // Tạo dữ liệu lỗi
	InvalidReqInput = 1000 // Lỗi dữ liệu đầu vào
	GetFailed       = 1002 // Lấy dữ liệu lỗi
	UpdateFailed    = 1004 // Cập nhật dữ liệu lỗi
	RegisterFailed  = 1007 // Đăng ký thất bại
	DataNotFound    = 1005 // Không tìm thấy dữ liệu
	Forbidden       = 1003 // Không có quyền thực hiện thao tác này
	DeletedFailed   = 1006 // Xóa dữ liệu lỗi

	//public data
	ERROR_INTERNAL_SERVER_ERROR = 6104 // Lỗi hệ thống
	ERROR_EXPIRED_TIME          = 6105 // Dữ liệu không phù hợp
	ERROR_UNAUTHORIZED          = 6101 // Không có quyền truy cập
	ERROR_PERMISSION_DENIED     = 6099 // Không có quyền truy cập

	// Auth
	LoginFailed        = 2005 // Đăng nhập thất bại
	UserIsInactive     = 2006 // Tài khoản không hoạt động
	SendOTPFailed      = 2007 // Lỗi gửi OTP
	EmailIsNotVerified = 2008 // Tài khoản chưa xác thực email

	// Docs category - Danh mục tài liệu
	DocsCategoryExistence  = 2013 // Tên danh mục tài liệu đã tồn tại
	DocsCategoryIsAssigned = 2014 // Danh mục đã có dữ liệu thuộc. Không thể xóa

	//User
	EmailDuplication        = 2000 // Lỗi tồn tại email
	UsernameDuplication     = 2002 // Lỗi tồn tại tên tài khoản
	PhoneDuplication        = 2001 // Lỗi tồn tại số điện thoại
	UpdateUserHasDataLinked = 2003 // Nhân sự đã được gán dữ liệu. Không thể chỉnh sửa.
)
