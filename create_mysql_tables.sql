-- InnovoSens Database Schema for MySQL
-- Database: innosense

USE innosense;

-- Drop tables if they exist (in reverse dependency order)
DROP TABLE IF EXISTS org_users;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS device_master;
DROP TABLE IF EXISTS home_images;
DROP TABLE IF EXISTS banner_images;
DROP TABLE IF EXISTS sweat_images;
DROP TABLE IF EXISTS sweat_rate_summary;
DROP TABLE IF EXISTS sweat_summary;
DROP TABLE IF EXISTS sweat_data;
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users_master;

-- Create users_master table
CREATE TABLE users_master (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    cnumber VARCHAR(20),
    userpin VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    gender VARCHAR(10),
    age INT,
    height DECIMAL(5,2),
    weight DECIMAL(5,2),
    role_id INT DEFAULT 2,
    ustatus INT DEFAULT 0,
    creation_datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create user_data table for hydration records
CREATE TABLE user_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    weight DECIMAL(5,2),
    height DECIMAL(5,2),
    sweat_position DECIMAL(5,2),
    time_taken DECIMAL(5,2),
    bmi DECIMAL(5,2),
    tbsa DECIMAL(5,2),
    image_path TEXT,
    sweat_rate DECIMAL(5,2),
    sweat_loss DECIMAL(5,2),
    device_type INT,
    image_id INT,
    creation_datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users_master(id) ON DELETE CASCADE
);

-- Create sweat_data table
CREATE TABLE sweat_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    image_id INT NOT NULL,
    sweat_rate DECIMAL(5,2),
    sweat_loss DECIMAL(5,2),
    creation_datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users_master(id) ON DELETE CASCADE
);

-- Create sweat_summary table
CREATE TABLE sweat_summary (
    id INT AUTO_INCREMENT PRIMARY KEY,
    low_limit DECIMAL(5,2),
    high_limit DECIMAL(5,2),
    hyd_status VARCHAR(50),
    comments TEXT,
    recomm TEXT,
    color VARCHAR(20)
);

-- Create sweat_rate_summary table
CREATE TABLE sweat_rate_summary (
    id INT AUTO_INCREMENT PRIMARY KEY,
    low_limit DECIMAL(5,2),
    high_limit DECIMAL(5,2),
    hyd_status VARCHAR(50),
    comments TEXT,
    recomm TEXT,
    color VARCHAR(20)
);

-- Create sweat_images table
CREATE TABLE sweat_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    image_path TEXT NOT NULL,
    sweat_range VARCHAR(50),
    implications TEXT,
    recomm TEXT,
    strategy TEXT,
    result TEXT,
    colorcode VARCHAR(20)
);

-- Create banner_images table
CREATE TABLE banner_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    image_path VARCHAR(500) NOT NULL
);

-- Create home_images table
CREATE TABLE home_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    image_path VARCHAR(500) NOT NULL
);

-- Create device_master table
CREATE TABLE device_master (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_name VARCHAR(255) NOT NULL,
    device_text TEXT
);

-- Create organizations table
CREATE TABLE organizations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    org_name VARCHAR(255) NOT NULL,
    org_desc TEXT,
    salt_key VARCHAR(255),
    api_key VARCHAR(255) UNIQUE
);

-- Create org_users table
CREATE TABLE org_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email_id VARCHAR(255) NOT NULL,
    user_pwd VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    org_id INT NOT NULL,
    FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- Insert sample data
-- Insert banner images
INSERT INTO banner_images (image_path) VALUES 
('/assets/banners/1.png'),
('/assets/banners/2.png'),
('/assets/banners/3.jpg'),
('/assets/banners/4.jpg'),
('/assets/banners/5.jpg'),
('/assets/banners/banner1.png'),
('/assets/banners/banner2.png'),
('/assets/banners/banner3.png');

-- Insert home images
INSERT INTO home_images (image_path) VALUES 
('/assets/banners/1.png'),
('/assets/banners/2.png'),
('/assets/banners/3.jpg'),
('/assets/banners/4.jpg'),
('/assets/banners/5.jpg'),
('/assets/banners/banner1.png'),
('/assets/banners/banner2.png'),
('/assets/banners/banner3.png');

-- Insert devices
INSERT INTO device_master (device_name, device_text) VALUES 
('Hydrosense (Classic)', 'Standard sweat hydration tracking for average exercise and sweat loss.'),
('Hydrosense Plus (+)', 'High-volume hydration tracking for longer or more intense workouts.'),
('Hydrosense Pro', 'Sweat hydration + electrolyte monitoring for regular training sessions.'),
('Hydrosense Pro Plus (+)', 'Full hydration and electrolyte analysis for high-intensity or extended activity.');

-- Insert sample organization
INSERT INTO organizations (org_name, org_desc, salt_key, api_key) VALUES 
('Test Organization', 'Test organization for API testing', 'test-salt-key', 'test-api-key'),
('InnoSense Organization', 'Main organization for InnoSense API', 'innosense-salt-key-2024', 'innosense-api-key-2024');

-- Insert sample organization user
INSERT INTO org_users (email_id, user_pwd, user_name, org_id) VALUES 
('test@example.com', 'test123', 'Test User', 1),
('admin@innosense.com', 'admin123', 'Admin User', 2);

-- Insert sample sweat summary data
INSERT INTO sweat_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES 
(0.0, 0.5, 'Low', 'Low hydration level', 'Increase fluid intake', 'red'),
(0.5, 1.0, 'Normal', 'Normal hydration level', 'Maintain current fluid intake', 'green'),
(1.0, 2.0, 'High', 'High hydration level', 'Consider reducing fluid intake', 'yellow');

-- Insert sample sweat rate summary data
INSERT INTO sweat_rate_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES 
(0.0, 0.5, 'Low Rate', 'Low sweat rate', 'Increase activity intensity', 'blue'),
(0.5, 1.5, 'Normal Rate', 'Normal sweat rate', 'Maintain current activity level', 'green'),
(1.5, 3.0, 'High Rate', 'High sweat rate', 'Consider reducing intensity', 'orange');

-- Insert sample sweat images
INSERT INTO sweat_images (image_path, sweat_range, implications, recomm, strategy, result, colorcode) VALUES 
('/assets/sweat/1.jpg', '0.0-0.5', 'Low sweat level', 'Increase activity', 'Gradual intensity increase', 'Improved hydration', 'blue'),
('/assets/sweat/2.jpg', '0.5-1.0', 'Normal sweat level', 'Maintain current routine', 'Consistent training', 'Optimal hydration', 'green'),
('/assets/sweat/3.jpg', '1.0-2.0', 'High sweat level', 'Monitor hydration', 'Adjust intensity', 'Watch for overexertion', 'yellow'),
('/assets/sweat/4.jpg', '2.0+', 'Very high sweat level', 'Reduce intensity', 'Rest and rehydrate', 'Risk of dehydration', 'red'),
('/assets/sweat/5.jpg', '0.0-0.3', 'Very low sweat level', 'Increase activity significantly', 'High intensity training', 'Need more activity', 'purple'),
('/assets/sweat/6.jpg', '1.5-2.5', 'Moderate-high sweat level', 'Monitor closely', 'Balanced approach', 'Good progress', 'orange');

-- Show table creation status
SELECT 'Database and tables created successfully!' as status;
