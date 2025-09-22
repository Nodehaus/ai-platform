-- Remove data_json column from training_datasets table since we now have a separate table
ALTER TABLE training_datasets DROP COLUMN data_json;