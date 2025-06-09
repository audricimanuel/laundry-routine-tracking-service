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
	"strings"
)

type (
	LaundryRepository interface {
		GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId ...string) ([]model.LaundryResponse, error)
		GetCategoryList(ctx context.Context, userId ...string) ([]model.CategoryResponse, error)
		GetCategoryById(ctx context.Context, userId string, ids ...string) ([]model.CategoryResponse, error)
		AddCategory(ctx context.Context, name string, userId ...string) error
		IsExistedCategoryName(ctx context.Context, name, userId string) bool
		AddLaundryData(ctx context.Context, userId string, request model.AddLaundryRequest) error
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

	query := squirrel.Select("id, title, laundry_date, total_items, status").
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

func (l *LaundryRepositoryImpl) GetCategoryList(ctx context.Context, userId ...string) ([]model.CategoryResponse, error) {
	var result []model.CategoryResponse
	log := logging.WithContext(ctx)

	baseQuery := squirrel.Select(`id, "name", user_id, is_active`).
		From("categories").
		OrderBy(`"name"`)

	if len(userId) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"user_id": userId})
	}

	query, args := baseQuery.PlaceholderFormat(squirrel.Dollar).MustSql()

	rows, err := l.db.PostgresDBSqlx.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error("error when getting category list:", err)
		return result, errorutils.DefineSQLError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.CategoryResponse
		if err := rows.StructScan(&temp); err != nil {
			log.Error("error when scanning row:", err)
			return result, errorutils.DefineSQLError(err)
		}
		result = append(result, temp)
	}

	return result, nil
}

func (l *LaundryRepositoryImpl) GetCategoryById(ctx context.Context, userId string, ids ...string) ([]model.CategoryResponse, error) {
	log := logging.WithContext(ctx)

	baseQuery := squirrel.Select(`id, "name", user_id, is_active`).
		From("categories").
		Where(squirrel.Eq{"id": ids}).
		OrderBy(`name`)

	if len(userId) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"user_id": userId})
	}

	query, args := baseQuery.PlaceholderFormat(squirrel.Dollar).MustSql()

	var result []model.CategoryResponse
	rows, err := l.db.PostgresDBSqlx.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error("error when getting category:", err)
		return result, errorutils.DefineSQLError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.CategoryResponse
		if err := rows.StructScan(&temp); err != nil {
			log.Error("error when scanning row:", err)
			return result, errorutils.DefineSQLError(err)
		}
		result = append(result, temp)
	}

	return result, nil
}

func (l *LaundryRepositoryImpl) AddCategory(ctx context.Context, name string, userId ...string) error {
	query, args := squirrel.Insert("categories").
		Columns("name", "user_id").
		Values(name, userId).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	_, err := l.db.PostgresDBSqlx.ExecContext(ctx, query, args...)
	if err != nil {
		logging.WithContext(ctx).Error("error when adding category:", err)
		return errorutils.DefineSQLError(err)
	}

	return nil
}

func (l *LaundryRepositoryImpl) IsExistedCategoryName(ctx context.Context, name, userId string) bool {
	query, args := squirrel.Select("id").
		From("categories").
		Where(squirrel.And{
			squirrel.Eq{"user_id": userId},
			squirrel.Expr(fmt.Sprintf(`LOWER("name") = %s`, strings.ToLower(name))),
		}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	var id string
	err := l.db.PostgresDBSqlx.QueryRowxContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return false
	}

	return id != ""
}

func (l *LaundryRepositoryImpl) AddLaundryData(ctx context.Context, userId string, request model.AddLaundryRequest) error {
	// TODO: add main laundry data

	// TODO: add laundry items

	return nil
}
