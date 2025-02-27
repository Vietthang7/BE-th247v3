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
	ERROR_DATA_MAX_SIZE_250                     = 6096 // Bạn đã nhập quá 250 ký tự.
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
	// lesson data
	ERROR_LESSON_DATA_TEST_NAME_CONTAIN_SPECIAL_CHARACTER = 6069 // Ô dữ liệu không hỗ trợ ký tự đặc biệt.
	ERROR_LESSON_DATA_TEST_NAME_EXIST_IN_CLASS            = 6070 // Tên đề thi đã tồn tại. Vui lòng nhập tên khác
	ERROR_LESSON_DATA_YOUTUBE_LINK_INVALID                = 6060 // Đường link youtube không hợp lệ .
	ERROR_LESSON_DATA_LESSON_DATA_MARSHAL_METADATA_FAILED = 6066 // Lỗi trong quá trình xử lý
	ERROR_LESSON_DATA_PARAGRAPH_CONTENT_INVALID           = 6061 // Nội dung đoạn văn không hợp lệ
	ERROR_LESSON_DATA_EXPIRED_TIME_INVALID                = 6063 // Thời gian kết thúc không hợp lệ
	ERROR_LESSON_DATA_TEST_ID_INVALID                     = 6068 // Đề thi không tìm thấy
	ERROR_LESSON_DATA_URL_DOCUMENT_INVALID                = 6064 // Đường dẫn tài liệu không hợp lệ
	ERROR_LESSON_DATA_TYPE_LESSON_DATA_NOT_FOUND          = 6065 // Kiểu bài giảng không hợp lệ
	ERROR_LESSON_DATA_POINT_TEST_INVALID                  = 6062 // Loại điểm không chính xác
	// schedule class
	ERROR_SCHEDULE_CLASS_LEARNED_CAN_NOT_DELETE = 6047 //Buổi học đã diễn ra. Không thể xóa

)
