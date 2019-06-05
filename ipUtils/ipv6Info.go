package ipUtils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Ipv6Addr [9]int

func ParseIpv6(str string) (addr Ipv6Addr, err error) {
	if str == "" {
		return addr, fmt.Errorf("input string is empty")
	}

	// There's no way I'm writing a regex for this one, so I'm going to "cheat"
	if _, _, err := net.ParseCIDR(str); err != nil && net.ParseIP(str) == nil {
		return addr, fmt.Errorf("invalid address by net pkg")
	}

	split := strings.Split(str, "/")
	if len(split) == 2 && split[1] != "" {
		prefix, err := strconv.Atoi(split[1])
		if err != nil {
			// I could continue anyway, but I don't want to deal with channels or something, so I'll just quit
			// This should be prevented by the checks above, anyway
			return addr, fmt.Errorf("prefix is NaN")
		}

		addr[8] = prefix
	} else {
		addr[8] = -1
	}

	spl := strings.Split(split[0], ":")

	i := 0
	ellipsis := false
	for si, hextet := range spl { // Don't you quibble with me!
		if hextet == "" {
			// Position of "ellipsis" (this is what Go's net package calls the :: as far as I can tell)
			if ellipsis {
				// Try and get around weird behavior with "::1"
				i++
				continue
			}
			ellipsis = true

			for j := i; j < i+(8-len(spl)+1); j++ {
				addr[j] = 0
			}
			i += 8 - len(spl) + 1
			continue
		}

		parsed, err := strconv.ParseInt(hextet, 16, 64)
		if err != nil {
			return addr, fmt.Errorf("hextet #%d is NaN: %s", si+1, hextet)
		}

		addr[i] = int(parsed)
		i++
	}

	return addr, err
}

func (ip Ipv6Addr) PrintBinary() (s string) {
	for i, hextet := range ip {
		formatted := strconv.FormatInt(int64(hextet), 2)
		for ; 16-len(formatted) > 0; {
			formatted = "0" + formatted
		}
		s += formatted
		if i < 7 {
			s += ":"
		} else if i == 7 && ip.IsCidrFormatted() {
			s += "/"
		} else {
			break
		}
	}

	return s
}

func (ip Ipv6Addr) Print() (s string) {
	for i, hextet := range ip[:8] {
		s += strconv.FormatInt(int64(hextet), 16)
		if i < 7 {
			s += ":"
		} else if i == 7 && ip.IsCidrFormatted() {
			s += "/"
		} else {
			break
		}
	}
	s += strconv.Itoa(ip[8])

	return s
}

func (ip Ipv6Addr) IsCidrFormatted() bool {
	return ip[8] != -1
}

func (ip Ipv6Addr) GetPrefix() int {
	return ip[8]
}

func (ip Ipv6Addr) GetClass() int {
	return 0
}

func (ip Ipv6Addr) GetType() int {
	if ip[0] == 0 && ip[1] == 0 && ip[2] == 0 && ip[3] == 0 && ip[4] == 0 && ip[5] == 0 && ip[6] == 0 {
		if ip[7] == 1 {
			return LOOPBACK
		} else if ip[7] == 0 {
			return UNSPECIFIED_UNICAST
		}
	} else if ip.isOfType(0xFE80, 10) {
		return LINK_LOCAL_UNICAST
	} else if ip.isOfType(0xFC00, 7) {
		return UNIQUE_LOCAL_UNICAST
	} else if ip.isOfType(0, 80) {
		return EMBEDDED_IPV4
	} else if ip.isOfType(0xFF00, 12) {
		return WELL_KNOWN_MULTICAST
	} else if ip.isOfType(0xFF10, 12) {
		return TRANSIENT_MULTICAST
	} else if ip.isOfType(0x2000, 3) {
		return UNICAST
	}// TODO Solicited-Node multicast (FF02:0:0:0:0:1:FF00::/104
	return -1
}

// Takes a hextet and a prefix and returns true if the IP address is of that type
// For example: [2001:db8::ff00:42:8329].isOfType(0x2000, 3) -> true
// See for usage: https://ptgmedia.pearsoncmg.com/images/chap4_9781587144776/elementLinks/04fig06_alt.jpg
func (ip Ipv6Addr) isOfType(mask int, prefix int) bool {
	m := ip.getMask(prefix)
	ok := true

	for i := 0; ok && i < 8; i++ {
		if m[i] & ip[i] != m[i] & mask {
			ok = false
		}
		fmt.Printf("%b & %b (%b) == %b & %b (%b): %v\n", m[i], ip[i], m[i] & ip[i], m[i], mask, m[i] & mask, m[i] & ip[i] == m[i] & mask)
	}

	return ok
}

// See IPv4's implementation
// This is just a helper for IPv6, since it doesn't actually have a subnet mask like IPv4 does
func (ip Ipv6Addr) getMask(prefix int) (ret [8]int) {
	if !ip.IsCidrFormatted() {
		return ret
	}

	for i := 0; i < 8; i++ {
		j := 0xFFFF
		for ; j > 0 && prefix > 0; j >>= 1 {
			prefix--
		}
		ret[i] = 0xFFFF - j
	}
	return ret
}

func (ip Ipv6Addr) PrintNetworkAddress() (s string) {
	netmask := ip.getMask(ip[8])

	for i := 0; i < 8; i++ {
		s += strconv.FormatInt(int64(netmask[i]&ip[i]), 16)

		if i < 7 {
			s += ":"
		} else if i == 7 && ip.IsCidrFormatted() {
			s += "/"
		}
	}
	s += strconv.Itoa(ip[8])

	return s
}