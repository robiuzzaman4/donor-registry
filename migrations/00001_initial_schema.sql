-- Up Migration
-- 1. Create Custom Enum Types
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_blood_group AS ENUM ('A+', 'A-', 'B+', 'B-', 'O+', 'O-', 'AB+', 'AB-');

CREATE TYPE user_role AS ENUM ('ADMIN', 'USER');

CREATE TYPE user_gender AS ENUM ('MALE', 'FEMALE');

-- 2. Create Users Table
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(11) NOT NULL UNIQUE,
        password TEXT NOT NULL,
        blood_group user_blood_group NOT NULL,
        role user_role NOT NULL DEFAULT 'USER',
        gender user_gender NOT NULL,
        date_of_birth DATE NOT NULL,
        zila VARCHAR(100) NOT NULL,
        upazila VARCHAR(100) NOT NULL,
        local_address TEXT NOT NULL,
        total_donate_count INT NOT NULL DEFAULT 0,
        is_verified BOOLEAN NOT NULL DEFAULT FALSE,
        is_available BOOLEAN NOT NULL DEFAULT TRUE,
        is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
        last_donated_at TIMESTAMP
        WITH
            TIME ZONE,
            created_at TIMESTAMP
        WITH
            TIME ZONE NOT NULL DEFAULT NOW (),
            updated_at TIMESTAMP
        WITH
            TIME ZONE NOT NULL DEFAULT NOW ()
    );

-- 3. Create indexes for common searches
CREATE INDEX idx_users_phone ON users (phone);

CREATE INDEX idx_users_blood_group ON users (blood_group);

CREATE INDEX idx_users_location ON users (zila, upazila);
