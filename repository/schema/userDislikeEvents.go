package schema

type UserDisLikeGenreAndSubGenre struct {
	UserID   string `gorm:"type:varchar(255);default:uuid_generate_v4()"`
	Genre    string `gorm:"type:varchar(255);"`
	SubGenre string `gorm:"type:varchar(255);"`
}
