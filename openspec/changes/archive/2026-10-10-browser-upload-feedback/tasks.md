## 1. Observable browser upload

- [x] 1.1 Add bounded typed API failures carrying HTTP status and guarded server reasons
- [x] 1.2 Add localized pending, committed, and rejected upload presentation with accessible live semantics
- [x] 1.3 Disable overlapping selections and reset the input after both success and failure

## 2. Regression coverage

- [x] 2.1 Cover API decoding and component pending/success/failure/retry behavior
- [x] 2.2 Extend compiled browser acceptance to assert visible success, visible collision failure, same-file retry, and exact committed bytes
- [x] 2.3 Strengthen compiled incoming-webhook acceptance for guarded success/collision outcomes and preserved bytes
- [x] 2.4 Run focused data-app, Go web-delivery, and real browser runtime tests

## 3. Verification and release

- [x] 3.1 Run strict OpenSpec validation and the complete `make verify` gate
- [x] 3.2 Ship the archived change as the next patch and verify GitHub Release, npm, Scoop, Homebrew, Pages, deleted feature branch, and clean synchronized main
