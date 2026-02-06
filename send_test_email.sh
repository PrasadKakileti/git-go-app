#!/bin/bash

echo "=========================================="
echo "   Sending Test Email"
echo "=========================================="
echo ""
echo "To: babusurendra500@gmail.com"
echo "From: prasadkakileti105@gmail.com"
echo ""

cd cmd/test_email && go run test_email.go

echo ""
echo "✅ Done! Check the inbox (and spam folder)"
echo ""
