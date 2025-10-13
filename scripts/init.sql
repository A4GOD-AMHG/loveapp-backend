-- LoveApp Database Initialization Script
-- This script sets up the initial database configuration

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC';

-- Create custom types if needed
-- (Currently not needed, but can be added here)

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON DATABASE loveapp TO loveapp;

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'LoveApp database initialized successfully at %', NOW();
END $$;