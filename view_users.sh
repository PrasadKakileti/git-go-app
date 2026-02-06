#!/bin/bash

echo "======================================"
echo "   JobHub - Registered Users"
echo "======================================"
echo ""

mysql -u root -ppassword job_portal -e "
SELECT 
    id as 'ID',
    email as 'Email',
    location as 'Location',
    notification_frequency as 'Frequency',
    CASE WHEN is_active = 1 THEN 'Active' ELSE 'Inactive' END as 'Status',
    DATE_FORMAT(created_at, '%Y-%m-%d %H:%i') as 'Registered'
FROM users
ORDER BY created_at DESC;
" 2>/dev/null

echo ""
echo "======================================"
echo "   Statistics"
echo "======================================"

mysql -u root -ppassword job_portal -e "
SELECT 
    COUNT(*) as 'Total Users',
    SUM(CASE WHEN notification_frequency = 'daily' THEN 1 ELSE 0 END) as 'Daily',
    SUM(CASE WHEN notification_frequency = 'weekly' THEN 1 ELSE 0 END) as 'Weekly'
FROM users;
" 2>/dev/null

echo ""
