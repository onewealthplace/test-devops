from flask import Flask, Response
from prometheus_client import Counter, generate_latest, CONTENT_TYPE_LATEST
import time

app = Flask(__name__)

# Global leak
leak = []

# Prometheus metrics
requests_total = Counter('worker_requests_total', 'Total number of /do-work requests')

@app.route('/do-work', methods=['POST'])
def do_work():
    print("Processing /do-work request...")
    requests_total.inc()
    leak.append('XXXX' * 10**6)
    return f"Work done!"

@app.route('/metrics')
def metrics():
    return Response(generate_latest(), mimetype=CONTENT_TYPE_LATEST)

@app.route('/health')
def health_check():
    if all_required_services_are_running():
        return 'OK', 200
    else:
        return 'Service Unavailable', 500

def all_required_services_are_running():
    return True

if __name__ == '__main__':
    print("Lauching startup process...")
    time.sleep(10)
    print("Startup complete. Starting Flask app.")
    app.run(host='0.0.0.0', port=5000)
