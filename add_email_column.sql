-- SQL Script to Add Email Column to users_master Table
-- Database: innosense

USE innosense;

-- Check if email column already exists
-- If the table doesn't have an email column, add it
-- If it already exists, this will show an error which can be ignored

-- Option 1: Add email column if it doesn't exist (MySQL 8.0+)
-- This will only add the column if it doesn't already exist
ALTER TABLE users_master 
ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE AFTER id;

-- Option 2: For older MySQL versions, use this approach
-- First check if column exists, then add if needed
-- (This requires manual verification)

-- Add email column (uncomment if needed)
-- ALTER TABLE users_master 
-- ADD COLUMN email VARCHAR(255) UNIQUE AFTER id;

-- Make email column NOT NULL if needed (after adding data)
-- ALTER TABLE users_master 
-- MODIFY COLUMN email VARCHAR(255) NOT NULL UNIQUE;

-- Add index on email column for better performance
CREATE INDEX IF NOT EXISTS idx_users_master_email ON users_master(email);

-- Verify the table structure
DESCRIBE users_master;

-- Show current table structure
SHOW CREATE TABLE users_master;
