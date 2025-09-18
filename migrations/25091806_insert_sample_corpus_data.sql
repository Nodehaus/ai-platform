-- Insert sample corpus data
-- These are example corpus entries that can be used for testing and development

INSERT INTO corpus (id, name, s3_path, created_at, updated_at) VALUES
    (uuid_generate_v4(), 'eurlex', 's3://nodehaus/documents/eurlex', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)

-- Note: In production, corpus data would be managed through a separate data ingestion process
-- This migration provides sample data for development and testing purposes