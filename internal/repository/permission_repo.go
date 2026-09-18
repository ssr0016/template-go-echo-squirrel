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

type PermissionRepo struct{ db *database.DB }

func NewPermissionRepo(db *database.DB) *PermissionRepo { return &PermissionRepo{db: db} }

const permissionColumns = "id, name, resource, action, created_at, updated_at"

func (r *PermissionRepo) Create(ctx context.Context, name, resource, action string) (*model.Permission, error) {
	query, args, err := r.db.Builder.
		Insert("permissions").
		Columns("name", "resource", "action").
		Values(name, resource, action).
		Suffix("RETURNING " + permissionColumns).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var p model.Permission
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert permission: %w", err)
	}
	return &p, nil
}

func (r *PermissionRepo) GetByID(ctx context.Context, id int64) (*model.Permission, error) {
	query, args, err := r.db.Builder.
		Select(permissionColumns).
		From("permissions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var p model.Permission
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query permission: %w", err)
	}
	return &p, nil
}

func (r *PermissionRepo) GetByName(ctx context.Context, name string) (*model.Permission, error) {
	query, args, err := r.db.Builder.
		Select(permissionColumns).
		From("permissions").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var p model.Permission
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query permission: %w", err)
	}
	return &p, nil
}

func (r *PermissionRepo) List(ctx context.Context) ([]model.Permission, error) {
	query, args, err := r.db.Builder.
		Select(permissionColumns).
		From("permissions").
		OrderBy("resource ASC, action ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	return r.scanPermissions(ctx, query, args...)
}

func (r *PermissionRepo) ListByResource(ctx context.Context, resource string) ([]model.Permission, error) {
	query, args, err := r.db.Builder.
		Select(permissionColumns).
		From("permissions").
		Where(sq.Eq{"resource": resource}).
		OrderBy("action ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	return r.scanPermissions(ctx, query, args...)
}

func (r *PermissionRepo) Delete(ctx context.Context, id int64) error {
	query, args, err := r.db.Builder.
		Delete("permissions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	result, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("permission not found")
	}
	return nil
}

func (r *PermissionRepo) scanPermissions(ctx context.Context, query string, args ...interface{}) ([]model.Permission, error) {
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query permissions: %w", err)
	}
	defer rows.Close()

	var perms []model.Permission
	for rows.Next() {
		var p model.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}
