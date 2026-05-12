-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create Enums
CREATE TYPE role_type AS ENUM ('admin', 'karyawan');
CREATE TYPE leave_type AS ENUM ('planned', 'unplanned');

-- Create Toko Table
CREATE TABLE toko (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(255) NOT NULL,
    is_pusat BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Karyawan Table
CREATE TABLE karyawan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(255) NOT NULL,
    role role_type NOT NULL,
    toko_id UUID REFERENCES toko(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Jadwal Libur Table
CREATE TABLE jadwal_libur (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID REFERENCES karyawan(id),
    tanggal DATE NOT NULL,
    tipe leave_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Backup Assignment Table
CREATE TABLE backup_assignment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jadwal_libur_id UUID REFERENCES jadwal_libur(id) ON DELETE CASCADE,
    backup_karyawan_id UUID REFERENCES karyawan(id),
    assigned_by UUID REFERENCES karyawan(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
