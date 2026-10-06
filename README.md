# shortlink-platform

A Go URL shortener delivered through a full CI/CD and observability platform on local Kubernetes.

## Scope

A small URL shortener written in Go, used as the workload for a complete
delivery platform that runs on a local kind cluster at no cloud cost.

**In scope**
- A Go service that creates short links and redirects to the original URL,
  with a PostgreSQL database for storage
- CI on GitHub Actions: lint, test, vulnerability scan, image build and push
- GitOps delivery with Helm and Argo CD
- Metrics and alerting with Prometheus and Grafana, logs with ELK
- One command to bring everything up and one to tear it down
- Documentation: architecture diagram, runbook, failure and load experiments

**Out of scope**
- Authentication, user accounts, a web frontend
- Cloud-specific services (the goal is to run anywhere with Docker)
- Multi-cluster, multi-region, or production-grade high availability of the database

**Goal**
Anyone can clone the repository, run one command, and see the full path from
a code change to a running, monitored service.

## Status

Work in progress. See the commit history for progress.


<img width="1884" height="858" alt="image" src="https://github.com/user-attachments/assets/57e518a2-ef55-4161-8b7a-d5ec053fe252" />
