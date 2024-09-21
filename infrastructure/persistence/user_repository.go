package persistence

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/infrastructure/security"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
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
	conninfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbname,
		password)

	base, err := NewPgProc(conninfo)
	if err != nil {
		//fmt.Print("cannot connect to database")
		panic("cannot connect to database")
	}

	type para struct {
		ctnId   int
		ctnName string
	}
	// Call an SQL procedures returning a composite value
	type payments struct {
		order_id string
		card_id  string
	}
	type Content struct {
		CntId   int    `pgproc:"cnt_id"`
		CntName string `pgproc:"cnt_name"`
	}
	var item Content
	var par1 int = 1
	var par2 string = "thong"

	err = base.Call(&item, "public", "content_get_test1", par1, par2)
	if err != nil {
		log.Printf("Error calling procedure: %v", err)
		return nil, err
	}

	fmt.Printf("Result: %+v\n", item)

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
