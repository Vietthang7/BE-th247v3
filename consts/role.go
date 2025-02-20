package consts

// Role
const (
	Root = iota + 1
	CenterOwner
	CenterAdmin
	CenterHR // Nhân sự
	Student  // Học viên
)

// Vị trí
const (
	Teacher           = iota + 1 // Giảng viên
	TeachingAssistant            // Trợ giảng
	NonePosition
	CareAssignee // Phân công chăm sóc
)
const (
	USER_ROOT = "root"
)
const PermissionsGoingWith = `{
            "create": [
                "list",
                "read",
                "update",
                "delete",
                "print"
            ],
            "update": [
                "list",
                "read",
                "print",
                "delete"
            ],
            "delete": [
                "list",
                "read"
            ],
            "read": [
                "list"
            ],
            "import": [
                "create",
                "list",
                "delete",
                "update"
            ]
        }`
