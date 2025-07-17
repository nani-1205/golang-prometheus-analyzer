// web/static/js/app.js
document.addEventListener('DOMContentLoaded', function() {
    const reportsContainer = document.getElementById('reports-container');
    const runSpikeAnalysisBtn = document.getElementById('run-spike-analysis');
    const runHighLoadAnalysisBtn = document.getElementById('run-high-load-analysis');

    const severityClasses = {
        'critical': 'danger',
        'warning': 'warning',
        'info': 'info',
    };

    function createReportCard(report) {
        const severityClass = severityClasses[report.severity] || 'secondary';
        const startTime = new Date(report.start_time).toLocaleString();
        const endTime = new Date(report.end_time).toLocaleString();
        const timeWindow = startTime === endTime ? startTime : `${startTime} to ${endTime}`;

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
                        <li class="list-group-item"><strong>Time:</strong> ${timeWindow}</li>
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

    // --- UPDATED BUTTON HANDLERS ---

    // Handler for the original "Spike" analysis
    runSpikeAnalysisBtn.addEventListener('click', async () => {
        runSpikeAnalysisBtn.disabled = true;
        runSpikeAnalysisBtn.innerHTML = '<span class="spinner-border spinner-border-sm"></span> Analyzing...';
        
        try {
            const response = await fetch('/api/analyze/cpu-spike', { method: 'POST' });
            if (response.status === 202) {
                alert('Spike analysis started! Reports will refresh soon.');
                setTimeout(fetchReports, 5000);
            } else {
                alert('Failed to start spike analysis.');
            }
        } catch (error) {
            console.error('Error triggering spike analysis:', error);
        } finally {
            runSpikeAnalysisBtn.disabled = false;
            runSpikeAnalysisBtn.innerHTML = '<i class="bi bi-activity"></i> Run Spike Analysis';
        }
    });

    // Handler for the NEW "High Load" analysis
    runHighLoadAnalysisBtn.addEventListener('click', async () => {
        runHighLoadAnalysisBtn.disabled = true;
        runHighLoadAnalysisBtn.innerHTML = '<span class="spinner-border spinner-border-sm"></span> Analyzing...';
        
        try {
            const response = await fetch('/api/analyze/cpu-high-load', { method: 'POST' });
            if (response.status === 202) {
                alert('High load analysis started! Reports will refresh soon.');
                setTimeout(fetchReports, 5000);
            } else {
                alert('Failed to start high load analysis.');
            }
        } catch (error) {
            console.error('Error triggering high load analysis:', error);
        } finally {
            runHighLoadAnalysisBtn.disabled = false;
            runHighLoadAnalysisBtn.innerHTML = '<i class="bi bi-fire"></i> Run High Load Analysis';
        }
    });


    // Initial load
    fetchReports();
});