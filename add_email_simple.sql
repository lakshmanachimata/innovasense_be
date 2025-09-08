-- Simple SQL Script to Add Email Column to users_master Table
-- Database: innosense

USE innosense;

-- Add email column to users_master table
ALTER TABLE users_master 
ADD COLUMN email VARCHAR(255) UNIQUE AFTER id;

-- Add index for better performance
CREATE INDEX idx_users_master_email ON users_master(email);

-- Verify the change
DESCRIBE users_master;
