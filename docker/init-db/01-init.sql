-- PostgreSQL initialization script for Calendar API
-- This script runs when the database container is first created

-- Create extensions that might be useful
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Set some performance-oriented settings for development
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 1000; -- Log queries > 1 second

-- Log successful initialization
DO $$
BEGIN
    RAISE NOTICE 'Calendar API database initialized successfully';
END
$$; 