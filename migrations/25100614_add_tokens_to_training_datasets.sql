-- Add tokens_in and tokens_out fields to training_datasets table
ALTER TABLE training_datasets
ADD COLUMN tokens_in INTEGER,
ADD COLUMN tokens_out INTEGER;
