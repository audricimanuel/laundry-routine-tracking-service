package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/database"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/logging"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
)

type (
	LaundryRepository interface {
		GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId ...string) ([]model.LaundryResponse, error)
	}

	LaundryRepositoryImpl struct {
		db database.DBCollection
	}
)

func NewLaundryRepository(db database.DBCollection) LaundryRepository {
	return &LaundryRepositoryImpl{
		db: db,
	}
}

func (l *LaundryRepositoryImpl) GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId ...string) ([]model.LaundryResponse, error) {
	result := []model.LaundryResponse{}

	log := logging.WithContext(ctx)

	offset := 0
	limit := 10
	if queryParam.Page > 1 {
		offset += (queryParam.Page - 1) * limit
	}

	query := squirrel.Select("id, detail_number, title, laundry_date, total_items, status").
		From("laundries").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	// filter user id
	if len(userId) > 0 {
		query = query.Where(squirrel.Eq{"user_id": userId})
	}

	// filter by laundry date
	if laundryDateFrom := queryParam.LaundryDateFrom; laundryDateFrom != nil {
		query = query.Where(squirrel.Expr(fmt.Sprintf(`DATE(laundry_date) >= %s`, laundryDateFrom.Format(constants.FORMAT_DATE_DEFAULT))))
	}

	if laundryDateTo := queryParam.LaundryDateTo; laundryDateTo != nil {
		query = query.Where(squirrel.Expr(fmt.Sprintf(`DATE(laundry_date) <= %s`, laundryDateTo.Format(constants.FORMAT_DATE_DEFAULT))))
	}

	// filter by detail number
	if detailNumber := queryParam.DetailNumber; detailNumber != "" {
		query = query.Where(squirrel.Eq{"detail_number": detailNumber})
	}

	sql, args := query.PlaceholderFormat(squirrel.Dollar).MustSql()

	rows, err := l.db.PostgresDBSqlx.QueryxContext(ctx, sql, args...)
	if err != nil {
		log.Error("error when getting laundry list:", err)
		return result, errorutils.DefineSQLError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.LaundryResponse
		if err := rows.StructScan(&temp); err != nil {
			log.Error("error when scanning row:", err)
			return result, errorutils.DefineSQLError(err)
		}
		result = append(result, temp)
	}

	return result, nil
}
