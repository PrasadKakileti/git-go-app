document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const email = document.getElementById('email').value;
    const location = document.getElementById('location').value;
    const frequency = document.getElementById('frequency').value;
    const messageDiv = document.getElementById('message');
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                location: location,
                notification_frequency: frequency
            })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            messageDiv.className = 'message success';
            messageDiv.textContent = 'Registration successful! You will receive job updates at ' + email;
            document.getElementById('registerForm').reset();
        } else {
            messageDiv.className = 'message error';
            messageDiv.textContent = data.message || 'Registration failed. Please try again.';
        }
    } catch (error) {
        messageDiv.className = 'message error';
        messageDiv.textContent = 'Network error. Please try again.';
    }
});
