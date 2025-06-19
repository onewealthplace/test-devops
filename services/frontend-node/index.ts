import express from 'express';
import { diag, DiagConsoleLogger, DiagLogLevel } from '@opentelemetry/api';
import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';

// OpenTelemetry setup (minimal, OTLP exporter assumed via env vars)
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.ERROR);
const sdk = new NodeSDK({
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();

const app = express();
const PORT = Number(process.env.PORT) || 8080;
const SERVICE_A_URL = process.env.SERVICE_A_URL || 'http://service-a:8081';

app.get('/api', async (_req, res) => {
  try {
    const response = await fetch(`${SERVICE_A_URL}/process`);
    const data = await response.text();
    res.send(`frontend → ${data}`);
  } catch (err) {
    console.error(err);
    res.status(500).send('Error calling service-a');
  }
});

app.listen(PORT, () => {
  console.log(`frontend listening on ${PORT}`);
});
