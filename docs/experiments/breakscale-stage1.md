# Breakscale — Stage 1 Baseline

## Architecture

Client → Go API Server → PostgreSQL

No caching, queues, replicas, or horizontal scaling.

## Purpose

Use Breakscale to build an initial intuition for how a simple synchronous architecture behaves as offered load approaches a dependency's capacity.

## Experiment

At 200 requests/second:

- Offered load: 200 req/s
- Throughput: 191 req/s
- API utilisation: 7.8%
- Database utilisation: 97%
- Database lock wait: 2.1 ms
- System P99: 250 ms
- Failures: 0%

## Observation

The database approached saturation while the API server remained lightly utilised.

This suggests the database can become the limiting resource before API compute capacity is exhausted.

At higher load, intermittent request failures were observed, with reported failure levels including 11%, 9.4%, and a worst observed value of 16%. Requests subsequently recovered.

## Interpretation

Increasing offered load beyond the system's sustainable capacity does not necessarily increase throughput. Work can accumulate around the saturated dependency, increasing latency and eventually causing failures.

## Limitation

These are Breakscale simulation results, not measurements from the Flashbulb implementation.

The simulation is being used as architectural intuition and hypothesis generation. Actual Flashbulb performance results will be measured separately using the real Go application and PostgreSQL deployment.

## Initial Hypothesis

For the Stage 1 Flashbulb architecture, database capacity may become a bottleneck before application-server compute capacity.

This hypothesis must be tested experimentally rather than assumed.
