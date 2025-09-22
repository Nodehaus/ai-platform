-- Create training_data_items table
CREATE TABLE training_data_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    training_dataset_id UUID NOT NULL REFERENCES training_datasets(id) ON DELETE CASCADE,
    values_json TEXT NOT NULL,
    corrects_id UUID REFERENCES training_data_items(id) ON DELETE SET NULL,
    source_document TEXT,
    source_document_start TEXT,
    source_document_end TEXT,
    generation_time_seconds DECIMAL(10,2) NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_training_data_items_training_dataset_id ON training_data_items(training_dataset_id);
CREATE INDEX idx_training_data_items_corrects_id ON training_data_items(corrects_id);
CREATE INDEX idx_training_data_items_deleted ON training_data_items(deleted);

-- Create trigger to update updated_at column
CREATE TRIGGER update_training_data_items_updated_at BEFORE UPDATE ON training_data_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();