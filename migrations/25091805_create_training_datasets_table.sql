-- Create training_datasets table
CREATE TABLE training_datasets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    generate_model VARCHAR(255),
    generate_model_runner VARCHAR(255),
    generate_gpu_info_card VARCHAR(255),
    generate_gpu_info_total_gb DECIMAL(10,2),
    generate_gpu_info_cuda_version VARCHAR(50),
    input_field VARCHAR(255) NOT NULL,
    output_field VARCHAR(255) NOT NULL,
    total_generation_time_seconds DECIMAL(10,2),
    generate_prompt_history_ids_json TEXT NOT NULL DEFAULT '[]',
    generate_prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE RESTRICT,
    corpus_id UUID NOT NULL REFERENCES corpus(id) ON DELETE RESTRICT,
    language_iso CHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PLANNING' CHECK (status IN ('PLANNING', 'RUNNING', 'ABORTED', 'FAILED', 'DONE', 'DELETED')),
    field_names_json TEXT NOT NULL,
    data_json TEXT NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Composite unique constraint for project versioning
    UNIQUE(project_id, version)
);

-- Create indexes for performance
CREATE INDEX idx_training_datasets_project_id ON training_datasets(project_id);
CREATE INDEX idx_training_datasets_status ON training_datasets(status);
CREATE INDEX idx_training_datasets_version ON training_datasets(project_id, version);
CREATE INDEX idx_training_datasets_corpus_id ON training_datasets(corpus_id);
CREATE INDEX idx_training_datasets_prompt_id ON training_datasets(generate_prompt_id);
CREATE INDEX idx_training_datasets_language_iso ON training_datasets(language_iso);

-- Create trigger to update updated_at column
CREATE TRIGGER update_training_datasets_updated_at BEFORE UPDATE ON training_datasets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();