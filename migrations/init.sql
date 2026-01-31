-- =============================================
-- Kasir API - Database Initialization Script
-- This script runs automatically when the 
-- PostgreSQL container starts for the first time
-- =============================================

-- Create category table
CREATE TABLE IF NOT EXISTS category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create product table
CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES category(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_product_category_id ON product(category_id);
CREATE INDEX IF NOT EXISTS idx_product_name ON product(name);
CREATE INDEX IF NOT EXISTS idx_category_name ON category(name);

-- Insert sample data for category
INSERT INTO category (name, description) VALUES 
    ('Makanan', 'Berbagai jenis makanan'),
    ('Minuman', 'Berbagai jenis minuman'),
    ('Snack', 'Makanan ringan dan cemilan')
ON CONFLICT DO NOTHING;

-- Insert sample data for product
INSERT INTO product (name, price, stock, category_id) VALUES 
    ('Nasi Goreng', 15000, 100, 1),
    ('Mie Ayam', 12000, 50, 1),
    ('Es Teh Manis', 5000, 200, 2),
    ('Kopi Hitam', 8000, 150, 2),
    ('Keripik Singkong', 10000, 80, 3),
    ('Kacang Goreng', 7000, 100, 3)
ON CONFLICT DO NOTHING;

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully!';
END $$;
