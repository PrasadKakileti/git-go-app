#!/bin/bash

echo "=========================================="
echo "   JobHub System Test"
echo "=========================================="
echo ""

# Check if app is running
if lsof -ti:8080 > /dev/null 2>&1; then
    echo "✅ Application is running on port 8080"
else
    echo "❌ Application is NOT running"
    echo "   Run: ./start.sh"
    exit 1
fi

echo ""
echo "Testing components..."
echo ""

# Test database
echo "1. Database Connection:"
if mysql -u root -ppassword job_portal -e "SELECT COUNT(*) FROM users;" 2>/dev/null | grep -q "[0-9]"; then
    USER_COUNT=$(mysql -u root -ppassword job_portal -e "SELECT COUNT(*) FROM users;" 2>/dev/null | tail -1)
    echo "   ✅ Connected - $USER_COUNT users registered"
else
    echo "   ❌ Database connection failed"
fi

# Test API
echo ""
echo "2. API Health Check:"
if curl -s http://localhost:8080/api/health | grep -q "ok"; then
    echo "   ✅ API is responding"
else
    echo "   ❌ API not responding"
fi

# Test email configuration
echo ""
echo "3. Email Service:"
if grep -q "SENDGRID_API_KEY=SG\." .env 2>/dev/null; then
    echo "   ✅ SendGrid configured"
elif grep -q "SMTP_PASSWORD=..*" .env 2>/dev/null; then
    echo "   ✅ SMTP configured"
else
    echo "   ⚠️  No email service configured"
    echo "      See SENDGRID_SETUP.md"
fi

# Check jobs in database
echo ""
echo "4. Job Listings:"
JOB_COUNT=$(mysql -u root -ppassword job_portal -e "SELECT COUNT(*) FROM jobs;" 2>/dev/null | tail -1)
if [ "$JOB_COUNT" -gt 0 ]; then
    echo "   ✅ $JOB_COUNT jobs in database"
else
    echo "   ⚠️  No jobs yet (scraper will run every 6 hours)"
fi

echo ""
echo "=========================================="
echo "   Test Complete!"
echo "=========================================="
echo ""
echo "📱 Access URLs:"
echo "   Local:   http://localhost:8080"
echo "   Network: http://10.21.12.100:8080"
echo "   Admin:   http://localhost:8080/admin"
echo ""
