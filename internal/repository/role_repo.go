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

type RoleRepo struct{ db *database.DB }

func NewRoleRepo(db *database.DB) *RoleRepo { return &RoleRepo{db: db} }

const roleColumns = "id, name, description, created_at, updated_at"

// Create inserts a new role.
func (r *RoleRepo) Create(ctx context.Context, req model.CreateRoleRequest) (*model.Role, error) {
	query, args, err := r.db.Builder.
		Insert("roles").
		Columns("name", "description").
		Values(req.Name, req.Description).
		Suffix("RETURNING " + roleColumns).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var role model.Role
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert role: %w", err)
	}
	return &role, nil
}

// GetByID returns a role by its ID.
func (r *RoleRepo) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	query, args, err := r.db.Builder.
		Select(roleColumns).
		From("roles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var role model.Role
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query role: %w", err)
	}
	return &role, nil
}

// GetByName returns a role by its name.
func (r *RoleRepo) GetByName(ctx context.Context, name string) (*model.Role, error) {
	query, args, err := r.db.Builder.
		Select(roleColumns).
		From("roles").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var role model.Role
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query role: %w", err)
	}
	return &role, nil
}

// List returns all roles.
func (r *RoleRepo) List(ctx context.Context) ([]model.Role, error) {
	query, args, err := r.db.Builder.
		Select(roleColumns).
		From("roles").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// Update updates a role.
func (r *RoleRepo) Update(ctx context.Context, id int64, req model.UpdateRoleRequest) (*model.Role, error) {
	builder := r.db.Builder.Update("roles").Where(sq.Eq{"id": id})

	if req.Name != "" {
		builder = builder.Set("name", req.Name)
	}
	if req.Description != "" {
		builder = builder.Set("description", req.Description)
	}
	builder = builder.Set("updated_at", sq.Expr("NOW()"))

	query, args, err := builder.
		Suffix("RETURNING " + roleColumns).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var role model.Role
	err = r.db.Pool.QueryRow(ctx, query, args...).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	return &role, nil
}

// Delete deletes a role.
func (r *RoleRepo) Delete(ctx context.Context, id int64) error {
	query, args, err := r.db.Builder.
		Delete("roles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	result, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

// GetPermissions returns all permissions for a role.
func (r *RoleRepo) GetPermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	query, args, err := r.db.Builder.
		Select("p.id", "p.name", "p.resource", "p.action", "p.created_at", "p.updated_at").
		From("permissions p").
		Join("role_permissions rp ON rp.permission_id = p.id").
		Where(sq.Eq{"rp.role_id": roleID}).
		OrderBy("p.name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

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

// AssignPermission assigns a permission to a role (idempotent).
func (r *RoleRepo) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	query, args, err := r.db.Builder.
		Insert("role_permissions").
		Columns("role_id", "permission_id").
		Values(roleID, permissionID).
		Suffix("ON CONFLICT (role_id, permission_id) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("assign permission: %w", err)
	}
	return nil
}

// RevokePermission removes a permission from a role.
func (r *RoleRepo) RevokePermission(ctx context.Context, roleID, permissionID int64) error {
	query, args, err := r.db.Builder.
		Delete("role_permissions").
		Where(sq.Eq{"role_id": roleID, "permission_id": permissionID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("revoke permission: %w", err)
	}
	return nil
}
