#!/bin/sh

set -eu

repository=${COURIER_REPOSITORY:-iwonz/courier}
version=${COURIER_VERSION:-latest}
install_dir=${COURIER_INSTALL_DIR:-"${HOME}/.local/bin"}
release_api=${COURIER_RELEASE_API_URL:-"https://api.github.com/repos/${repository}/releases/latest"}
release_base=${COURIER_RELEASE_BASE_URL:-"https://github.com/${repository}/releases/download"}

fail() {
  printf '%s\n' "courier installer: $*" >&2
  exit 1
}

download() {
  source_url=$1
  destination_file=$2
  if command -v curl >/dev/null 2>&1; then
    curl --fail --silent --show-error --location "$source_url" --output "$destination_file"
  elif command -v wget >/dev/null 2>&1; then
    wget --quiet --output-document="$destination_file" "$source_url"
  else
    fail "curl or wget is required to download the release"
  fi
}

temporary_dir=$(mktemp -d "${TMPDIR:-/tmp}/courier-install.XXXXXX") || fail "cannot create a temporary directory"
chmod 700 "$temporary_dir"
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM

if [ "$version" = "latest" ]; then
  download "$release_api" "$temporary_dir/release.json"
  version=$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\(v[^"]*\)".*/\1/p' "$temporary_dir/release.json" | sed -n '1p')
fi

case "$version" in
  v*) tag=$version ;;
  *) tag="v$version" ;;
esac
printf '%s\n' "$tag" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$' || fail "release version must be vMAJOR.MINOR.PATCH"
version_number=${tag#v}

case $(uname -s) in
  Darwin) target_os=darwin ;;
  Linux) target_os=linux ;;
  *) fail "unsupported operating system: $(uname -s)" ;;
esac

case $(uname -m) in
  x86_64|amd64) target_arch=amd64 ;;
  arm64|aarch64) target_arch=arm64 ;;
  *) fail "unsupported architecture: $(uname -m)" ;;
esac

asset="courier_${version_number}_${target_os}_${target_arch}"
asset_url="${release_base}/${tag}/${asset}"
checksum_url="${release_base}/${tag}/checksums.txt"
download "$asset_url" "$temporary_dir/$asset"
download "$checksum_url" "$temporary_dir/checksums.txt"

expected=$(awk -v name="$asset" '$2 == name || $2 == "*" name { print tolower($1); exit }' "$temporary_dir/checksums.txt")
[ -n "$expected" ] || fail "checksums.txt has no entry for $asset"

if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$temporary_dir/$asset" | awk '{ print tolower($1) }')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$temporary_dir/$asset" | awk '{ print tolower($1) }')
elif command -v openssl >/dev/null 2>&1; then
  actual=$(openssl dgst -sha256 "$temporary_dir/$asset" | awk '{ print tolower($NF) }')
else
  fail "sha256sum, shasum, or openssl is required to verify the release"
fi
[ "$actual" = "$expected" ] || fail "checksum mismatch for $asset"

umask 077
mkdir -p "$install_dir"
staged="$install_dir/.courier-install.$$"
trap 'rm -f "$staged"; rm -rf "$temporary_dir"' EXIT HUP INT TERM
cp "$temporary_dir/$asset" "$staged"
chmod 755 "$staged"
mv -f "$staged" "$install_dir/courier"
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM

printf '%s\n' "Installed courier $version_number to $install_dir/courier"
case ":${PATH}:" in
  *":${install_dir}:"*) ;;
  *) printf '%s\n' "Add $install_dir to PATH to invoke courier directly." ;;
esac
