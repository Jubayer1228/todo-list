CREATE USER username with password 'pass';
CREATE DATABASE my_database;
GRANT ALL PRIVILEGES ON DATABASE my_database TO username;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS todos (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  text TEXT,
  checked BOOLEAN
);
