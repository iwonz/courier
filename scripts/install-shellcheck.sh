#!/bin/sh

set -eu

version=${SHELLCHECK_VERSION:-v0.11.0}
destination=${SHELLCHECK_INSTALL_PATH:-.cache/tools/shellcheck}

if [ -x "$destination" ] && "$destination" --version 2>/dev/null | grep -q "version: ${version#v}"; then
  exit 0
fi

case $(uname -s) in
  Darwin) target_os=darwin ;;
  Linux) target_os=linux ;;
  *) printf '%s\n' "ShellCheck bootstrap supports macOS and Linux" >&2; exit 1 ;;
esac
case $(uname -m) in
  x86_64|amd64) target_arch=x86_64 ;;
  arm64|aarch64) target_arch=aarch64 ;;
  *) printf '%s\n' "Unsupported ShellCheck host architecture: $(uname -m)" >&2; exit 1 ;;
esac

case "$version:$target_os:$target_arch" in
  v0.11.0:darwin:aarch64) expected=339b930feb1ea764467013cc1f72d09cd6b869ebf1013296ba9055ab2ffbd26f ;;
  v0.11.0:darwin:x86_64) expected=c2c15e08df0e8fbc374c335b230a7ee958c313fa5714817a59aa59f1aa594f51 ;;
  v0.11.0:linux:aarch64) expected=68a8133197a50beb8803f8d42f9908d1af1c5540d4bb05fdfca8c1fa47decefc ;;
  v0.11.0:linux:x86_64) expected=b7af85e41cc99489dcc21d66c6d5f3685138f06d34651e6d34b42ec6d54fe6f6 ;;
  *) printf '%s\n' "No trusted ShellCheck checksum for $version $target_os/$target_arch" >&2; exit 1 ;;
esac

temporary_dir=$(mktemp -d "${TMPDIR:-/tmp}/courier-shellcheck.XXXXXX")
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM
archive="shellcheck-${version}.${target_os}.${target_arch}.tar.gz"
base="https://github.com/koalaman/shellcheck/releases/download/${version}"

if command -v curl >/dev/null 2>&1; then
  curl --fail --silent --show-error --location "$base/$archive" --output "$temporary_dir/$archive"
elif command -v wget >/dev/null 2>&1; then
  wget --quiet --output-document="$temporary_dir/$archive" "$base/$archive"
else
  printf '%s\n' "curl or wget is required to bootstrap ShellCheck" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$temporary_dir/$archive" | awk '{ print tolower($1) }')
else
  actual=$(shasum -a 256 "$temporary_dir/$archive" | awk '{ print tolower($1) }')
fi
[ "$actual" = "$expected" ] || { printf '%s\n' "ShellCheck checksum mismatch" >&2; exit 1; }

tar -xzf "$temporary_dir/$archive" -C "$temporary_dir" "shellcheck-${version}/shellcheck"
mkdir -p "$(dirname "$destination")"
chmod 755 "$temporary_dir/shellcheck-${version}/shellcheck"
mv "$temporary_dir/shellcheck-${version}/shellcheck" "$destination"
