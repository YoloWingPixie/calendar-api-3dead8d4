-- Calendar API Database Schema
-- PostgreSQL 16
--
-- This reference schema defines the database structure for the Calendar API.
-- It will be managed via SQLAlchemy and Alembic migrations.

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    access_key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for API key lookups
CREATE INDEX idx_users_access_key ON users(access_key);

-- Calendars Table
CREATE TABLE calendars (
    calendar_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id UUID NOT NULL REFERENCES users(user_id),
    calendar_name VARCHAR(255) NOT NULL,
    editor_ids UUID[] NOT NULL DEFAULT '{}',
    reader_ids UUID[] NOT NULL DEFAULT '{}',
    public_read BOOLEAN NOT NULL DEFAULT false,
    public_write BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_editor_ids CHECK (
        NOT EXISTS (
            SELECT 1 FROM unnest(editor_ids) AS editor_id
            WHERE NOT EXISTS (
                SELECT 1 FROM users WHERE user_id = editor_id
            )
        )
    ),
    CONSTRAINT valid_reader_ids CHECK (
        NOT EXISTS (
            SELECT 1 FROM unnest(reader_ids) AS reader_id
            WHERE NOT EXISTS (
                SELECT 1 FROM users WHERE user_id = reader_id
            )
        )
    )
);

-- Index for owner lookups
CREATE INDEX idx_calendars_owner_user_id ON calendars(owner_user_id);

-- Index for array contains operations (GIN index for array columns)
CREATE INDEX idx_calendars_editor_ids ON calendars USING GIN (editor_ids);
CREATE INDEX idx_calendars_reader_ids ON calendars USING GIN (reader_ids);

-- Calendar Events Table
CREATE TABLE calendar_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES calendars(calendar_id),
    creator_user_id UUID NOT NULL REFERENCES users(user_id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_all_day BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT end_after_start CHECK (end_time > start_time),
    CONSTRAINT min_duration CHECK (end_time - start_time >= interval '1 minute'),
    CONSTRAINT all_day_constraint CHECK (
        (is_all_day = false) OR
        (is_all_day = true AND
         start_time::time = '00:00:00' AND
         end_time::time = '23:59:59' AND
         start_time::date = end_time::date)
    )
);

-- Index for calendar lookups
CREATE INDEX idx_calendar_events_calendar_id ON calendar_events(calendar_id);

-- Index for creator lookups
CREATE INDEX idx_calendar_events_creator_user_id ON calendar_events(creator_user_id);

-- Index for time-based queries
CREATE INDEX idx_calendar_events_start_time ON calendar_events(start_time);
CREATE INDEX idx_calendar_events_end_time ON calendar_events(end_time);

-- Composite index for common time range queries
CREATE INDEX idx_calendar_events_time_range ON calendar_events(calendar_id, start_time, end_time);

-- Update trigger function for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for all tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_calendars_updated_at BEFORE UPDATE ON calendars
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_calendar_events_updated_at BEFORE UPDATE ON calendar_events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
