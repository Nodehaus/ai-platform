package persistence

import (
	"database/sql"
	"time"

	"ai-platform/internal/application/domain/entities"
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
