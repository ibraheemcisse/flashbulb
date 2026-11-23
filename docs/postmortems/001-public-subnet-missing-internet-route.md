# Postmortem 001 — Missing Internet Route

**Date:** 2026-11-23  
**Stage:** Stage 1 — EC2 deployment  
**Status:** Resolved

## Summary

The Stage 1 EC2 instance could not be reached over SSH.

The instance was deployed in `us-east-1a` within the `evershop-dev-public-us-east-1a` subnet. The VPC had an attached Internet Gateway, but the route table associated with the subnet did not contain a default route to the Internet Gateway.

## Symptoms

SSH connections failed during connection establishment:

```text
kex_exchange_identification: read: Connection reset
Connection reset by <instance-ip> port 22


EC2 system and instance status checks were passing.
Investigation
The subnet's route table was:
rtb-0dcc9520927da7355


Initially, it contained only:
10.0.0.0/16 → local


The VPC had an attached Internet Gateway:
igw-0df56fc0050e52b45


However, there was no default internet route.
Root Cause
The subnet's route table was missing:
0.0.0.0/0 → igw-0df56fc0050e52b45


Therefore, traffic from the subnet destined for the public internet had no route through the Internet Gateway.
Resolution
Added the following route to rtb-0dcc9520927da7355:
Destination: 0.0.0.0/0
Target:      igw-0df56fc0050e52b45


Connectivity was restored after the route was added.
Impact
The missing route blocked external connectivity to the EC2 instance and prevented SSH access during the initial deployment.
Lesson
The subnet was named as a public subnet, but its effective routing configuration did not provide internet access. The route table configuration was the actual determining factor.
