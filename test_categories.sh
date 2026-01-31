#!/bin/bash

# Category API Testing Script
# Run this script to test all category endpoints

BASE_URL="http://127.0.0.1:8080"

echo "🧪 Testing Category API Endpoints"
echo "=================================="

echo ""
echo "1. GET /api/categories - Get all categories"
echo "------------------------------------------"
curl -X GET "${BASE_URL}/api/categories" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "2. POST /api/categories - Create category (JSON)"
echo "------------------------------------------------"
curl -X POST "${BASE_URL}/api/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Category",
    "description": "Test category description"
  }' \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "3. POST /api/categories - Create category (Form data)"
echo "----------------------------------------------------"
curl -X POST "${BASE_URL}/api/categories" \
  -F "name=Form Category" \
  -F "description=Created via form data" \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "4. GET /api/categories/1 - Get category by ID"
echo "--------------------------------------------"
curl -X GET "${BASE_URL}/api/categories/1" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "5. PUT /api/categories/1 - Update category"
echo "------------------------------------------"
curl -X PUT "${BASE_URL}/api/categories/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Category",
    "description": "Updated description"
  }' \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "6. DELETE /api/categories/1 - Delete category"
echo "---------------------------------------------"
curl -X DELETE "${BASE_URL}/api/categories/1" \
  -w "\nStatus: %{http_code}\n\n"

echo ""
echo "✅ Testing completed!"