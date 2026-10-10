## 1. Specification and IPC behavior

- [x] 1.1 Specify idle-safe foreground subscriptions and compiled-runtime timeout-boundary acceptance
- [x] 1.2 Add distinct streaming IPC deadline configuration while retaining the bounded RPC default
- [x] 1.3 Route progress subscription reads through streaming configuration with cancellation and error coverage

## 2. Regression acceptance

- [x] 2.1 Add a compiled foreground hosted-route idle regression beyond five seconds followed by clean external stop
- [x] 2.2 Run focused IPC, worker, application, and browser runtime tests

## 3. Verification and release

- [x] 3.1 Run strict OpenSpec validation and the complete `make verify` gate
- [x] 3.2 Ship the archived change as the next patch and verify release surfaces, Pages, deleted feature branch, and clean synchronized main
