// Smooth scrolling
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
    });
});

// Animated counter
function animateCounter(element, target, duration = 2000) {
    let start = 0;
    const increment = target / (duration / 16);
    const timer = setInterval(() => {
        start += increment;
        if (start >= target) {
            element.textContent = target + '+';
            clearInterval(timer);
        } else {
            element.textContent = Math.floor(start) + '+';
        }
    }, 16);
}

// Trigger counter animation when in view
const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            const counter = entry.target;
            const target = parseInt(counter.textContent);
            animateCounter(counter, target);
            observer.unobserve(counter);
        }
    });
});

const jobCounter = document.getElementById('jobCount');
if (jobCounter) {
    observer.observe(jobCounter);
}

// Form handling
const form = document.getElementById('registerForm');
const messageDiv = document.getElementById('message');

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const submitBtn = form.querySelector('.submit-btn');
    const btnText = submitBtn.querySelector('.btn-text');
    const originalText = btnText.textContent;
    
    // Disable button and show loading
    submitBtn.disabled = true;
    btnText.textContent = 'Processing...';
    messageDiv.style.display = 'none';
    
    const formData = {
        email: document.getElementById('email').value,
        location: document.getElementById('location').value,
        domain: document.getElementById('domain').value,
        notification_frequency: document.querySelector('input[name="frequency"]:checked').value
    };
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });
        
        const data = await response.json();
        
        if (response.ok && data.success) {
            messageDiv.className = 'message success';
            messageDiv.innerHTML = `
                <strong>🎉 Success!</strong><br>
                ${data.message}<br>
                Check <strong>${formData.email}</strong> for job alerts!
            `;
            form.reset();
            
            // Confetti effect (simple)
            createConfetti();
        } else {
            throw new Error(data.error || data.message || 'Registration failed');
        }
    } catch (error) {
        messageDiv.className = 'message error';
        messageDiv.innerHTML = `
            <strong>⚠️ Oops!</strong><br>
            ${error.message}. Please try again.
        `;
    } finally {
        submitBtn.disabled = false;
        btnText.textContent = originalText;
        messageDiv.style.display = 'block';
        messageDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
});

// Simple confetti effect
function createConfetti() {
    const colors = ['#6366f1', '#8b5cf6', '#10b981', '#f59e0b'];
    for (let i = 0; i < 50; i++) {
        setTimeout(() => {
            const confetti = document.createElement('div');
            confetti.style.cssText = `
                position: fixed;
                width: 10px;
                height: 10px;
                background: ${colors[Math.floor(Math.random() * colors.length)]};
                left: ${Math.random() * 100}vw;
                top: -10px;
                opacity: 1;
                border-radius: 50%;
                pointer-events: none;
                z-index: 9999;
                animation: fall ${2 + Math.random() * 2}s linear forwards;
            `;
            document.body.appendChild(confetti);
            setTimeout(() => confetti.remove(), 4000);
        }, i * 30);
    }
}

// Add fall animation
const style = document.createElement('style');
style.textContent = `
    @keyframes fall {
        to {
            transform: translateY(100vh) rotate(360deg);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// Navbar scroll effect
let lastScroll = 0;
const navbar = document.querySelector('.navbar');

window.addEventListener('scroll', () => {
    const currentScroll = window.pageYOffset;
    
    if (currentScroll > lastScroll && currentScroll > 100) {
        navbar.style.transform = 'translateY(-100%)';
    } else {
        navbar.style.transform = 'translateY(0)';
    }
    
    lastScroll = currentScroll;
});
