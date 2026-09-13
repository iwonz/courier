## MODIFIED Requirements

### Requirement: Endpoint parsing

Courier SHALL classify local paths, SSH paths, exact `web://` and `webhook://` service endpoints, and HTTP(S) URLs into explicit endpoint kinds while preserving their source text and rejecting unknown URI schemes.

#### Scenario: SSH alias without user

- **WHEN** `source-server:/opt/data` is parsed
- **THEN** its kind is SSH, host is `source-server`, path is `/opt/data`, and user remains unset for SSH config resolution

#### Scenario: Windows path is local

- **WHEN** `C:\Users\me\data` or `C:/Users/me/data` is parsed
- **THEN** its kind is local and it is not interpreted as host `C` or as an URI scheme

#### Scenario: Remote Windows path

- **WHEN** `user@server:C:/Users/me/data` is parsed
- **THEN** its kind is SSH and its remote path is `C:/Users/me/data`

#### Scenario: Bracketed IPv6 SSH host

- **WHEN** `root@[2001:db8::1]:/opt/data` is parsed
- **THEN** its kind is SSH and the bracketed IPv6 host is retained

#### Scenario: Exact service endpoint

- **WHEN** `web://` or `webhook://` is parsed
- **THEN** it is classified as its Courier service kind without an authority or path

#### Scenario: Unknown scheme

- **WHEN** an endpoint begins with an unrecognized `name://` scheme
- **THEN** parsing fails instead of treating it as a local or SSH path

#### Scenario: Unicode and spaces

- **WHEN** a valid endpoint contains Unicode or spaces in its path
- **THEN** the path is retained byte-for-byte
