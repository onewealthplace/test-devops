from flask import Flask, Response
from prometheus_client import Counter, generate_latest, CONTENT_TYPE_LATEST
import time
import logging

app = Flask(__name__)

# Global leak
leak = []

# Prometheus metrics
requests_total = Counter('worker_requests_total', 'Total number of /do-work requests')

@app.route('/do-work', methods=['POST'])
def do_work():
    app.logger.info("Processing /do-work request...")
    requests_total.inc()
    leak.append('XXXX' * 10**6)
    return f'Work done!', 200

@app.route('/clear', methods=['POST'])
def clear_leak():
    app.logger.info("Clearing leak...")
    leak.clear()
    return f'Leak cleared!', 200

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
    app.run(host='0.0.0.0', port=5000, debug=True)

if __name__ != '__main__':
    gunicorn_logger = logging.getLogger('gunicorn.error')
    app.logger.handlers = gunicorn_logger.handlers
    app.logger.setLevel(gunicorn_logger.level)
