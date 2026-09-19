package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/ssr0016/template/internal/database"
	"github.com/ssr0016/template/internal/model"
)

type UserRepo struct{ db *database.DB }

func NewUserRepo(db *database.DB) *UserRepo { return &UserRepo{db: db} }

const userColumns = "id, email, name, password_hash, role_id, created_at, updated_at"

// scanUser scans a single user from a row.
func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.RoleID, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query, args, err := r.db.Builder.
		Select(userColumns).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return scanUser(r.db.Pool.QueryRow(ctx, query, args...))
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := r.db.Builder.
		Select(userColumns).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return scanUser(r.db.Pool.QueryRow(ctx, query, args...))
}

func (r *UserRepo) CreateWithPassword(ctx context.Context, email, name, hash string) (*model.User, error) {
	query, args, err := r.db.Builder.
		Insert("users").
		Columns("email", "name", "password_hash").
		Values(email, name, hash).
		Suffix("RETURNING " + userColumns).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return scanUser(r.db.Pool.QueryRow(ctx, query, args...))
}

func (r *UserRepo) List(ctx context.Context, emailFilter string, limit int) ([]model.User, error) {
	// Validate limit to prevent overflow and DoS
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	q := r.db.Builder.
		Select(userColumns).
		From("users").
		OrderBy("id DESC").
		Limit(uint64(limit)) // #nosec G115 - limit is validated to be 1..100

	if emailFilter != "" {
		q = q.Where(sq.ILike{"email": "%" + emailFilter + "%"})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}

// UpdateRole changes a user's role.
func (r *UserRepo) UpdateRole(ctx context.Context, userID, roleID int64) error {
	query, args, err := r.db.Builder.
		Update("users").
		Set("role_id", roleID).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	result, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// GetWithRole returns a user with their role loaded.
func (r *UserRepo) GetWithRole(ctx context.Context, id int64) (*model.User, error) {
	query, args, err := r.db.Builder.
		Select(
			"u.id", "u.email", "u.name", "u.password_hash", "u.role_id",
			"u.created_at", "u.updated_at",
			"r.id", "r.name", "r.description", "r.created_at", "r.updated_at",
		).
		From("users u").
		LeftJoin("roles r ON r.id = u.role_id").
		Where(sq.Eq{"u.id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u model.User
	var rID *int64
	var rName, rDesc *string
	var rCreatedAt, rUpdatedAt *time.Time

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.RoleID,
		&u.CreatedAt, &u.UpdatedAt,
		&rID, &rName, &rDesc, &rCreatedAt, &rUpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}

	// Attach role if loaded
	if rID != nil {
		u.Role = &model.Role{
			ID:          *rID,
			Name:        *rName,
			Description: *rDesc,
			CreatedAt:   *rCreatedAt,
			UpdatedAt:   *rUpdatedAt,
		}
	}

	return &u, nil
}
