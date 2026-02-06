#!/bin/bash

LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)

echo "Starting Job Portal Application..."
echo "=================================="
echo ""
echo "Database: job_portal (MySQL)"
echo ""
echo "📱 Access from ANY device:"
echo "   http://$LOCAL_IP:8080"
echo ""
echo "💻 Local access:"
echo "   http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run main.go
