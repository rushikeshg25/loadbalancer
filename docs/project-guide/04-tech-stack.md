# Tech stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `go 1.23.2` minimum language/toolchain directive; a newer compatible installed toolchain may be used | [go.mod:3](../../go.mod#L3) |
| Module | `github.com/rushikeshg25/loadbalancer`; no module release version declared here | [go.mod:1](../../go.mod#L1) |

## Frameworks and major libraries

There is no web framework and no third-party dependency requirement in [go.mod](../../go.mod). The following are Go standard-library packages, so their implementation version follows the selected Go toolchain rather than independent dependency pins.

| Library | Version | Used for | Used in |
| --- | --- | --- | --- |
| `net` | Go standard library | TCP listeners, accepted connections, address formatting, context-aware dialing with timeout | [lib.go:33](../../lib.go#L33), [lib.go:45](../../lib.go#L45), [lib.go:96](../../lib.go#L96) |
| `io` | Go standard library | Concurrent byte-stream copying and test reads | [lib.go:107](../../lib.go#L107), [lib_test.go:44](../../lib_test.go#L44) |
| `context` | Go standard library | Caller cancellation, listener/socket callbacks, derived service context | [lib.go:27](../../lib.go#L27), [lib.go:84](../../lib.go#L84) |
| `sync` | Go standard library | Mutex protecting shared reservations; wait group joining accepted connections | [types.go:13](../../types.go#L13), [lib.go:30](../../lib.go#L30) |
| `time` | Go standard library | Three-second dial timeout and test deadlines | [lib.go:96](../../lib.go#L96), [lib_test.go:40](../../lib_test.go#L40) |
| `testing`, `bytes` | Go standard library | Behavioral assertions and deterministic 140 KB payload | [lib_test.go:3](../../lib_test.go#L3), [lib_test.go:41](../../lib_test.go#L41) |

## Data and infrastructure

| Service | Role | Configured at |
| --- | --- | --- |
| Frontend TCP listener | Supplied by caller to `Serve`; legacy wrapper derives its bind address from `Host`/`Port` | [lib.go:14](../../lib.go#L14), [lib.go:45](../../lib.go#L45) |
| Backend TCP servers | Application-provided addresses and static health eligibility; sample expects two independently running localhost services | [types.go:5](../../types.go#L5), [example/example.go:8](../../example/example.go#L8) |
| In-memory counters | Reservation count per backend and rotation cursor per balancer; no datastore dependency | [lib.go:54](../../lib.go#L54), [lib.go:77](../../lib.go#L77) |

No database, cache, message queue, TLS termination, or health-probing library is present in the [module](../../go.mod) or [runtime path](../../lib.go). The README explicitly excludes TLS termination and dynamic discovery ([scope](../../README.md#L17)).

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| `go build ./...` | Compile library and example using the Go module | [go.mod](../../go.mod), [example import](../../example/example.go#L4) |
| `go test -race ./...` | Run socket/state tests with Go's race detector | [README command](../../README.md#L19), [tests](../../lib_test.go#L12) |
| `go run ./example` | Start the legacy example process | [README command](../../README.md#L19), [entry point](../../example/example.go#L7) |

## Notes

- There is no separate toolchain pin, CI workflow, Makefile, linter configuration, or deployment recipe among the [tracked files](03-structure.md). Do not infer deployment automation from the historical test claim in [HISTORY.md](../../HISTORY.md#L50).
- The example's backend addresses are compiled Go literals, while library callers construct configuration themselves; no environment loader is present ([example](../../example/example.go#L8)).
- `context.AfterFunc` and `DialContext` are central lifecycle mechanisms, not optional framework hooks ([lib.go](../../lib.go#L27)).
