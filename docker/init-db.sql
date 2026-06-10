-- Creates the three service databases on a single PostgreSQL instance.
-- Runs automatically on first container start via docker-entrypoint-initdb.d.
-- The default POSTGRES_DB (spotnearr_user) is created by the postgres image itself.

CREATE DATABASE spotnearr_vendor;
CREATE DATABASE spotnearr_search;
