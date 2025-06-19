# syntax=docker/dockerfile:1

FROM python:3.12-slim AS runtime

WORKDIR /app

# Install python dependencies first for better cache utilisation
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . ./

# Expose FastAPI default port
ENV PORT=8000
EXPOSE 8000

# Run the service
CMD ["python", "main.py"]
