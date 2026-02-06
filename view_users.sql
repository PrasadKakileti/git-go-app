-- View all registered users
SELECT 
    id,
    email,
    location,
    notification_frequency,
    CASE WHEN is_active = 1 THEN 'Active' ELSE 'Inactive' END as status,
    created_at
FROM users
ORDER BY created_at DESC;

-- Count users by location
SELECT 
    location,
    COUNT(*) as user_count
FROM users
GROUP BY location
ORDER BY user_count DESC;

-- Count users by notification frequency
SELECT 
    notification_frequency,
    COUNT(*) as user_count
FROM users
GROUP BY notification_frequency;
