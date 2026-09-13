package policy

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strings"
)

type Admission struct {
	prefixes []netip.Prefix
}

func NewAdmission(entries []string) (*Admission, error) {
	prefixes := make([]netip.Prefix, 0, len(entries))
	seen := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		prefix, err := parseAdmissionEntry(entry)
		if err != nil {
			return nil, err
		}
		key := prefix.String()
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		prefixes = append(prefixes, prefix)
	}
	sort.Slice(prefixes, func(left, right int) bool { return prefixes[left].String() < prefixes[right].String() })
	return &Admission{prefixes: prefixes}, nil
}

func parseAdmissionEntry(value string) (netip.Prefix, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return netip.Prefix{}, errors.New("empty IP admission rule")
	}
	if address, err := netip.ParseAddr(value); err == nil {
		address = address.Unmap()
		return netip.PrefixFrom(address, address.BitLen()), nil
	}
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("invalid IP admission rule %q: %w", value, err)
	}
	if prefix.Addr().Is4In6() {
		if prefix.Bits() < 96 {
			return netip.Prefix{}, fmt.Errorf("mapped IPv6 admission prefix %q is ambiguous", value)
		}
		prefix = netip.PrefixFrom(prefix.Addr().Unmap(), prefix.Bits()-96)
	}
	return prefix.Masked(), nil
}

func (admission *Admission) Allows(address netip.Addr) bool {
	address = address.Unmap()
	if len(admission.prefixes) == 0 {
		return true
	}
	for _, prefix := range admission.prefixes {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}

func PeerIP(peer net.Addr) (netip.Addr, error) {
	if peer == nil {
		return netip.Addr{}, errors.New("connection peer is required")
	}
	if tcp, ok := peer.(*net.TCPAddr); ok {
		address, ok := netip.AddrFromSlice(tcp.IP)
		if !ok {
			return netip.Addr{}, errors.New("connection peer has an invalid IP")
		}
		return address.Unmap(), nil
	}
	return ParsePeerAddress(peer.String())
}

func ParsePeerAddress(remoteAddress string) (netip.Addr, error) {
	host, _, err := net.SplitHostPort(remoteAddress)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("invalid connection peer %q: %w", remoteAddress, err)
	}
	address, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("invalid connection peer IP %q: %w", host, err)
	}
	return address.Unmap(), nil
}
