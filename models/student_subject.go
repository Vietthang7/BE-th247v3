package models

import "github.com/google/uuid"

type StudentSubject struct {
	StudentId    uuid.UUID `json:"student_id" gorm:"index:idx_student_subjects_student_id"`
	StudyNeedsId uuid.UUID `json:"study_needs_id" gorm:"index:idx_student_subjects_study_needs_id"`
	SubjectId    uuid.UUID `json:"subject_id" gorm:"index:idx_student_subjects_subject_id"`
}

type InputReserved struct {
	StudentIds uuid.UUIDs `json:"student_ids"`
	ClassId    uuid.UUID  `json:"class_id"`
}
