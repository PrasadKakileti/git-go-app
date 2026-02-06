// Check authentication
const userEmail = localStorage.getItem('userEmail');
if (!userEmail) {
    window.location.href = '/login.html';
}

document.getElementById('userEmail').textContent = userEmail;

// Fetch and display jobs
async function loadJobs() {
    try {
        const response = await fetch(`/api/jobs?email=${encodeURIComponent(userEmail)}`);
        const data = await response.json();
        
        document.getElementById('loading').style.display = 'none';
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load jobs');
        }
        
        displayJobs(data.jobs || []);
    } catch (error) {
        document.getElementById('loading').style.display = 'none';
        document.getElementById('error').textContent = error.message;
        document.getElementById('error').style.display = 'block';
    }
}

function displayJobs(jobs) {
    const container = document.getElementById('jobsContainer');
    
    if (jobs.length === 0) {
        container.innerHTML = '<div class="no-jobs"><h3>No jobs available yet</h3><p>Check back soon for new opportunities!</p></div>';
        return;
    }
    
    container.innerHTML = jobs.map(job => `
        <div class="job-card">
            <h3 class="job-title">${escapeHtml(job.title)}</h3>
            <div class="job-company">🏢 ${escapeHtml(job.company)}</div>
            <div class="job-meta">
                <span>📍 ${escapeHtml(job.location)}</span>
                <span>💼 ${escapeHtml(job.domain)}</span>
            </div>
            <p class="job-description">${escapeHtml(job.description)}</p>
            <div class="job-date">📅 Posted: ${formatDate(job.posted_at)}</div>
            <a href="${escapeHtml(job.source_url)}" target="_blank" class="btn-apply">View & Apply →</a>
        </div>
    `).join('');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function logout() {
    localStorage.removeItem('userEmail');
    window.location.href = '/login.html';
}

// Load jobs on page load
loadJobs();
