const userEmail = localStorage.getItem('userEmail');
if (!userEmail) {
    window.location.href = '/login.html';
}
document.getElementById('userEmail').textContent = userEmail;

// ── Job loading ────────────────────────────────────────────────
async function loadJobs() {
    document.getElementById('loading').style.display = 'block';
    document.getElementById('error').style.display = 'none';

    try {
        const res = await fetch(`/api/jobs?email=${encodeURIComponent(userEmail)}`);

        // Guard: make sure it's JSON before calling .json()
        const ct = res.headers.get('content-type') || '';
        if (!ct.includes('application/json')) {
            throw new Error(`Server error (status ${res.status})`);
        }

        const data = await res.json();
        document.getElementById('loading').style.display = 'none';

        if (!res.ok) throw new Error(data.error || 'Failed to load jobs');

        displayJobs(data.jobs || []);
    } catch (err) {
        document.getElementById('loading').style.display = 'none';
        document.getElementById('error').textContent = err.message;
        document.getElementById('error').style.display = 'block';
    }
}

// ── Refresh button ─────────────────────────────────────────────
async function refreshJobs() {
    const btn = document.getElementById('refreshBtn');
    btn.disabled = true;
    btn.textContent = '⟳ Fetching…';

    try {
        const res = await fetch(`/api/refresh-jobs?email=${encodeURIComponent(userEmail)}`);
        const data = await res.json();
        showToast(data.message || 'Fetching fresh jobs — reload in ~30 seconds.');
        // Auto-reload the job list after 30 s
        setTimeout(loadJobs, 30000);
    } catch {
        showToast('Could not connect to server.');
    } finally {
        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = '⟳ Refresh Jobs';
        }, 31000);
    }
}

// ── Display ────────────────────────────────────────────────────
function displayJobs(jobs) {
    const container = document.getElementById('jobsContainer');
    const statsBar  = document.getElementById('statsBar');

    if (jobs.length === 0) {
        statsBar.style.display = 'none';
        container.innerHTML = `
          <div class="no-jobs">
            <h3>No jobs found yet</h3>
            <p>We haven't fetched jobs for your profile yet, or none match your exact preferences right now.</p>
            <p class="hint">
              👉 If you just signed up, click <strong>⟳ Refresh Jobs</strong> above and wait ~30 seconds.<br>
              👉 Make sure your <strong>JSEARCH_API_KEY</strong> is set in <code>.env</code> —
              get a free key at <a href="https://rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch" target="_blank">RapidAPI JSearch</a>.
            </p>
          </div>`;
        return;
    }

    statsBar.style.display = 'flex';
    document.getElementById('jobCount').textContent = jobs.length;

    container.innerHTML = jobs.map(job => {
        const source   = job.source || 'Job Board';
        const applyURL = job.source_url && job.source_url.startsWith('http') ? job.source_url : '';

        return `
        <div class="job-card">
          <div class="job-source-badge">${escapeHtml(source)}</div>
          <h3 class="job-title">${escapeHtml(job.title)}</h3>
          <div class="job-company">🏢 ${escapeHtml(job.company || '')}</div>
          <div class="job-meta">
            <span>📍 ${escapeHtml(job.location)}</span>
            <span>💼 ${escapeHtml(job.domain)}</span>
          </div>
          <p class="job-description">${escapeHtml(job.description || '')}</p>
          <div class="job-date">📅 Posted: ${formatDate(job.posted_at)}</div>
          ${applyURL
            ? `<a href="${applyURL}" target="_blank" rel="noopener noreferrer" class="btn-apply">Apply Now →</a>`
            : `<span class="btn-apply disabled">Link unavailable</span>`}
        </div>`;
    }).join('');
}

// ── Helpers ────────────────────────────────────────────────────
function formatDate(ds) {
    if (!ds) return 'Unknown';
    const d = new Date(ds);
    return isNaN(d) ? ds : d.toLocaleDateString('en-IN', { year: 'numeric', month: 'short', day: 'numeric' });
}

function escapeHtml(text) {
    const d = document.createElement('div');
    d.textContent = String(text || '');
    return d.innerHTML;
}

function showToast(msg) {
    const t = document.getElementById('toast');
    t.textContent = msg;
    t.classList.add('show');
    setTimeout(() => t.classList.remove('show'), 4000);
}

function logout() {
    localStorage.removeItem('userEmail');
    localStorage.removeItem('userToken');
    window.location.href = '/login.html';
}

loadJobs();
