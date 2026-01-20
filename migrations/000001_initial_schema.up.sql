-- Migration: 001_initial_schema.sql
-- Create features table

CREATE TABLE IF NOT EXISTS features (
    id SERIAL PRIMARY KEY,
    icon TEXT,
    title TEXT,
    description TEXT
);

-- Insert initial data
INSERT INTO features (icon, title, description) VALUES 
('💻', 'Web Development', 'Kami membuat website yang responsif dan modern.'),
('📱', 'Mobile Apps', 'Aplikasi mobile untuk iOS dan Android.'),
('☁️', 'Cloud Solutions', 'Solusi cloud untuk bisnis Anda.'),
('🔒', 'Cybersecurity', 'Melindungi data dan sistem Anda.')
ON CONFLICT DO NOTHING;