# Tasks

## 1. Specify the fallback boundary

- [x] 1.1 Require a typed SFTP capability gap and explicit consent
- [x] 1.2 Specify verified temporary staging and cleanup
- [x] 1.3 Specify Windows agent and update handoff behavior

## 2. Implement helper acquisition and transport

- [x] 2.1 Preserve the authenticated SSH connection for the fallback callback
- [x] 2.2 Acquire an exact-platform helper and verify release checksums
- [x] 2.3 Upload, remotely verify, and run the embedded SFTP server
- [x] 2.4 Guarantee local and remote cleanup on every terminal path

## 3. Integrate consent and hidden runtime entry points

- [x] 3.1 Add serialized interactive confirmation with refusal by default
- [x] 3.2 Serve SFTP only through an internal process mode
- [x] 3.3 Keep native SFTP free of helper prompts and deployment

## 4. Harden Windows behavior

- [x] 4.1 Dial the Windows OpenSSH agent named pipe natively
- [x] 4.2 Implement staged Windows update handoff and cleanup

## 5. Verify and document

- [x] 5.1 Add deterministic unit and integration coverage with cleanup assertions
- [x] 5.2 Retain exact 100% first-party Go statement coverage and race safety
- [x] 5.3 Cross-build all supported OS/architecture targets
- [x] 5.4 Update English user, installation, and security documentation
- [x] 5.5 Pass strict OpenSpec validation and the complete local release dry run
