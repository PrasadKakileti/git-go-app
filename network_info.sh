#!/bin/bash

echo "=========================================="
echo "   JobHub - Network Access Information"
echo "=========================================="
echo ""

# Get local IP
LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)
PORT=8080

echo "📱 Access from ANY device on your network:"
echo ""
echo "   Main Site:    http://$LOCAL_IP:$PORT"
echo "   Admin Panel:  http://$LOCAL_IP:$PORT/admin"
echo ""
echo "💻 Local access (this computer):"
echo ""
echo "   Main Site:    http://localhost:$PORT"
echo "   Admin Panel:  http://localhost:$PORT/admin"
echo ""
echo "=========================================="
echo "   Instructions for Other Devices"
echo "=========================================="
echo ""
echo "1. Make sure your device is on the SAME WiFi network"
echo "2. Open browser on your phone/tablet/other computer"
echo "3. Enter: http://$LOCAL_IP:$PORT"
echo ""
echo "✅ Server is configured to accept connections from:"
echo "   - This computer (localhost)"
echo "   - Any device on your local network"
echo ""
echo "🔒 Security Note:"
echo "   Only devices on your WiFi can access this"
echo "   Not accessible from the internet"
echo ""
