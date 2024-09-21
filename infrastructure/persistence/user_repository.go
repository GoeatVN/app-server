package persistence

import (
	"context"
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/infrastructure/security"
	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
	"time"
)

type UserRepo struct {
	// sử dụng baseRepo
	//persistence.Repository[entity.User]
	db *gorm.DB
}

// NewUserRepository creates and returns a new instance of UserRepo.
// It initializes the UserRepo with a database connection.
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db}
}

// UserRepo implements the repository.UserRepository interface
var _ repository.UserRepository = &UserRepo{}

func (r *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error

	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	//err := r.db.Debug().Find(&users).Error

	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	connInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbname,
		password)
	// Example of a struct result (each procedure can have different results)
	type Result struct {
		CntId          int       `pgColumn:"cnt_id"`
		CntName        string    `pgColumn:"cnt_name"`
		CntInt         int       `pgColumn:"cnt_int"`
		CntTest        string    `pgColumn:"cnt_test"`
		CntDate        time.Time `pgColumn:"cnt_date"`
		CntTimeStamp   time.Time `pgColumn:"cnt_timestamp"`
		CntTime        time.Time `pgColumn:"cnt_time"`
		CtnNumeric     float64   `pgColumn:"cnt_numeric"`
		CtnBoolean     bool      `pgColumn:"cnt_boolean"`
		CtnTimestampTz time.Time `pgColumn:"cnt_timestamptz"`
	}

	// Struct containing parameters with `pgParam` tags
	type ContentParams struct {
		CntId   int    `pgParam:"prm_id"`
		CntName string `pgParam:"prm_name"`
		CurRef  string `pgCur:"cur_ref"`
	}
	// Kết nối đến database qua biến môi trường DATABASE_URL
	conn, err := pgx.Connect(context.Background(), connInfo)
	if err != nil {
		log.Fatalf("Không thể kết nối đến database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Định nghĩa struct chứa tham số đầu vào cho stored procedure
	params := ContentParams{
		CntId:   1,
		CntName: "thong123",
	}

	// Khởi tạo slice để chứa kết quả trả về
	var results []Result

	// Gọi hàm ExecuteStoredProcedure từ module db
	err = ExecuteStoredProcedureWithCursor(conn, "content_get_test", params, &results)
	if err != nil {
		log.Fatalf("Lỗi khi gọi stored procedure: %v\n", err)
	}

	// In kết quả
	for _, result := range results {
		fmt.Printf("CntId: %d, CntName: %s, CntInt: %d, CntDate: %s\n", result.CntId, result.CntName, result.CntInt, result.CntDate)
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := map[string]string{}
	err := r.db.Debug().Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	//Verify the password
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}
	return &user, nil
}
