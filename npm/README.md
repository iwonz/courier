# @iwonz/courier

This package installs the verified standalone Courier binary for macOS, Linux, or Windows on amd64 or arm64.

```sh
npx @iwonz/courier --help
```

The postinstall script downloads the matching asset from [GitHub Releases](https://github.com/iwonz/courier/releases), verifies its SHA-256 entry from `checksums.txt`, and safely extracts only the Courier executable. Installation with `--ignore-scripts` is not supported.

See the [main repository](https://github.com/iwonz/courier) for transfer and security documentation.
