-- Add generate_examples_number column to training_datasets table
ALTER TABLE training_datasets
ADD COLUMN generate_examples_number INTEGER NOT NULL DEFAULT 0;