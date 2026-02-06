#!/bin/bash

LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)
URL="http://$LOCAL_IP:8080"

echo "=========================================="
echo "   Quick Mobile Access"
echo "=========================================="
echo ""
echo "📱 Scan this URL with your phone:"
echo ""
echo "   $URL"
echo ""
echo "Or manually type in your mobile browser:"
echo "   $LOCAL_IP:8080"
echo ""
echo "✅ Make sure your phone is on the same WiFi!"
echo ""
