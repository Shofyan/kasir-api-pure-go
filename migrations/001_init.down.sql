-- Migration: 001_init (Rollback)
-- Description: Rollback initial database schema

-- Drop indexes
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_categories_name;

-- Drop tables (order matters due to foreign key)
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
