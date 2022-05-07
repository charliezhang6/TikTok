package Model

type User struct {
	ID       uint64 `gorm:"column:user_id"`
	Name     string `gorm:"column:user_name"`
	Password string `gorm:"column:user_password"`
	Salt     string `gorm:":column:salt"`
	Token    string `gorm:"column:token"`
}
