package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ssr0016/template/internal/database"
	"github.com/ssr0016/template/internal/model"
)

type UserRepo struct{ db *database.DB }

func NewUserRepo(db *database.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query, args, err := r.db.Builder.
		Select("id", "email", "name", "password_hash", "created_at", "updated_at").
		From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var u model.User
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := r.db.Builder.
		Select("id", "email", "name", "password_hash", "created_at", "updated_at").
		From("users").Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var u model.User
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) CreateWithPassword(ctx context.Context, email, name, hash string) (*model.User, error) {
	query, args, err := r.db.Builder.
		Insert("users").Columns("email", "name", "password_hash").
		Values(email, name, hash).
		Suffix("RETURNING id, email, name, password_hash, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var u model.User
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, emailFilter string, limit int) ([]model.User, error) {
	q := r.db.Builder.
		Select("id", "email", "name", "password_hash", "created_at", "updated_at").
		From("users").OrderBy("id DESC").Limit(uint64(limit))
	if emailFilter != "" {
		q = q.Where(sq.ILike{"email": "%" + emailFilter + "%"})
	}
	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
