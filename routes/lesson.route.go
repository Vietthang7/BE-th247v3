package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func LessonRoute(lesson fiber.Router) {
	// lesson
	lesson.Post("/lesson", controllers.CreateLessons)
	lesson.Patch("/lesson", controllers.UpdateLessons)
	lesson.Delete("/lesson", controllers.DeleteLesson)
	lesson.Get("/lesson", controllers.GetDetailLesson)
	lesson.Get("/lessons", controllers.GetListLessons)
	//lesson-data
	lesson.Post("/lesson-data", controllers.CreateLessonData)
	lesson.Patch("/lesson-data", controllers.UpdateLessonData)
	lesson.Get("/lesson-data", controllers.GetDetailLessonData)
	lesson.Delete("/lesson-data", controllers.DeleteLessonData)
	lesson.Get("/lesson-datas", controllers.GetListLessonDatas)
}
