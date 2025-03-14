package models

import (
	"api/pkg/config"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin" gorm:"default:false"`
	IsEmployee bool   `json:"is_employee" gorm:"default:false"` // Add this line
}

func init() {
	db = config.GetDB()
	db.AutoMigrate(&User{})

	// Check if admin user exists, if not create one
	var adminCount int64
	db.Model(&User{}).Where("is_admin = ?", true).Count(&adminCount)
	if adminCount == 0 {
		// Create default admin user
		CreateAdminUser()
	}

	// Check if employee user exists, if not create one
	var employeeCount int64
	db.Model(&User{}).Where("is_employee = ?", true).Count(&employeeCount)
	if employeeCount == 0 {
		// Create default employee user
		_, err := CreateEmployeeUser("employee@gmail.com", "Testing123")
		if err != nil {
			panic("Failed to create employee user: " + err.Error())
		}
	}
}

func CreateAdminUser() {
	// Default admin credentials - you should change these!
	adminEmail := "admin@urcreativesolutions.com"
	adminPassword := "admin123"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash admin password")
	}

	// Create admin user
	adminUser := User{
		Email:    adminEmail,
		Password: string(hashedPassword),
		IsAdmin:  true,
	}

	result := db.Create(&adminUser)
	if result.Error != nil {
		panic("Failed to create admin user: " + result.Error.Error())
	}
}

func CreateEmployeeUser(email, password string) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create employee user
	employeeUser := User{
		Email:      email,
		Password:   string(hashedPassword),
		IsAdmin:    false,
		IsEmployee: true,
	}

	result := db.Create(&employeeUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &employeeUser, nil
}

func (u *User) HashPassword(plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

func CreateUser(user *User) (*User, error) {
	// Check if user already exists
	var existingUser User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	// Create the user
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
