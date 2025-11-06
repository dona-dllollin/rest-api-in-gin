package database

import "database/sql"

type Models struct {
	Users    UserModel
	Events   EventModel
	Attendes AttendeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Events:   EventModel{DB: db},
		Attendes: AttendeModel{DB: db},
	}
}
