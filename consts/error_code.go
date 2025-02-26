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
	ERROR_INTERNAL_SERVER_ERROR                 = 6104 // Lỗi hệ thống
	ERROR_EXPIRED_TIME                          = 6105 // Dữ liệu không phù hợp
	ERROR_UNAUTHORIZED                          = 6101 // Không có quyền truy cập
	ERROR_PERMISSION_DENIED                     = 6099 // Không có quyền truy cập
	ERROR_DATA_LONGER                           = 6100 // Dữ liệu không phù hợp
	ERROR_START_TIME_MUST_SMALLER_THAN_END_TIME = 6097 // Thời gian bắt đầu phải nhỏ hơn thời gian kết thúc

	// Auth
	LoginFailed        = 2005 // Đăng nhập thất bại
	UserIsInactive     = 2006 // Tài khoản không hoạt động
	SendOTPFailed      = 2007 // Lỗi gửi OTP
	EmailIsNotVerified = 2008 // Tài khoản chưa xác thực email

	// Docs category - Danh mục tài liệu
	DocsCategoryExistence  = 2013 // Tên danh mục tài liệu đã tồn tại
	DocsCategoryIsAssigned = 2014 // Danh mục đã có dữ liệu thuộc. Không thể xóa

	//User
	EmailDuplication = 2000
	// Lỗi tồn tại email
	UsernameDuplication     = 2002 // Lỗi tồn tại tên tài khoản
	PhoneDuplication        = 2001 // Lỗi tồn tại số điện thoại
	UpdateUserHasDataLinked = 2003 // Nhân sự đã được gán dữ liệu. Không thể chỉnh sửa.
	UserIsArranged          = 2004 // Nhân sự đã được phân công giảng dạy. Không thể xóa.

	ERROR_STUDENT_HAS_ASSIGNED = "ERROR_STUDENT_HAS_ASSIGNED"

	// category
	ERROR_CATEGORY_HAS_CHILDREN        = 6072 // Danh mục đã có dữ liệu thuộc. Không thể xóa
	ERROR_CATEGORY_HAS_DATA_DEPENDENCY = 6073 // Danh mục đã có dữ liệu thuộc. Không thể xóa
	ERROR_CATEGORY_EXISTS              = 6074 // Danh mục đã tồn tại
	ERROR_CATEGORY_NOT_FOUND           = 6075 // Không tìm thấy danh mục
	ERROR_CATEGORY_HAS_DEPENDENCIES    = 6076 // Category đang bị ràng buộc

	// teacher
	ERROR_TEACHER_NOT_FOUND = 6084 // Không tìm thấy giảng viên
	// subject
	ERROR_SUBJECT_EXISTS                     = 6000 // Môn học đã tồn tại
	ERROR_CAN_NOT_DELETE_SUBJECT_HAS_CLASS   = 6002 // Đã có lớp của môn học. Không thể xóa
	ERROR_CAN_NOT_DELETE_SUBJECT_HAS_STUDENT = 6003 // Môn học đã có học viên đăng ký. Không thể xóa

	// Classroom - Phòng học
	ClassroomExistence  = 2009 // Phòng học đã tồn tại
	ClassroomIsArranged = 2010 // Đã có lớp học được gán. Không thể xóa

	// work session
	ERROR_WORK_SESSION_NAME_EXIST        = 6048 // Tên ca làm đã tồn tại
	ERROR_WORK_SESSION_HAVE_DATA_DEPENDS = 6049 // Ca làm đã có dữ liệu thuộc. Không thể cập nhật
)
