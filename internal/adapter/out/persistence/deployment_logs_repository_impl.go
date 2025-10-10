package persistence

import (
	"database/sql"
	"time"

	"ai-platform/internal/application/domain/entities"

	"github.com/google/uuid"
)

type DeploymentLogsRepositoryImpl struct {
	Db *sql.DB
}

func (r *DeploymentLogsRepositoryImpl) Create(log *entities.DeploymentLogs) error {
	query := `INSERT INTO deployment_logs (id, deployment_id, tokens_in, tokens_out, input, output, delay_time, execution_time, source, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	now := time.Now()
	log.CreatedAt = now
	log.UpdatedAt = now

	_, err := r.Db.Exec(query,
		log.ID,
		log.DeploymentID,
		log.TokensIn,
		log.TokensOut,
		log.Input,
		log.Output,
		log.DelayTime,
		log.ExecutionTime,
		log.Source,
		log.CreatedAt,
		log.UpdatedAt,
	)

	return err
}

func (r *DeploymentLogsRepositoryImpl) GetLatest(deploymentID uuid.UUID, limit int) ([]*entities.DeploymentLogs, error) {
	query := `SELECT id, deployment_id, tokens_in, tokens_out, input, output, delay_time, execution_time, source, created_at, updated_at
			  FROM deployment_logs
			  WHERE deployment_id = $1
			  ORDER BY created_at DESC
			  LIMIT $2`

	rows, err := r.Db.Query(query, deploymentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entities.DeploymentLogs
	for rows.Next() {
		log := &entities.DeploymentLogs{}
		err := rows.Scan(
			&log.ID,
			&log.DeploymentID,
			&log.TokensIn,
			&log.TokensOut,
			&log.Input,
			&log.Output,
			&log.DelayTime,
			&log.ExecutionTime,
			&log.Source,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *DeploymentLogsRepositoryImpl) GetAll(deploymentID uuid.UUID) ([]*entities.DeploymentLogs, error) {
	query := `SELECT id, deployment_id, tokens_in, tokens_out, input, output, delay_time, execution_time, source, created_at, updated_at
			  FROM deployment_logs
			  WHERE deployment_id = $1
			  ORDER BY created_at DESC`

	rows, err := r.Db.Query(query, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entities.DeploymentLogs
	for rows.Next() {
		log := &entities.DeploymentLogs{}
		err := rows.Scan(
			&log.ID,
			&log.DeploymentID,
			&log.TokensIn,
			&log.TokensOut,
			&log.Input,
			&log.Output,
			&log.DelayTime,
			&log.ExecutionTime,
			&log.Source,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
