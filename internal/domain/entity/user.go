package entity

//// User đại diện cho một người dùng trong hệ thống
//type User struct {
//	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
//	Name      string    `json:"name" binding:"required"`
//	Email     string    `json:"email" binding:"required,email" gorm:"unique"`
//	Password  string    `json:"-" binding:"required"` // Mật khẩu sẽ không được trả về qua API
//	Role      string    `json:"role"`                 // Vai trò của người dùng, ví dụ: admin, user
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}

//// BeforeCreate gán thời gian tạo cho User
//func (u *User) BeforeCreate() {
//	u.CreatedAt = time.Now()
//}
//
//// BeforeUpdate gán thời gian cập nhật cho User
//func (u *User) BeforeUpdate() {
//	u.UpdatedAt = time.Now()
//}
