-- Migration: 000002_add_testimonials.up.sql
-- Create testimonials table and insert initial data

CREATE TABLE testimonials (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    company TEXT NOT NULL,
    text TEXT NOT NULL,
    avatar TEXT
);

INSERT INTO testimonials (name, company, text, avatar) VALUES
('Budi Santos', 'ABC Corporation', 'Pelayanan sangat memuaskan, website kami jadi lebih modern.', '👨'),
('Sari Dewi', 'XYZ Enterprises', 'Tim yang profesional dan hasil kerja berkualitas tinggi.', '👩');