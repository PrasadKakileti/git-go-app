#!/bin/bash

echo "Starting server..."
go run main.go > /tmp/jobhub_debug.log 2>&1 &
SERVER_PID=$!
sleep 3

echo ""
echo "Testing registration with your email..."
echo ""

RESPONSE=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"prasadkakileti105@gmail.com","location":"Hyderabad","domain":"Java","notification_frequency":"daily"}')

echo "Response: $RESPONSE"
echo ""

if echo "$RESPONSE" | grep -q "already registered"; then
    echo "✅ Email is already registered!"
    echo "   This is expected - you're already in the system."
    echo ""
    echo "Try with a different email or check admin panel:"
    echo "   http://localhost:8080/admin"
elif echo "$RESPONSE" | grep -q "success"; then
    echo "✅ Registration successful!"
else
    echo "❌ Error occurred:"
    echo "$RESPONSE"
    echo ""
    echo "Server logs:"
    tail -10 /tmp/jobhub_debug.log
fi

echo ""
kill $SERVER_PID 2>/dev/null
