-- Enable UUID extension if not enabled (it's already enabled in 000001, but just in case)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE konfigurasi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) UNIQUE NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO konfigurasi (key, value) VALUES 
('maks_libur_per_bulan', '3'),
('min_available_per_hari', '2');
