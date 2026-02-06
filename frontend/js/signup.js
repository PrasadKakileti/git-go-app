document.getElementById('signupForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const emailOrMobile = document.getElementById('emailOrMobile').value.trim();
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const location = document.getElementById('location').value;
    const domain = document.getElementById('domain').value;
    const experience = document.getElementById('experience').value;
    const frequency = document.querySelector('input[name="frequency"]:checked').value;
    const messageDiv = document.getElementById('message');
    
    // Validate passwords match
    if (password !== confirmPassword) {
        messageDiv.className = 'message error';
        messageDiv.textContent = 'Passwords do not match!';
        messageDiv.style.display = 'block';
        return;
    }
    
    // Determine if email or mobile
    const isEmail = emailOrMobile.includes('@');
    const isMobile = /^[0-9]{10}$/.test(emailOrMobile);
    
    if (!isEmail && !isMobile) {
        messageDiv.className = 'message error';
        messageDiv.textContent = 'Please enter a valid email or 10-digit mobile number';
        messageDiv.style.display = 'block';
        return;
    }
    
    const data = {
        email: isEmail ? emailOrMobile : '',
        mobile: isMobile ? emailOrMobile : '',
        password: password,
        location: location,
        domain: domain,
        experience: experience,
        notification_frequency: frequency
    };
    
    const submitBtn = document.querySelector('.btn-primary');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Signing up...';
    
    try {
        const response = await fetch('/api/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok && result.success) {
            messageDiv.className = 'message success';
            messageDiv.textContent = result.message;
            messageDiv.style.display = 'block';
            
            // Redirect to login after 2 seconds
            setTimeout(() => {
                window.location.href = 'login.html';
            }, 2000);
        } else {
            throw new Error(result.error || 'Signup failed');
        }
    } catch (error) {
        messageDiv.className = 'message error';
        messageDiv.textContent = error.message;
        messageDiv.style.display = 'block';
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Sign Up';
    }
});
