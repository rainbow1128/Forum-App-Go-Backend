package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) []string {

	var errorMessages []string
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			err = errors.New("Required Nickname")
			errorMessages = append(errorMessages, err.Error())
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages = append(errorMessages, err.Error())
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages = append(errorMessages, err.Error())
		}
		if err = checkmail.ValidateFormat(u.Email); err != nil {
			err = errors.New("Invalid Email")
			errorMessages = append(errorMessages, err.Error())
		}
	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages = append(errorMessages, err.Error())
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages = append(errorMessages, err.Error())
		}
		if err = checkmail.ValidateFormat(u.Email); err != nil {
			err = errors.New("Invalid Email")
			errorMessages = append(errorMessages, err.Error())
		}
	default:
		if u.Nickname == "" {
			err = errors.New("Required Nickname")
			errorMessages = append(errorMessages, err.Error())
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages = append(errorMessages, err.Error())
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages = append(errorMessages, err.Error())
		}
		if err = checkmail.ValidateFormat(u.Email); err != nil {
			err = errors.New("Invalid Email")
			errorMessages = append(errorMessages, err.Error())
		}
	}
	return errorMessages
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}