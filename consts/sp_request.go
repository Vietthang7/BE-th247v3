package consts

// CareAssignment
const (
	Manual      = iota + 1 // Phân công thủ công
	EquallyEven            // Chia đều cho nhân viên
	EvenByRate             // Chia cho nhân viên theo tỷ lệ
)

// Support Request
const (
	SRForLeave     = iota + 1 // Xin nghỉ
	SRReserve                 // Bảo lưu
	SRStopStudying            // Dừng học
	SROther                   // Yêu cầu hỗ trợ khác
)

// Product's type in cart
const (
	CartCurriculum = iota + 1
	CartSubject
)
