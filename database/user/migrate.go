package user

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Claim{},
		&Review{},
		&BusinessFollow{},
		&ProductLike{},
		&ProductSave{},
		&SpotlightLike{},
		&SpotlightSave{},
	)
}
