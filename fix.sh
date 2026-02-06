#!/bin/bash

echo "🔧 Fixing common issues..."
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file missing - creating from template..."
    cp .env.example .env
    sed -i '' 's/DB_PASSWORD=/DB_PASSWORD=password/' .env
    echo "✅ .env file created"
else
    echo "✅ .env file exists"
fi

# Test MySQL connection
echo ""
echo "Testing MySQL connection..."
if mysql -u root -ppassword -e "SELECT 1;" 2>/dev/null | grep -q "1"; then
    echo "✅ MySQL connection successful"
else
    echo "❌ MySQL connection failed"
    echo "   Resetting MySQL password..."
    mysql -u root -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';" 2>/dev/null
    mysql -u root -ppassword -e "FLUSH PRIVILEGES;" 2>/dev/null
    echo "✅ MySQL password reset to 'password'"
fi

# Check database
echo ""
echo "Checking database..."
if mysql -u root -ppassword -e "USE job_portal;" 2>/dev/null; then
    echo "✅ Database 'job_portal' exists"
else
    echo "❌ Database missing - creating..."
    mysql -u root -ppassword < schema.sql 2>/dev/null
    echo "✅ Database created"
fi

echo ""
echo "🎉 All checks complete! Run: ./start.sh"
