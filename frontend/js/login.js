document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const emailOrMobile = document.getElementById('emailOrMobile').value.trim();
    const password = document.getElementById('password').value;
    const messageDiv = document.getElementById('message');
    
    const data = {
        emailOrMobile: emailOrMobile,
        password: password
    };
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok && result.success) {
            messageDiv.className = 'message success';
            messageDiv.textContent = 'Login successful! Redirecting...';
            messageDiv.style.display = 'block';
            
            // Store session
            localStorage.setItem('user', JSON.stringify(result.user));
            localStorage.setItem('userEmail', result.user.email);
            localStorage.setItem('token', result.token);
            
            // Redirect to dashboard
            window.location.href = '/dashboard.html';
        } else {
            throw new Error(result.error || 'Login failed');
        }
    } catch (error) {
        messageDiv.className = 'message error';
        messageDiv.textContent = error.message;
        messageDiv.style.display = 'block';
    }
});
