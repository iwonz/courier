#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

command -v docker >/dev/null 2>&1 || {
  printf '%s\n' "Docker is required for Linux package acceptance" >&2
  exit 1
}
[ -f dist/metadata.json ] || {
  printf '%s\n' "dist/metadata.json is missing; run the GoReleaser snapshot first" >&2
  exit 1
}

version=$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' dist/metadata.json)
[ -n "$version" ] || {
  printf '%s\n' "cannot read snapshot version" >&2
  exit 1
}

case $(docker info --format '{{.Architecture}}') in
  x86_64|amd64) architecture=amd64; platform=linux/amd64 ;;
  aarch64|arm64) architecture=arm64; platform=linux/arm64 ;;
  *) printf '%s\n' "Docker server architecture is unsupported" >&2; exit 1 ;;
esac
arch_architecture=$architecture
arch_platform=$platform
if [ "$architecture" = arm64 ]; then
  # The official Arch image is amd64-only. CI exercises it on its native amd64
  # runner; ARM hosts still verify the arm64 Arch package with Manjaro/pacman.
  arch_architecture=amd64
  arch_platform=linux/amd64
fi

if command -v uuidgen >/dev/null 2>&1; then
  operation_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
else
  random_hex=$(od -An -N16 -tx1 /dev/urandom | tr -d ' \n')
  operation_id=$(printf '%s\n' "$random_hex" | awk '{ print substr($0,1,8) "-" substr($0,9,4) "-" substr($0,13,4) "-" substr($0,17,4) "-" substr($0,21,12) }')
fi

label_key=com.iwonz.courier.acceptance
label="$label_key=$operation_id"
temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/courier-packages.XXXXXX")

owned_ids() {
  docker "$1" ls -q --filter "label=$label"
}

cleanup() {
  cleanup_status=0
  if ! container_ids=$(docker ps -aq --filter "label=$label" 2>/dev/null); then
    container_ids=
    cleanup_status=1
  fi
  if [ -n "$container_ids" ]; then
    docker rm -f $container_ids >/dev/null 2>&1 || cleanup_status=1
  fi
  if ! network_ids=$(owned_ids network 2>/dev/null); then
    network_ids=
    cleanup_status=1
  fi
  if [ -n "$network_ids" ]; then
    docker network rm $network_ids >/dev/null 2>&1 || cleanup_status=1
  fi
  if ! volume_ids=$(owned_ids volume 2>/dev/null); then
    volume_ids=
    cleanup_status=1
  fi
  if [ -n "$volume_ids" ]; then
    docker volume rm -f $volume_ids >/dev/null 2>&1 || cleanup_status=1
  fi
  rm -rf "$temporary_directory" || cleanup_status=1

  if ! remaining_containers=$(docker ps -aq --filter "label=$label" 2>/dev/null); then
    cleanup_status=1
  elif [ -n "$remaining_containers" ]; then
    cleanup_status=1
  fi
  if ! remaining_networks=$(owned_ids network 2>/dev/null); then
    cleanup_status=1
  elif [ -n "$remaining_networks" ]; then
    cleanup_status=1
  fi
  if ! remaining_volumes=$(owned_ids volume 2>/dev/null); then
    cleanup_status=1
  elif [ -n "$remaining_volumes" ]; then
    cleanup_status=1
  fi
  [ ! -e "$temporary_directory" ] || cleanup_status=1
  return "$cleanup_status"
}

finish() {
  status=$?
  trap - EXIT HUP INT TERM
  if ! cleanup; then
    printf '%s\n' "Courier-owned Docker resources remain for $operation_id" >&2
    status=1
  fi
  exit "$status"
}
trap finish EXIT
trap 'exit 130' HUP INT TERM

deb_package=$(pwd)/dist/courier_${version}_linux_${architecture}.deb
rpm_package=$(pwd)/dist/courier_${version}_linux_${architecture}.rpm
apk_package=$(pwd)/dist/courier_${version}_linux_${architecture}.apk
arch_package=$(pwd)/dist/courier_${version}_linux_${architecture}.pkg.tar.zst
arch_compatible_package=$(pwd)/dist/courier_${version}_linux_${arch_architecture}.pkg.tar.zst
for package in "$deb_package" "$rpm_package" "$apk_package" "$arch_package" "$arch_compatible_package"; do
  [ -f "$package" ] || {
    printf '%s\n' "missing Linux package: $package" >&2
    exit 1
  }
done

run_install() {
  distribution=$1
  image=$2
  package=$3
  container_platform=$4
  install_command=$5
  container_package=/tmp/$(basename "$package")
  printf '%s\n' "Verifying $distribution package in $image"
  docker run --rm \
    --platform "$container_platform" \
    --label "$label" \
    --name "courier-acceptance-$operation_id-$distribution" \
    --env "COURIER_PACKAGE=$container_package" \
    --mount "type=bind,source=$package,target=$container_package,readonly" \
    "$image" sh -ec "$install_command; courier version"
}

run_install ubuntu ubuntu:24.04 "$deb_package" "$platform" 'dpkg -i "$COURIER_PACKAGE" >/dev/null'
run_install debian debian:bookworm-slim "$deb_package" "$platform" 'dpkg -i "$COURIER_PACKAGE" >/dev/null'
if [ "$architecture" = amd64 ]; then
  run_install arch archlinux:base "$arch_compatible_package" "$arch_platform" 'pacman --noconfirm -U "$COURIER_PACKAGE" >/dev/null'
else
  printf '%s\n' "Skipping the amd64-only official Arch image on an ARM Docker server; CI runs it natively"
fi
run_install manjaro manjarolinux/base:20260322 "$arch_package" "$platform" 'pacman --noconfirm -U "$COURIER_PACKAGE" >/dev/null'
run_install fedora fedora:43 "$rpm_package" "$platform" 'rpm -Uvh --replacepkgs "$COURIER_PACKAGE" >/dev/null'
run_install rhel registry.access.redhat.com/ubi9/ubi-minimal:9.6 "$rpm_package" "$platform" 'rpm -Uvh --replacepkgs "$COURIER_PACKAGE" >/dev/null'
run_install alpine alpine:3.22 "$apk_package" "$platform" 'apk add --allow-untrusted --no-network "$COURIER_PACKAGE" >/dev/null'

printf '%s\n' "Verified native Linux packages on $architecture"
