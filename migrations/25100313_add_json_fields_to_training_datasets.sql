-- Add json_object_fields and expected_output_size_chars to training_datasets table
ALTER TABLE training_datasets
ADD COLUMN json_object_fields_json TEXT,
ADD COLUMN expected_output_size_chars INTEGER;

-- Update existing rows with default values (empty JSON object and 0)
UPDATE training_datasets
SET json_object_fields_json = '{}',
    expected_output_size_chars = 0
WHERE json_object_fields_json IS NULL;

-- Make the columns NOT NULL after setting defaults
ALTER TABLE training_datasets
ALTER COLUMN json_object_fields_json SET NOT NULL,
ALTER COLUMN expected_output_size_chars SET NOT NULL;
