package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"intern_247/models"
	"log"
	"os"
	"time"
)

const CTimeOut = 10 * time.Second

type DatabaseConfig struct {
	*gorm.DB
	Driver      string `yaml:"driver" env:"DB_DRIVER"`
	Host        string `yaml:"host" env:"DATABASE_HOST"`
	Username    string `yaml:"username" env:"DATABASE_USERNAME"`
	Password    string `yaml:"password" env:"DATABASE_PASSWORD"`
	DBName      string `yaml:"db_name" env:"DATABASE_NAME"`
	Port        string `yaml:"port" env:"DATABASE_PORT"`
	Connections int    `yaml:"connections" env:"DB_CONNECTIONS"`
	Debug       bool   `yaml:"debug"`
	MaxIdleConn int    `env:"MAX_IDLE_CONN"`
	MaxOpenConn int    `env:"MAX_OPEN_CONN"`
}

func (cg *DatabaseConfig) Setup() {
	logrus.SetLevel(logrus.DebugLevel) // tất cả các thông báo từ Debug trở lên sẽ được ghi lại
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Đặt thời gian trễ cho các truy vấn SQL chậm(1 giây)
			LogLevel:                  logger.Silent, // Đặt mức độ log của GORM (tắt log).
			IgnoreRecordNotFoundError: true,          // Bỏ qua các lỗi khi không tìm thấy bản ghi.
			Colorful:                  false,         // Tắt màu sắc trong log.
		},
	)

	mainDbDNS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cg.Username, cg.Password, cg.Host, cg.Port, cg.DBName)
	DB, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:               mainDbDNS,
			DefaultStringSize: 256, // default size for string fields
			// DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			// DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			// DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			// SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		}),
		&gorm.Config{
			PrepareStmt:                              true,
			Logger:                                   newLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		logrus.Panic("Failed to connect database: "+mainDbDNS, err)
	}
	sqlDB, _ := DB.DB()
	if cg.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(cg.MaxIdleConn) // MAX_IDLE_CONN
	}
	if cg.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(cg.MaxOpenConn) // MAX_OPEN_CONN
	}
	if cg.Debug {
		DB = DB.Debug()
	}
	// Lưu đối tượng DB vào trường DB của cấu trúc DatabaseConfig.
	cg.DB = DB

	if err = MigrateDatabase(DB); err != nil {
		logrus.Fatal(err)
	}
	fmt.Println("*************** DB AUTO MIGRATE FINISHED  ***************")
}

func MigrateDatabase(DB *gorm.DB) error {
	if err := DB.AutoMigrate(&models.Center{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.LoginInfo{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.OrganStruct{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Branch{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Certificate{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Student{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Province{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Ward{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.SessionAttendance{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Category{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Class{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.StudyProgress{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.StudentClasses{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Classroom{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.ContactChannel{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Curriculum{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.CustomerSource{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.District{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.EnrollmentPlan{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Lesson{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.LessonData{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.LessonData{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.RoomSchedule{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.ScheduleClass{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Shift{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.StudentCertificates{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.StudentLog{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.ExamResult{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.TestService{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Subject{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.TimeSlot{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.WorkSession{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.SalaryStatement{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.CareAssignment{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.TuitionFeePkgCf{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.OTPLog{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.Permission{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.PermissionTag{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.PermissionGroup{}); err != nil {
		logrus.Debug(err)
	}
	if err := DB.AutoMigrate(&models.SalaryHistory{}); err != nil {
		logrus.Debug(err)
	}

	logrus.Info("Migrating finish")
	return nil
}
