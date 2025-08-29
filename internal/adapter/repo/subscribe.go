package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"subscription-service/internal/app/entity"
	"subscription-service/internal/port"
)

var _ port.SubscriptionRepo = (*Subscription)(nil)

type Subscription struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewSubscription(pool *pgxpool.Pool, logger *zap.Logger) (*Subscription, error) {
	return &Subscription{pool: pool, logger: logger}, nil
}

func (r *Subscription) GetSubscription(ctx context.Context, id string) (*entity.Subscription, error) {
	var sub entity.Subscription

	err := r.pool.QueryRow(ctx, `
    SELECT id, title, price, user_id, start_date, end_date, created_at, updated_at 
    FROM subscriptions 
    WHERE id = $1
`, id).Scan(
		&sub.ID,
		&sub.Title,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, port.ErrNotFound
		}
		return nil, err
	}

	return &sub, nil
}

func (r *Subscription) Create(ctx context.Context, post entity.CreateSubscriptionRequest) error {
	_, err := r.pool.Exec(
		ctx, "INSERT INTO subscriptions"+
			"(id, title, price, user_id, start_date, end_date, created_at, updated_at)"+
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		&post.ID,
		&post.Title,
		&post.Price,
		&post.UserID,
		&post.StartDate,
		&post.EndDate,
		&post.CreatedAt,
		&post.UpdatedAt)
	if err != nil {
		return port.ErrSubscriptionAlreadyExists
	}
	return nil
}

func (r *Subscription) Update(ctx context.Context, post entity.UpdateSubscriptionRequest) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE subscriptions "+
			"SET title = $2, price = $3, start_date = $4, end_date= $5, updated_at = $6 "+
			"WHERE id = $1",
		post.ID,
		post.Title,
		post.Price,
		post.StartDate,
		post.EndDate,
		post.UpdatedAt)
	if err != nil {
		return port.ErrSubscriptionAlreadyExists
	}
	if tag.RowsAffected() == 0 {
		return nil
	}

	return nil
}

func (r *Subscription) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return port.ErrNotFound
	}

	return nil
}

func (r *Subscription) List(ctx context.Context, filter entity.ListSubscriptionFilter) ([]entity.Subscription, error) {
	query := sqlbuilder.Select(
		"id",
		"title",
		"price",
		"user_id",
		"start_date",
		"end_date", "created_at",
		"updated_at",
	).From("subscriptions")

	var and []string

	if filter.Title != nil {
		and = append(and, query.EQ("title", *filter.Title))
	}
	if filter.UserID != nil {
		and = append(and, query.EQ("user_id", *filter.UserID))
	}
	if filter.Price != nil {
		and = append(and, query.EQ("price", *filter.Price))
	}
	if filter.StartDate != nil {
		and = append(and, query.GE("start_date", *filter.StartDate))
	}

	if filter.EndDate != nil {
		and = append(and, query.LE("end_date", *filter.EndDate))
	}

	if filter.Limit != nil {
		query.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		query.Offset(*filter.Offset)
	}

	queryString, args := query.Where(and...).BuildWithFlavor(sqlbuilder.PostgreSQL)
	res, err := r.pool.Query(ctx, queryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer res.Close()

	if !res.Next() {
		return nil, nil
	}

	var subs []entity.Subscription

	for {
		var s entity.Subscription
		if err := res.Scan(
			&s.ID, &s.Title, &s.Price, &s.UserID,
			&s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)

		if !res.Next() {
			break
		}
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *Subscription) Sum(ctx context.Context, filter entity.ListSubscriptionFilter) (int64, error) {
	query := sqlbuilder.Select("SUM(price)").From("subscriptions")

	var and []string

	if filter.Title != nil {
		and = append(and, query.EQ("title", *filter.Title))
	}
	if filter.UserID != nil {
		and = append(and, query.EQ("user_id", *filter.UserID))
	}
	if filter.StartDate != nil {
		and = append(and, query.GE("start_date", *filter.StartDate))
	}
	if filter.EndDate != nil {
		and = append(and, query.LE("end_date", *filter.EndDate))
	}

	queryString, args := query.Where(and...).BuildWithFlavor(sqlbuilder.PostgreSQL)

	var sum int64
	err := r.pool.QueryRow(ctx, queryString, args...).Scan(&sum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return sum, nil
}
