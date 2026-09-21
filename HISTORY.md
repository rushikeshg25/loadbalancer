# TCP load balancer history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2024-12-23T13:15:32+05:30: Initial commit

- **What happened:** The repository records `Initial commit`.
- **Evidence:** [commit 83eb67d0b1](https://github.com/rushikeshg25/loadbalancer/commit/83eb67d0b12fd8f3e69b94b0fa580219b354970d).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:21:37+05:30: docs: define loadbalancer v1 contract

- **What happened:** The repository records `docs: define loadbalancer v1 contract`.
- **Evidence:** [commit da8c5784e4](https://github.com/rushikeshg25/loadbalancer/commit/da8c5784e4feaff29b868e427b7dc0ef3c6dcb48).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:26:56+05:30: feat: stream duplex TCP with connection reservations and shutdown

- **What happened:** The repository records `feat: stream duplex TCP with connection reservations and shutdown`.
- **Evidence:** [commit 4d56ebfd27](https://github.com/rushikeshg25/loadbalancer/commit/4d56ebfd277942d9f2f9f0d2aa7d09f237b8349c).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:26:56+05:30: test: verify large streaming transfers balancing and cancellation

- **What happened:** The repository records `test: verify large streaming transfers balancing and cancellation`.
- **Evidence:** [commit 34f399e558](https://github.com/rushikeshg25/loadbalancer/commit/34f399e558b09bb9db5e9f715b972ec2c9fdecc1).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:32:07+05:30: fix: cancel active streams on listener failure and preserve strategy APIs

- **What happened:** The repository records `fix: cancel active streams on listener failure and preserve strategy APIs`.
- **Evidence:** [commit d055365224](https://github.com/rushikeshg25/loadbalancer/commit/d0553652243454fe32346cb2dd62c60eafe3d54c).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:33:04+05:30: docs: document loadbalancer v1 usage and limitations

- **What happened:** The repository records `docs: document loadbalancer v1 usage and limitations`.
- **Evidence:** [commit bbfc290a65](https://github.com/rushikeshg25/loadbalancer/commit/bbfc290a650fd0a398887d1bc02706205d5398b5).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed `go test -race ./...` with real localhost sockets; [lib_test.go](lib_test.go) checks a 140 KB duplex transfer, least-connection reservations and active-stream cancellation.

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.

## 2026-09-21: Maintainer project guide

Added the six-file [project guide](docs/project-guide/README.md), tracing architecture, runtime flows, source structure, dependencies and decisions against the v1 code. Relative paths, source/heading anchors and Mermaid syntax were checked. The guide distinguishes observed behavior from inferred rationale and records remaining limitations.
