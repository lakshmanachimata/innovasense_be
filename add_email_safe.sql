-- Safe SQL Script to Add Email Column (handles existing column)
-- Database: innosense

USE innosense;

-- Check if email column exists and add if it doesn't
-- This approach works for all MySQL versions

-- Step 1: Check if email column exists
SELECT COUNT(*) as email_column_exists
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'innosense' 
  AND TABLE_NAME = 'users_master' 
  AND COLUMN_NAME = 'email';

-- Step 2: Add email column only if it doesn't exist
-- (Run this only if the above query returns 0)
-- ALTER TABLE users_master 
-- ADD COLUMN email VARCHAR(255) UNIQUE AFTER id;

-- Step 3: Add index for better performance
-- CREATE INDEX idx_users_master_email ON users_master(email);

-- Step 4: Verify the table structure
DESCRIBE users_master;
