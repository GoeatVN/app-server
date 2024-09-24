package persistence

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	persistence2 "food-app/infrastructure/database"
	"food-app/infrastructure/security"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo implements the repository.UserRepository interface
// It will hold the connection to the database
// and implement the methods in the UserRepository interface
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
	// Example of a struct result (each procedure can have different results)
	type Result struct {
		CntName        string    `pgColumn:"cnt_name" json:"cntName"`
		CntInt         int       `pgColumn:"cnt_int" json:"cntInt"`
		CntTest        string    `pgColumn:"cnt_test" json:"cntTest"`
		CntDate        time.Time `pgColumn:"cnt_date" json:"cntDate"`
		CntTimeStamp   time.Time `pgColumn:"cnt_timestamp"`
		CntTime        time.Time `pgColumn:"cnt_time"`
		CtnNumeric     float64   `pgColumn:"cnt_numeric"`
		CtnBoolean     bool      `pgColumn:"cnt_boolean"`
		CtnTimestampTz time.Time `pgColumn:"cnt_timestamptz"`
		ctnTemp        float64   `pgColumn:"cnt_temp"`
		CntId          int       `pgColumn:"cnt_id" json:"cntId"`
		/*
			CntId        int       `pgColumn:"cnt_id" json:"cntId"`
			CntTimeStamp time.Time `pgColumn:"cnt_timestamp"`
		*/
	}

	// Struct containing parameters with `pgParam` tags
	type ContentParams struct {
		CntId   int       `pgParam:"prm_id"`
		CntName string    `pgParam:"prm_name"`
		CntDate time.Time `pgParam:"prm_date"`
		CurRef  string    `pgCur:"cur_ref"`
	}
	// Định nghĩa struct chứa tham số đầu vào cho stored procedure
	params := ContentParams{
		CntId:   1,
		CntName: "thong",
		CntDate: time.Now(),
	}

	// Khởi tạo slice để chứa kết quả trả về
	var results []Result
	//// Tạo một context cơ bản
	//ctx := context.Background()
	//
	//// Nếu bạn muốn đặt timeout
	//ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()
	// Gọi hàm ExecuteStoredProcedure từ module db
	results, err := persistence2.ExecuteStoredProcedureWithCursor[Result]("content_get_test", params)
	if err != nil {
		fmt.Println("Lỗi khi gọi stored procedure: %v\n", err)
	}

	// In kết quả
	for _, result := range results {
		fmt.Printf("CntId: %d, CntTimeStamp: %s\n", result.CntId, result.CntTimeStamp)
		//fmt.Printf("CntId: %d, CntName: %s, CntInt: %d, CntDate: %s\n", result.CntId, result.CntName, result.CntInt, result.CntDate)
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	users = append(users, entity.User{ID: 1, FirstName: "thong", LastName: "le", Email: "thongle@gmail.com"})
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
