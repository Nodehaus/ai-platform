-- Create corpus table
CREATE TABLE corpus (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    s3_path VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name for faster lookups
CREATE INDEX idx_corpus_name ON corpus(name);

-- Create trigger to update updated_at column
CREATE TRIGGER update_corpus_updated_at BEFORE UPDATE ON corpus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();