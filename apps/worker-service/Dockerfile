FROM python:3.13-slim

WORKDIR /app

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py ./

ENV FLASK_ENV=production
EXPOSE 5000

CMD ["gunicorn", "-b", "0.0.0.0:5000", "app:app"]