-- Migration: 001_init
-- Description: Initial database schema for Kasir API
-- Created: 2026-01-31

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Insert sample data for categories
INSERT INTO categories (name, description) VALUES 
    ('Makanan', 'Berbagai jenis makanan'),
    ('Minuman', 'Berbagai jenis minuman'),
    ('Snack', 'Makanan ringan dan cemilan')
ON CONFLICT DO NOTHING;

-- Insert sample data for products
INSERT INTO products (name, price, stock, category_id) VALUES 
    ('Nasi Goreng', 15000, 100, 1),
    ('Mie Ayam', 12000, 50, 1),
    ('Es Teh Manis', 5000, 200, 2),
    ('Kopi Hitam', 8000, 150, 2),
    ('Keripik Singkong', 10000, 80, 3),
    ('Kacang Goreng', 7000, 100, 3)
ON CONFLICT DO NOTHING;
