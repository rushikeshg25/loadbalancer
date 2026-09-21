# TCP load balancer project guide

> Generated: 2026-09-21 from commit `daba656`.

## What this is

This Go library distributes incoming TCP connections across configured backends and streams bytes in both directions ([implementation](../../lib.go#L82)). Applications supply a listener and context to `Serve`, while a small runnable example uses the older blocking convenience method ([API](../../lib.go#L14), [example](../../example/example.go#L27)). Round robin and least active connections share mutex-protected reservations, with another eligible backend attempted after a failed dial ([selection](../../lib.go#L54), [failover](../../lib.go#L96)).

## Run it

From the repository root, with Go 1.23.2 or a compatible newer toolchain ([go.mod](../../go.mod#L3)):

```bash
go build ./...
go test -race ./...
go run ./example
```

The example listens on `localhost:9090` and expects independently running TCP backends at `localhost:8080` and `localhost:8081`; it does not start those backends ([configuration](../../example/example.go#L8)). There are no external module dependencies or environment-variable configuration in the manifest and entry point ([go.mod](../../go.mod), [main](../../example/example.go#L7)). Stop the example process to exit; applications needing controlled cancellation should use `Serve` with their own context ([legacy wrapper](../../lib.go#L44)). The tests start ephemeral localhost sockets, so require permission to bind loopback ports ([test setup](../../lib_test.go#L13)). These are reproducible commands, not a claim of a new test run during guide generation.

## The five-file tour

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [example/example.go](../../example/example.go#L7) | See the configuration and legacy entry point. | [Startup](02-flow.md#startup) |
| 2 | [types.go](../../types.go#L5) | Learn the server list, strategy, and shared state. | [Data model](01-architecture.md#data-model) |
| 3 | [lib.go](../../lib.go#L14) | Trace listener ownership, reservation, dialing, and streams. | [Connection lifecycle](02-flow.md#connection-lifecycle) |
| 4 | [strategy.go](../../strategy.go#L3) | Understand the public selection-only helpers. | [Selection helpers](02-flow.md#selection-helpers) |
| 5 | [lib_test.go](../../lib_test.go#L12) | See the behavior demonstrated with real sockets. | [Build and verification](02-flow.md#build-and-verification) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, contracts, and failure boundaries.
2. [Flow](02-flow.md) — startup, proxying, shutdown, and verification.
3. [Structure](03-structure.md) — every source file and its callers.
4. [Tech stack](04-tech-stack.md) — Go, standard-library facilities, and tooling.
5. [Decisions](05-decisions.md) — evidence, inferred rationale, and gotchas.

## Open questions

- Is validation failure intended to leave listener ownership with the caller? Validation returns before deferred cleanup is registered, and `InitLoadbalancer` does not close its listener when `Serve` rejects configuration ([validation](../../lib.go#L15), [wrapper](../../lib.go#L44)).
- Should selection-only helpers return the counter after releasing their temporary reservation? They currently copy before release, so the returned `ServerRequests` includes that temporary increment ([strategy.go](../../strategy.go#L8)).
- What additional acceptance evidence is planned for invalid addresses/strategies, round robin, failed-dial retry, half-closes, and unexpected listener failure? The [v1 contract](../../V1.md#L9) asks for invalid-input and lifecycle-failure tests, but the three [current tests](../../lib_test.go#L12) cover only a subset.
- What production limits and operational observability are needed? The README excludes connection admission limits, and history explicitly records no production or long-running validation ([README](../../README.md#L17), [HISTORY](../../HISTORY.md#L58)).
