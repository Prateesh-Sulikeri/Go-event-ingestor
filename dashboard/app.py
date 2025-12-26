import os
from flask import Flask, jsonify, render_template, request
from flask_cors import CORS
import psycopg2
import subprocess

# Resolve project root and scripts path
BASE_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
SHELL_DIR = os.path.join(BASE_DIR, "scripts", "shell")

app = Flask(
    __name__,
    static_folder=os.path.join(os.path.dirname(__file__), "static"),
    template_folder=os.path.join(os.path.dirname(__file__), "templates")
)

CORS(app)

# Database connection
def db():
    return psycopg2.connect(
        host=os.getenv("PG_HOST", "localhost"),
        database=os.getenv("PG_DATABASE", "event_ingestor"),
        user=os.getenv("PG_USER", "postgres"),
        password=os.getenv("PG_PASSWORD", ""),
    )


# ---------------------------------------------------------
# ROUTES
# ---------------------------------------------------------

@app.get("/")
def index():
    return render_template("index.html")


@app.get("/metrics/total_events")
def total_events():
    conn = db()
    cur = conn.cursor()
    cur.execute("SELECT COUNT(*) FROM events")
    count = cur.fetchone()[0]
    conn.close()
    return jsonify({"total": count})


@app.get("/metrics/events_last_10")
def events_last_10():
    conn = db()
    cur = conn.cursor()
    cur.execute("""
        SELECT to_char(date_trunc('minute', received_at), 'HH24:MI') AS minute,
               COUNT(*)
        FROM events
        WHERE received_at > NOW() - INTERVAL '10 minutes'
        GROUP BY 1
        ORDER BY 1
    """)
    rows = cur.fetchall()
    conn.close()
    return jsonify([{"minute": r[0], "count": r[1]} for r in rows])


@app.post("/admin/clear_events")
def clear_events():
    try:
        conn = db()
        cur = conn.cursor()
        cur.execute("TRUNCATE TABLE events RESTART IDENTITY;")
        conn.commit()
        return jsonify({"status": "ok"})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        conn.close()


@app.post("/action/run")
def run_action():
    data = request.json
    action = data.get("action")

    if action == "steady":
        rps = str(data.get("rps", 5))
        subprocess.Popen(["bash", f"{SHELL_DIR}/steady_load.sh", rps])
        return jsonify({"status": f"steady {rps} started"})

    if action == "burst":
        count = str(data.get("count", 200))
        subprocess.Popen(["bash", f"{SHELL_DIR}/burst_test.sh", count])
        return jsonify({"status": f"burst {count} started"})

    if action == "wave":
        min_r = str(data.get("min", 1))
        max_r = str(data.get("max", 20))
        subprocess.Popen(["bash", f"{SHELL_DIR}/wave.sh", min_r, max_r])
        return jsonify({"status": f"wave {min_r}-{max_r} started"})

    if action == "stop":
        subprocess.call(["pkill", "-f", "steady_load.sh"])
        subprocess.call(["pkill", "-f", "burst_test.sh"])
        subprocess.call(["pkill", "-f", "wave.sh"])
        return jsonify({"status": "all loads stopped"})

    return jsonify({"error": "invalid action"}), 400


# ---------------------------------------------------------
# MAIN
# ---------------------------------------------------------
if __name__ == "__main__":
    print("Dashboard running at http://localhost:5000")
    app.run(debug=True)
