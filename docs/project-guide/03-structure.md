# Structure

## What lives where

This small Go module has ten tracked files at the inspected commit. The root package owns the complete TCP lifecycle and shared backend state; `example/` is the only executable package. There are no nested service or persistence packages in this inventory.

```text
loadbalancer/
├── go.mod                 module path and Go directive
├── types.go               public configuration and state types
├── lib.go                 listener, reservation, dialing, and streams
├── strategy.go            public selection-only helpers
├── lib_test.go            socket and pool behavior tests
├── example/
│   └── example.go         legacy blocking sample
├── README.md              public usage and limitations
├── V1.md                  acceptance contract
├── HISTORY.md             delivery history and verification record
├── .gitignore             local artifact exclusions
└── docs/project-guide/    this guide
```

## Root package

All runtime root files share package `loadbalancer`. The tests use that same package and can exercise unexported reservation functions ([test package](../../lib_test.go#L1)).

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [types.go](../../types.go#L5) | Defines backend configuration, shared balancer state, and strategy names. | `BackendServer`, `Loadbalancer`, `StrategyType`, `ROUNDROBIN`, `LEASTCONN` | All library code, example, embedders, tests |
| [lib.go](../../lib.go#L14) | Validates/owns a listener after validation, accepts clients, reserves backends, retries dials, proxies streams, and handles shutdown. | `Serve`, `InitLoadbalancer`; internal `reserve`, `reserveWith`, `release`, `proxy`, `closeWrite` | Example/embedders invoke exported methods; strategy helpers use reservations; tests call service and reservation methods |
| [strategy.go](../../strategy.go#L3) | Selects a server value while releasing the temporary reservation before return. | `NextServer`, `RoundRobin`, `LeastConn`; internal `selectServer` | External library consumers; the proxy uses `reserve` directly instead |
| [lib_test.go](../../lib_test.go#L12) | Exercises 140 KB streaming, context shutdown, least-connections/eligibility, and empty configuration. | `TestStreamAndShutdown`, `TestStrategies`, `TestEmptyConfiguration` | Go test runner |

## Example

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [example/example.go](../../example/example.go#L7) | Configures two localhost backends and starts a round-robin proxy on port 9090. | Executable `main` | `go run ./example` or built sample executable |

## Repository documents and configuration

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [go.mod](../../go.mod#L1) | Declares module import path and Go 1.23.2. | Module `github.com/rushikeshg25/loadbalancer` | Go build/test tools and downstream consumers |
| [README.md](../../README.md#L3) | Documents `Serve`, streaming, strategy choices, cancellation, and unsupported features. | Public usage contract | Library users and maintainers |
| [V1.md](../../V1.md#L3) | Defines v1 scope and behavioral-test acceptance expectations. | Contract and acceptance criteria | Delivery reviewers |
| [HISTORY.md](../../HISTORY.md#L10) | Records baseline and delivery commits plus historical verification claims. | Timeline and evidence links | Maintainers tracing changes |

## Excluded

The tracked [.gitignore](../../.gitignore) is omitted from the design tables because it is artifact-filter boilerplate. Git internals are not source. An existing untracked `logs/` directory was present during inspection and is not part of the tracked source inventory or this guide's edits. There are no vendored dependencies, lockfiles, CI files, deployment manifests, or generated source in the tracked inventory. The six generated guide files describe the source and are not additional runtime components.
