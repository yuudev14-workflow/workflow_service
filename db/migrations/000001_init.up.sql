-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable gen_salt extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create tables
CREATE TABLE IF NOT EXISTS
  users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT
  );
  
CREATE TABLE
  IF NOT EXISTS workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(50),
    description TEXT,
    trigger_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
  );

CREATE TABLE
  IF NOT EXISTS schedulers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    workflow_id UUID REFERENCES workflows (id) ON DELETE CASCADE,
    cron VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
  );

CREATE TABLE
  IF NOT EXISTS workflow_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    workflow_id UUID REFERENCES workflows (id) ON DELETE CASCADE,
    status VARCHAR(20),
    triggered_at TIMESTAMP NOT NULL DEFAULT NOW()
  );

CREATE TABLE
  IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    workflow_id UUID NOT NULL REFERENCES workflows (id) ON DELETE CASCADE,
    name VARCHAR(50),
    description TEXT,
    parameters JSON,
    config VARCHAR(50),
    connector_name VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_workflow_name UNIQUE (workflow_id, name)
  );

CREATE TABLE
  IF NOT EXISTS edges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    destination_id UUID REFERENCES tasks (id) ON DELETE CASCADE,
    source_id UUID REFERENCES tasks (id) ON DELETE CASCADE,

    CONSTRAINT unique_source_destination UNIQUE (source_id, destination_id)
  );

CREATE TABLE
  IF NOT EXISTS task_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    workflow_history_id UUID REFERENCES workflow_history (id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks (id) ON DELETE CASCADE,
    status VARCHAR(20),
    triggered_at TIMESTAMP NOT NULL DEFAULT NOW()
  );