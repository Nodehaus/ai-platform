-- Insert sample corpus data
-- These are example corpus entries that can be used for testing and development

INSERT INTO corpus (id, name, s3_path, files_subset, created_at, updated_at) VALUES
    (uuid_generate_v4(), 'eurlex', '/documents/eurlex', ARRAY[
    '02016R0679-20160504_eng.json',
    '02002L0058-20091219_eng.json',
    '02016L0680-20160504_eng.json',
    '02023R1115-20241226_eng.json',
    '02010L0075-20240804_eng.json',
    '02006R1907-20250623_eng.json',
    '02008L0098-20180705_eng.json',
    '02014L0065-20250117_eng.json',
    '02013R0575-20250629_eng.json',
    '02013L0036-20250117_eng.json',
    '02015L2366-20250117_eng.json',
    '02016R1011-20220101_eng.json'
], CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)

-- Note: In production, corpus data would be managed through a separate data ingestion process
-- This migration provides sample data for development and testing purposes