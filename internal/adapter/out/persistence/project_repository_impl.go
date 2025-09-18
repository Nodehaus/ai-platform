package persistence

import (
	"database/sql"
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
	"time"
)

type ProjectRepositoryImpl struct {
	Db *sql.DB
}


func (r *ProjectRepositoryImpl) Create(project *entities.Project) error {
	query := `INSERT INTO projects (id, name, owner_id, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := r.Db.Exec(query,
		project.ID,
		project.Name,
		project.OwnerID,
		string(project.Status),
		project.CreatedAt,
		project.UpdatedAt,
	)

	return err
}

func (r *ProjectRepositoryImpl) GetByID(id uuid.UUID) (*entities.Project, error) {
	query := `SELECT id, name, owner_id, status, created_at, updated_at FROM projects WHERE id = $1`

	var model ProjectRepositoryModel
	err := r.Db.QueryRow(query, id).Scan(
		&model.ID,
		&model.Name,
		&model.OwnerID,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *ProjectRepositoryImpl) GetByOwnerID(ownerID uuid.UUID) ([]entities.Project, error) {
	query := `SELECT id, name, owner_id, status, created_at, updated_at FROM projects WHERE owner_id = $1 ORDER BY created_at DESC`

	rows, err := r.Db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var model ProjectRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.OwnerID,
			&model.Status,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *model.ToEntity())
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) GetActiveByOwnerID(ownerID uuid.UUID) ([]entities.Project, error) {
	query := `SELECT id, name, owner_id, status, created_at, updated_at FROM projects WHERE owner_id = $1 AND status = 'ACTIVE' ORDER BY created_at DESC`

	rows, err := r.Db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []entities.Project
	for rows.Next() {
		var model ProjectRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.OwnerID,
			&model.Status,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *model.ToEntity())
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) ExistsByNameAndOwnerID(name string, ownerID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE name = $1 AND owner_id = $2 AND status != 'DELETED')`

	var exists bool
	err := r.Db.QueryRow(query, name, ownerID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ProjectRepositoryImpl) Update(project *entities.Project) error {
	query := `UPDATE projects SET name = $2, status = $3, updated_at = $4 WHERE id = $1`

	project.UpdatedAt = time.Now()

	_, err := r.Db.Exec(query,
		project.ID,
		project.Name,
		string(project.Status),
		project.UpdatedAt,
	)

	return err
}

func (r *ProjectRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.Db.Exec(query, id)
	return err
}