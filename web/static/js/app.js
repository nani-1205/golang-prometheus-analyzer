document.addEventListener('DOMContentLoaded', function() {
    const reportsContainer = document.getElementById('reports-container');
    const runCpuAnalysisBtn = document.getElementById('run-cpu-analysis');

    const severityClasses = {
        'critical': 'danger',
        'warning': 'warning',
        'info': 'info',
    };

    function createReportCard(report) {
        const severityClass = severityClasses[report.severity] || 'secondary';
        const startTime = new Date(report.start_time).toLocaleString();
        const endTime = new Date(report.end_time).toLocaleString();

        return `
            <div class="col-md-6 col-lg-4 mb-4">
                <div class="card h-100 shadow-sm severity-${report.severity}">
                    <div class="card-header bg-dark text-white">
                        <h5 class="card-title mb-0">${report.pattern_detected}</h5>
                        <small>${report.metric_name}</small>
                    </div>
                    <div class="card-body">
                        <p class="card-text">${report.details}</p>
                    </div>
                    <ul class="list-group list-group-flush">
                        <li class="list-group-item"><strong>Severity:</strong> <span class="badge bg-${severityClass}">${report.severity}</span></li>
                        <li class="list-group-item"><strong>Start:</strong> ${startTime}</li>
                        <li class="list-group-item"><strong>End:</strong> ${endTime}</li>
                    </ul>
                    <div class="card-footer text-muted">
                        Reported on ${new Date(report.created_at).toLocaleDateString()}
                    </div>
                </div>
            </div>
        `;
    }

    async function fetchReports() {
        try {
            const response = await fetch('/api/reports');
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const reports = await response.json();
            reportsContainer.innerHTML = ''; // Clear spinner
            if (reports && reports.length > 0) {
                reports.forEach(report => {
                    reportsContainer.innerHTML += createReportCard(report);
                });
            } else {
                reportsContainer.innerHTML = '<div class="col-12"><div class="alert alert-info">No reports found yet.</div></div>';
            }
        } catch (error) {
            console.error('Failed to fetch reports:', error);
            reportsContainer.innerHTML = '<div class="col-12"><div class="alert alert-danger">Failed to load reports.</div></div>';
        }
    }

    runCpuAnalysisBtn.addEventListener('click', async () => {
        runCpuAnalysisBtn.disabled = true;
        runCpuAnalysisBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Analyzing...';
        
        try {
            const response = await fetch('/api/analyze/cpu', { method: 'POST' });
            if (response.status === 202) {
                alert('Analysis started! The page will refresh with new reports soon.');
                // You could implement polling or WebSockets here for a smoother experience
                setTimeout(fetchReports, 5000); // Refresh reports after 5 seconds
            } else {
                alert('Failed to start analysis.');
            }
        } catch (error) {
            console.error('Error triggering analysis:', error);
            alert('An error occurred.');
        } finally {
            runCpuAnalysisBtn.disabled = false;
            runCpuAnalysisBtn.innerHTML = '<i class="bi bi-play-circle"></i> Run CPU Analysis Now';
        }
    });

    // Initial load
    fetchReports();
});