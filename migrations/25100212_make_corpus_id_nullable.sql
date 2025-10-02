-- Make corpus_id nullable in training_datasets table
ALTER TABLE training_datasets ALTER COLUMN corpus_id DROP NOT NULL;
