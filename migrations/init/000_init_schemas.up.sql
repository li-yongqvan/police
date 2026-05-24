-- ============================================================================
-- Migration 000: Initialize three independent PostgreSQL schemas
-- ============================================================================
-- Each service owns its own schema for data isolation.
-- MVP: shared PostgreSQL instance, separate schemas.
-- Later: migrate schemas to independent database instances.
-- ============================================================================

-- Create schema_auth for user-service
CREATE SCHEMA IF NOT EXISTS schema_auth;

-- Create schema_forum for forum-service
CREATE SCHEMA IF NOT EXISTS schema_forum;

-- Create schema_admin for admin-service
CREATE SCHEMA IF NOT EXISTS schema_admin;
