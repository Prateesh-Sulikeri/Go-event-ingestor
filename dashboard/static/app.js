let chart;

// Fetch total event count
async function fetchTotal() {
  const r = await fetch("/metrics/total_events");
  const data = await r.json();
  document.getElementById("total").textContent = data.total;
}

// Fetch events per minute (last 10min)
async function fetchRecent() {
  const r = await fetch("/metrics/events_last_10");
  return r.json();
}

// Chart update logic
async function updateChart() {
  const data = await fetchRecent();
  const labels = data.length > 0 ? data.map(d => d.minute) : ["No Data"];
  const values = data.length > 0 ? data.map(d => d.count) : [0];

  const canvas = document.getElementById("eventsChart");
  const ctx = canvas.getContext("2d");

  if (!chart) {
    chart = new Chart(ctx, {
      type: "line",
      data: {
        labels,
        datasets: [{
          label: "Events per Min",
          data: values,
          borderColor: "#2563eb",
          backgroundColor: "rgba(37, 99, 235, 0.25)",
          borderWidth: 2,
          tension: 0.3,
          pointRadius: 3
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: false,
        scales: {
          x: {
            title: { display: true, text: "Time (HH:MM)" }
          },
          y: {
            beginAtZero: true,
            title: { display: true, text: "Events/Minute" }
          }
        }
      }
    });
  } else {
    chart.data.labels = labels;
    chart.data.datasets[0].data = values;
    chart.update("none");
  }
}


// ------------------------------
// Controller Actions
// ------------------------------

async function customSteady() {
  const rps = document.getElementById("steady_rps").value;
  const res = await fetch("/action/run", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({ action: "steady", rps })
  });
  const msg = await res.json();
  document.getElementById("status").textContent = msg.status;
}

async function customBurst() {
  const count = document.getElementById("burst_count").value;
  const res = await fetch("/action/run", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({ action: "burst", count })
  });
  document.getElementById("status").textContent = (await res.json()).status;
}

async function customWave() {
  const min = document.getElementById("wave_min").value;
  const max = document.getElementById("wave_max").value;
  const res = await fetch("/action/run", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({ action: "wave", min, max })
  });
  document.getElementById("status").textContent = (await res.json()).status;
}

async function triggerStop() {
  const r = await fetch("/action/run", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({ action: "stop" })
  });
  document.getElementById("status").textContent = (await r.json()).status;
}

async function clearDB() {
  const r = await fetch("/admin/clear_events", { method: "POST" });
  document.getElementById("status").textContent = (await r.json()).status;
  setTimeout(() => {
    fetchTotal();
    updateChart();
  }, 600);
}


// ------------------------------
// AUTO POLLING (REALTIME)
// ------------------------------

setInterval(() => {
  fetchTotal();
  updateChart();
}, 1500);


// Initial Load
fetchTotal();
updateChart();
