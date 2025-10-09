-- Create deployment_logs table
CREATE TABLE deployment_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deployment_id UUID NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,
    tokens_in INT NOT NULL,
    tokens_out INT NOT NULL,
    input TEXT NOT NULL,
    output TEXT NOT NULL,
    delay_time INT NOT NULL,
    execution_time INT NOT NULL,
    source VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_deployment_logs_deployment_id ON deployment_logs(deployment_id);
CREATE INDEX idx_deployment_logs_created_at ON deployment_logs(created_at);

-- Create trigger to update updated_at column
CREATE TRIGGER update_deployment_logs_updated_at BEFORE UPDATE ON deployment_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
