package helpers

import "intern_247/consts"

func IsTeacherOrAsistant(role int64, position int64) bool {
	return (role == consts.CenterHR) && position == consts.Teacher || position == consts.TeachingAssistant
}

func IsStudent(role int64) bool {
	return role == consts.Student
}
