package repository

import (
	"Recommendation-System/repository/schema"
	"bytes"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	strength = 10
	Salt     = "salt"
)

type UserRepo interface {
	InsertUser(user schema.User) error
	SelectUserByID(id uuid.UUID) (schema.User, error)
	SelectUserByUsername(username string) (schema.User, error)
	InsertUserLikeGenreAndSubGenre(userID string, genre, subgenre string) error
	InsertUserDisLikeGenreAndSubGenre(userID string, genre, subgenre string) error
	SelectUserLikedGenresAndSubGenre(userID string) ([]schema.UserLikeGenreAndSubGenre, error)
	SelectUserDisLikedGenreAndSubGenre(userID string) ([]schema.UserDisLikeGenreAndSubGenre, error)
}

type userRepo struct {
	db *GormDatabase
}

func NewUserRepo(db *GormDatabase) UserRepo {
	return &userRepo{db}
}

func (u *userRepo) SelectUserByID(id uuid.UUID) (schema.User, error) {
	var user schema.User
	if err := u.db.DB.Model(&schema.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) InsertUser(user schema.User) error {
	// add salt to password and hash it
	passwordBuf := bytes.Buffer{}
	passwordBuf.WriteString(user.Password)
	passwordBuf.WriteString(Salt)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBuf.Bytes(), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if err := u.db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepo) SelectUserByUsername(username string) (schema.User, error) {
	var user schema.User
	if err := u.db.DB.Model(&schema.User{}).Where("user_name = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) InsertUserLikeGenreAndSubGenre(userID string, genre, subgenre string) error {
	userLikeGenre := schema.UserLikeGenreAndSubGenre{
		UserID:   userID,
		Genre:    genre,
		SubGenre: subgenre,
	}
	if err := u.db.DB.Create(&userLikeGenre).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepo) InsertUserDisLikeGenreAndSubGenre(userID string, genre, subgenre string) error {
	userDisLikeGenre := schema.UserDisLikeGenreAndSubGenre{
		UserID:   userID,
		Genre:    genre,
		SubGenre: subgenre,
	}
	if err := u.db.DB.Create(&userDisLikeGenre).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepo) SelectUserLikedGenresAndSubGenre(userID string) ([]schema.UserLikeGenreAndSubGenre, error) {
	var userLikedGenres []schema.UserLikeGenreAndSubGenre
	if err := u.db.DB.Model(&schema.UserLikeGenreAndSubGenre{}).Where("user_id = ?", userID).Find(&userLikedGenres).Error; err != nil {
		return userLikedGenres, err
	}
	return userLikedGenres, nil
}

func (u *userRepo) SelectUserDisLikedGenreAndSubGenre(userID string) ([]schema.UserDisLikeGenreAndSubGenre, error) {
	var userDisLikedGenres []schema.UserDisLikeGenreAndSubGenre
	if err := u.db.DB.Model(&schema.UserDisLikeGenreAndSubGenre{}).Where("user_id = ?", userID).Find(&userDisLikedGenres).Error; err != nil {
		return userDisLikedGenres, err
	}
	return userDisLikedGenres, nil
}
