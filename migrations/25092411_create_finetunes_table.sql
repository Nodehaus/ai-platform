-- Create finetunes table
CREATE TABLE finetunes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    base_model_name VARCHAR(255) NOT NULL,
    model_size_gb INTEGER,
    model_size_parameter INTEGER,
    model_dtype VARCHAR(50),
    model_quantization VARCHAR(50),
    inference_samples_json TEXT NOT NULL DEFAULT '[]',
    training_dataset_id UUID NOT NULL REFERENCES training_datasets(id) ON DELETE RESTRICT,
    training_dataset_number_examples INTEGER,
    training_dataset_select_random BOOLEAN NOT NULL DEFAULT false,
    training_time_seconds DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'PLANNING' CHECK (status IN ('PLANNING', 'RUNNING', 'ABORTED', 'FAILED', 'DONE', 'DELETED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Composite unique constraint for project versioning
    UNIQUE(project_id, version)
);

-- Create indexes for performance
CREATE INDEX idx_finetunes_project_id ON finetunes(project_id);
CREATE INDEX idx_finetunes_status ON finetunes(status);
CREATE INDEX idx_finetunes_version ON finetunes(project_id, version);
CREATE INDEX idx_finetunes_training_dataset_id ON finetunes(training_dataset_id);
CREATE INDEX idx_finetunes_base_model_name ON finetunes(base_model_name);

-- Create trigger to update updated_at column
CREATE TRIGGER update_finetunes_updated_at BEFORE UPDATE ON finetunes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();