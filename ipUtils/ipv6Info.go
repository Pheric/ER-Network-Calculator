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

func (ip Ipv6Addr) GetPrivateClass() int {
	return 0
}

func (ip Ipv6Addr) GetType() int {
	return 0
}

// See IPv4's implementation
// This is just a helper for IPv6, since it doesn't actually have a subnet mask like IPv4 does
func (ip Ipv6Addr) getNetmask() (ret [8]int) {
	if !ip.IsCidrFormatted() {
		return ret
	}

	prefix := ip[8]
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
	netmask := ip.getNetmask()
	fmt.Println(netmask)

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
