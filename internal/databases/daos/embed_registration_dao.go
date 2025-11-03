package daos

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/databases/models"
)

type EmbedRegistrationDAO struct {
	db DatabaseExecutor
}

func NewEmbedRegistrationDAO(db DatabaseExecutor) *EmbedRegistrationDAO {
	return &EmbedRegistrationDAO{db: db}
}

// Create inserts a new embed registration
func (dao *EmbedRegistrationDAO) Create(ctx context.Context, registration *models.EmbedRegistration) error {
	sql, args, err := squirrel.
		Insert("embed_registrations").
		Columns("id", "form_id", "allowed_domains", "is_active", "created_at", "updated_at").
		Values(registration.ID, registration.FormID, registration.AllowedDomains, registration.IsActive, registration.CreatedAt, registration.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return jerrors.InternalServerError("failed to build insert query")
	}

	_, err = dao.db.Exec(ctx, sql, args...)
	if err != nil {
		return jerrors.InternalServerError("failed to create embed registration")
	}

	return nil
}

// GetByID retrieves an embed registration by ID
func (dao *EmbedRegistrationDAO) GetByID(ctx context.Context, id uuid.UUID) (*models.EmbedRegistration, error) {
	sql, args, err := squirrel.
		Select("id", "form_id", "allowed_domains", "is_active", "created_at", "updated_at").
		From("embed_registrations").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, jerrors.InternalServerError("failed to build select query")
	}

	row := dao.db.QueryRow(ctx, sql, args...)

	var registration models.EmbedRegistration
	err = row.Scan(
		&registration.ID,
		&registration.FormID,
		&registration.AllowedDomains,
		&registration.IsActive,
		&registration.CreatedAt,
		&registration.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, jerrors.NotFound("embed registration not found")
		}
		return nil, jerrors.InternalServerError("failed to scan embed registration")
	}

	return &registration, nil
}

// GetByFormID retrieves an embed registration by form ID
func (dao *EmbedRegistrationDAO) GetByFormID(ctx context.Context, formID uuid.UUID) (*models.EmbedRegistration, error) {
	sql, args, err := squirrel.
		Select("id", "form_id", "allowed_domains", "is_active", "created_at", "updated_at").
		From("embed_registrations").
		Where(squirrel.Eq{"form_id": formID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, jerrors.InternalServerError("failed to build select query")
	}

	row := dao.db.QueryRow(ctx, sql, args...)

	var registration models.EmbedRegistration
	err = row.Scan(
		&registration.ID,
		&registration.FormID,
		&registration.AllowedDomains,
		&registration.IsActive,
		&registration.CreatedAt,
		&registration.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, jerrors.NotFound("embed registration not found")
		}
		return nil, jerrors.InternalServerError("failed to scan embed registration")
	}

	return &registration, nil
}

// GetActiveByFormID retrieves an active embed registration by form ID
func (dao *EmbedRegistrationDAO) GetActiveByFormID(ctx context.Context, formID uuid.UUID) (*models.EmbedRegistration, error) {
	sql, args, err := squirrel.
		Select("id", "form_id", "allowed_domains", "is_active", "created_at", "updated_at").
		From("embed_registrations").
		Where(squirrel.Eq{"form_id": formID, "is_active": true}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, jerrors.InternalServerError("failed to build select query")
	}

	row := dao.db.QueryRow(ctx, sql, args...)

	var registration models.EmbedRegistration
	err = row.Scan(
		&registration.ID,
		&registration.FormID,
		&registration.AllowedDomains,
		&registration.IsActive,
		&registration.CreatedAt,
		&registration.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, jerrors.NotFound("active embed registration not found")
		}
		return nil, jerrors.InternalServerError("failed to scan embed registration")
	}

	return &registration, nil
}

// Update updates an existing embed registration
func (dao *EmbedRegistrationDAO) Update(ctx context.Context, registration *models.EmbedRegistration) error {
	sql, args, err := squirrel.
		Update("embed_registrations").
		Set("allowed_domains", registration.AllowedDomains).
		Set("is_active", registration.IsActive).
		Set("updated_at", registration.UpdatedAt).
		Where(squirrel.Eq{"id": registration.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return jerrors.InternalServerError("failed to build update query")
	}

	result, err := dao.db.Exec(ctx, sql, args...)
	if err != nil {
		return jerrors.InternalServerError("failed to update embed registration")
	}

	if result.RowsAffected() == 0 {
		return jerrors.NotFound("embed registration not found")
	}

	return nil
}

// Delete removes an embed registration by ID
func (dao *EmbedRegistrationDAO) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := squirrel.
		Delete("embed_registrations").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return jerrors.InternalServerError("failed to build delete query")
	}

	result, err := dao.db.Exec(ctx, sql, args...)
	if err != nil {
		return jerrors.InternalServerError("failed to delete embed registration")
	}

	if result.RowsAffected() == 0 {
		return jerrors.NotFound("embed registration not found")
	}

	return nil
}

// CheckDomainAllowed checks if a domain is allowed for a specific form
func (dao *EmbedRegistrationDAO) CheckDomainAllowed(ctx context.Context, formID uuid.UUID, domain string) (bool, error) {
	sql, args, err := squirrel.
		Select("1").
		From("embed_registrations").
		Where(squirrel.Eq{"form_id": formID, "is_active": true}).
		Where("? = ANY(allowed_domains)", domain).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, jerrors.InternalServerError("failed to build domain check query")
	}

	row := dao.db.QueryRow(ctx, sql, args...)

	var exists int
	err = row.Scan(&exists)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, jerrors.InternalServerError("failed to check domain permission")
	}

	return true, nil
}

// ListByFormIDs retrieves embed registrations for multiple form IDs
func (dao *EmbedRegistrationDAO) ListByFormIDs(ctx context.Context, formIDs []uuid.UUID) ([]*models.EmbedRegistration, error) {
	if len(formIDs) == 0 {
		return []*models.EmbedRegistration{}, nil
	}

	// Convert UUID slice to interface slice for squirrel
	interfaceFormIDs := make([]interface{}, len(formIDs))
	for i, id := range formIDs {
		interfaceFormIDs[i] = id
	}

	sql, args, err := squirrel.
		Select("id", "form_id", "allowed_domains", "is_active", "created_at", "updated_at").
		From("embed_registrations").
		Where(squirrel.Eq{"form_id": interfaceFormIDs}).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, jerrors.InternalServerError("failed to build list query")
	}

	rows, err := dao.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, jerrors.InternalServerError("failed to query embed registrations")
	}
	defer rows.Close()

	var registrations []*models.EmbedRegistration
	for rows.Next() {
		var registration models.EmbedRegistration
		err := rows.Scan(
			&registration.ID,
			&registration.FormID,
			&registration.AllowedDomains,
			&registration.IsActive,
			&registration.CreatedAt,
			&registration.UpdatedAt,
		)
		if err != nil {
			return nil, jerrors.InternalServerError("failed to scan embed registration row")
		}
		registrations = append(registrations, &registration)
	}

	if err = rows.Err(); err != nil {
		return nil, jerrors.InternalServerError("error iterating embed registration rows")
	}

	return registrations, nil
}
