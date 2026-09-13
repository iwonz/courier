## 1. Worker data plane

- [x] 1.1 Extend private registration with bounded ephemeral runtime definitions and host rollback/stop hooks
- [x] 1.2 Serve one chi data host on the worker listener and dispatch opaque delivery tokens
- [x] 1.3 Open local and SSH-backed endpoints without persisting runtime secrets

## 2. Browser routes

- [x] 2.1 Implement protected metadata and filtered directory navigation without pre-authentication leakage
- [x] 2.2 Implement staged multipart upload, conflict handling, limits, cancellation, and cleanup
- [x] 2.3 Implement bounded file and deterministic whole-directory tar.gz downloads with aggregate rate enforcement
- [x] 2.4 Add password/Basic flows, CSRF, no-UI behavior, and foreground/background ownership

## 3. UI and contract

- [x] 3.1 Add the Lit data application using shared components, themes, and English/Russian catalogs
- [x] 3.2 Embed deterministic production assets and add asset-freshness verification
- [x] 3.3 Mark only browser routes and applicable command flags shipped in the canonical contract

## 4. Quality

- [x] 4.1 Cover HTTP isolation, upload/download safety, archives, limits, rates, stop, and IPC rollback at 100%
- [x] 4.2 Document browser URLs, authentication, navigation, and background lifecycle
- [x] 4.3 Run full verification, strict OpenSpec validation, archive the change, and commit once
