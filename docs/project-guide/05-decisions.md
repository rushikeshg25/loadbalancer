# Decisions

The behavior below is confirmed by source. Motivations are marked inferred unless an explicit comment, contract, or README supports them.

## Make the listener and context caller supplied

- **What:** `Serve` takes a `net.Listener` and context; the legacy initializer creates its own listener with a background context.
- **Evidence:** [Serve signature and ownership comment](../../lib.go#L13), [legacy wrapper](../../lib.go#L44), [README recommendation](../../README.md#L17).
- **Why, apparently:** Explicit cancellation and error handling are intended for new applications; this rationale is confirmed by the README.
- **Tradeoff:** Embedders control bind configuration and service lifetime, but must handle listener cleanup if validation fails before ownership cleanup is registered.
- **Confidence:** Confirmed API recommendation; ownership edge case derived from source ordering.

## Reserve before dialing

- **What:** Selection increments the chosen server's count under the shared mutex before attempting a network connection, and releases on dial failure or stream completion.
- **Evidence:** [reserveWith](../../lib.go#L54), [dial failure and deferred release](../../lib.go#L96), [README](../../README.md#L3).
- **Why, apparently:** Simultaneous connections should observe active load; confirmed by the README. Counting in-progress dials as load follows from the implementation.
- **Tradeoff:** Avoids multiple concurrent selectors seeing the same idle count, but serializes pool selection through one mutex and includes slow dials in the load estimate.
- **Confidence:** Confirmed mechanism and stated concurrency purpose.

## Share a rotating scan across strategies

- **What:** Both strategies start scanning at `TotalRequests`; least-connections keeps the smallest count, while other strategy values take the first eligible entry. A successful selection moves the cursor past the selected index.
- **Evidence:** [scan and tie handling](../../lib.go#L59), [cursor update](../../lib.go#L73), [Serve strategy validation](../../lib.go#L23).
- **Why, apparently:** Inferred: one selection implementation keeps eligibility rules consistent and rotates least-connections ties.
- **Tradeoff:** Compact shared logic, but `TotalRequests` is a misleading name for a cursor. An empty strategy behaves as round robin; standalone selector helpers do not validate unknown strategies.
- **Confidence:** Behavior confirmed; rationale inferred.

## Stream bytes and preserve TCP half-closes

- **What:** Two concurrent `io.Copy` operations stream the connection; clean EOF calls destination `CloseWrite` if available, while errors close both sockets.
- **Evidence:** [copy lifecycle](../../lib.go#L105), [optional CloseWrite interface](../../lib.go#L127), [v1 contract](../../V1.md#L5), [README](../../README.md#L15).
- **Why, apparently:** Full-duplex streaming and peers finishing responses after a half-close are explicit requirements/documentation.
- **Tradeoff:** Handles payloads beyond a fixed request buffer and supports opaque TCP protocols, but does not inspect requests, retry partial application transactions, or limit stream duration. A custom connection wrapper lacking `CloseWrite` cannot propagate a half-close through this helper.
- **Confidence:** Confirmed behavior and stated streaming intent.

## Keep failed-dial retry local to one client

- **What:** Each client tracks attempted backend indices. A failed dial releases its reservation and retries selection without mutating `IsHealthy`.
- **Evidence:** [exclusion map and dial loop](../../lib.go#L86), [static health scope](../../README.md#L17).
- **Why, apparently:** Inferred: allow immediate failover while keeping active health probes and dynamic pool management outside v1.
- **Tradeoff:** A single failed backend does not immediately abort a client, but every new client can attempt that same unavailable server, and sequential dial timeouts accumulate.
- **Confidence:** Mechanism and static eligibility confirmed; rationale inferred.

## Preserve selection-only compatibility methods

- **What:** `NextServer`, `RoundRobin`, and `LeastConn` use shared reservation logic, copy the backend, then release immediately.
- **Evidence:** [strategy.go](../../strategy.go#L3), [history entry](../../HISTORY.md#L36).
- **Why, apparently:** The history identifies preservation of strategy APIs; using the same pool logic is an inferred way to retain consistent eligibility and ordering.
- **Tradeoff:** Existing callers can still select a backend, but these methods do not reserve a connection for the caller's later use. Their returned counter snapshot includes the temporary reservation and calling them advances the shared cursor.
- **Confidence:** API preservation supported by recorded commit title; implementation consequences confirmed by source.

## Gotchas

- **Listener rejection and legacy errors:** Validation occurs before `defer listener.Close()`. The legacy wrapper neither closes after a rejected `Serve` call nor reports listen/service errors ([lib.go:15](../../lib.go#L15), [lib.go:44](../../lib.go#L44)).
- **Health defaults to false:** A syntactically valid pool with every `IsHealthy` unset passes startup checks but cannot select a backend; each accepted client is closed ([types.go:8](../../types.go#L8), [filter](../../lib.go#L61), [proxy exit](../../lib.go#L89)).
- **Counters are implementation state:** Initialize counters/cursor to their zero defaults and do not read or mutate them concurrently with service. They are exported but synchronized only inside library operations, and arbitrary negative cursor values can produce invalid slice indices ([README contract](../../README.md#L17), [index calculation](../../lib.go#L59)).
- **Cancellation terminates active streams:** The documented shutdown waits for goroutines after closing sockets, rather than letting existing streams finish normally. The example has no signal-to-context integration ([cleanup](../../lib.go#L31), [client callback](../../lib.go#L84), [example call](../../example/example.go#L27)).
- **No service error for per-client failures:** Dial/copy errors cause retry or closure locally; only validation and accept errors reach `Serve`'s caller ([dial handling](../../lib.go#L97), [copy handling](../../lib.go#L108), [accept handling](../../lib.go#L34)).
- **Acceptance is broader than direct test coverage:** Current tests demonstrate a 140 KB echo, cancellation, least-connections reservations/eligibility, and empty configuration. They do not directly exercise failed-dial failover, half-close completion, explicit round-robin order, unexpected accept failure, or other validation branches ([tests](../../lib_test.go#L12), [acceptance](../../V1.md#L9)).

## Conventions

- Public configuration/types live in [types.go](../../types.go); lifecycle/network operations live in [lib.go](../../lib.go); public strategy wrappers live in [strategy.go](../../strategy.go).
- Keep shared selection and count changes inside the mutex-protected reservation/release methods ([lib.go:54](../../lib.go#L54), [lib.go:77](../../lib.go#L77)).
- Tests use package-internal helpers and real ephemeral localhost listeners, with explicit deadlines around socket verification ([lib_test.go:1](../../lib_test.go#L1), [lib_test.go:13](../../lib_test.go#L13), [lib_test.go:40](../../lib_test.go#L40)).
