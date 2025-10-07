package persistence

import (
	"database/sql"
	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
	"time"
)

type DeploymentRepositoryImpl struct {
	Db *sql.DB
}

func (r *DeploymentRepositoryImpl) Create(deployment *entities.Deployment) error {
	query := `INSERT INTO deployments (id, model_name, api_key, project_id, finetune_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	deployment.CreatedAt = now
	deployment.UpdatedAt = now

	_, err := r.Db.Exec(query,
		deployment.ID,
		deployment.ModelName,
		deployment.APIKey,
		deployment.ProjectID,
		deployment.FinetuneID,
		deployment.CreatedAt,
		deployment.UpdatedAt,
	)

	return err
}

func (r *DeploymentRepositoryImpl) GetByID(id uuid.UUID) (*entities.Deployment, error) {
	query := `SELECT id, model_name, api_key, project_id, finetune_id, created_at, updated_at
			  FROM deployments WHERE id = $1`

	var model DeploymentRepositoryModel
	err := r.Db.QueryRow(query, id).Scan(
		&model.ID,
		&model.ModelName,
		&model.APIKey,
		&model.ProjectID,
		&model.FinetuneID,
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

func (r *DeploymentRepositoryImpl) GetByProjectID(projectID uuid.UUID) ([]entities.Deployment, error) {
	query := `SELECT id, model_name, api_key, project_id, finetune_id, created_at, updated_at
			  FROM deployments WHERE project_id = $1 ORDER BY created_at DESC`

	rows, err := r.Db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []entities.Deployment
	for rows.Next() {
		var model DeploymentRepositoryModel
		err := rows.Scan(
			&model.ID,
			&model.ModelName,
			&model.APIKey,
			&model.ProjectID,
			&model.FinetuneID,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, *model.ToEntity())
	}

	return deployments, nil
}

func (r *DeploymentRepositoryImpl) GetByFinetuneID(finetuneID uuid.UUID) (*entities.Deployment, error) {
	query := `SELECT id, model_name, api_key, project_id, finetune_id, created_at, updated_at
			  FROM deployments WHERE finetune_id = $1`

	var model DeploymentRepositoryModel
	err := r.Db.QueryRow(query, finetuneID).Scan(
		&model.ID,
		&model.ModelName,
		&model.APIKey,
		&model.ProjectID,
		&model.FinetuneID,
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

func (r *DeploymentRepositoryImpl) GetByProjectIDAndModelName(projectID uuid.UUID, modelName string) (*entities.Deployment, error) {
	query := `SELECT id, model_name, api_key, project_id, finetune_id, created_at, updated_at
			  FROM deployments WHERE project_id = $1 AND model_name = $2`

	var model DeploymentRepositoryModel
	err := r.Db.QueryRow(query, projectID, modelName).Scan(
		&model.ID,
		&model.ModelName,
		&model.APIKey,
		&model.ProjectID,
		&model.FinetuneID,
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

func (r *DeploymentRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM deployments WHERE id = $1`
	_, err := r.Db.Exec(query, id)
	return err
}
