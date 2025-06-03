package repository

import "github.com/audricimanuel/laundry-routine-tracking-service/internal/database"

type (
	EmailRepository interface {
	}

	EmailRepositoryImpl struct {
		db database.DBCollection
	}
)

func NewEmailRepository(db database.DBCollection) EmailRepository {
	return &EmailRepositoryImpl{
		db: db,
	}
}
