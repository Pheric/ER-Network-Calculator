package ipUtils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Ipv4Addr [4]int
const (
	NOT_CLASSFUL = iota
	CLASS_A
	CLASS_B
	CLASS_C
	CLASS_E
	PUBLIC
	PRIVATE
	SPECIAL
	LOOPBACK
	APIPA
	BROADCAST
	MULTICAST
	NETWORK
	UNICAST
)

// Parses an IPv4 address out of a string. Must not be in CIDR notation, have any protocol, or have a port.
func ParseIpv4(str string) (addr Ipv4Addr, err error) {
	if str == "" {
		return addr, fmt.Errorf("input string is empty")
	}

	if net.ParseIP(str) == nil {
		return addr, fmt.Errorf("ip format is invalid")
	}

	split := strings.Split(str, ".")
	for i, octet := range split {
		parsed, err := strconv.Atoi(octet)
		if err != nil {
			return addr, fmt.Errorf("octet #%d is NaN", i)
		}

		addr[i] = parsed
	}

	return addr, nil
}

// Returns the padded binary representation of the IP address
func (ip Ipv4Addr) PrintBinary() (s string) {
	for i, octet := range ip {
		formatted := strconv.FormatInt(int64(octet), 2)
		for ; 8-len(formatted) > 0; {
			formatted = "0" + formatted
		}
		s += formatted
		if i < 3 {
			s += "."
		}
	}

	return s
}

// If the address is a normal private IP address, this will return its network's class
// For example, the address '192.168.1.1' is a class B; '240.23.18.1' is a class E
func (ip Ipv4Addr) GetPrivateClass() int {
	if ip[0] == 10 || ip[0] == 127 {
		return CLASS_A
	} else if ip[0] == 172 && (ip[1] >= 16 && ip[1] <= 31) {
		return CLASS_B
	} else if ip[0] == 192 && ip[1] == 168 {
		return CLASS_C
	} else if ip[0] >= 240 {
		return CLASS_E
	} else {
		return NOT_CLASSFUL
	}
}

func (ip Ipv4Addr) GetType() int {
	if ip[0] == 127 {
		return LOOPBACK
	} else if ip[0] == 169 && ip[1] == 254 {
		return APIPA
	} else if ip[0] >= 224 && ip[0] <= 239 {
		return MULTICAST
	} else if ip[3] == 255 {
		return BROADCAST
	} else if ip[3] == 0 { // TODO: support subnetting
		return NETWORK
	} else {
		return UNICAST
	}
}

// Describes the IP address as best it can:
// - Public / private
// - Class if private
// - Type (APIPA / loopback / multicast / broadcast / unicast / network)
func (ip Ipv4Addr) Describe() (ret []int) {
	class := ip.GetPrivateClass()
	ret = append(ret, class)

	if class != NOT_CLASSFUL {
		ret = append(ret, PRIVATE)
	} else {
		ret = append(ret, PUBLIC)
	}

	ret = append(ret, ip.GetType())

	return ret
}

// TODO
/*// Returns 'true' if the IP address is special, meaning:
// - it is a normal private IP (see GetPrivateClass())
// - it is reserved or commonly used in a special way according to Wikipedia's IPv4#Addressing section
func (ip Ipv4Addr) isSpecial() bool {
	// See https://en.wikipedia.org/wiki/IPv4#Addressing
	return ip.GetPrivateClass() != NOT_CLASSFUL &&
		!(ip[0] == 100 && (ip[1] >= 64 && ip[1] <= 127)) && // "Shared address space for communications between a service provider and its subscribers when using a carrier-grade NAT."
		!(ip[0] == 192 && ip[1] == 0 && ip[2] == 0) && // "IETF Protocol Assignments."
		!(ip[0] == 198 && (ip[1] == 18 || ip[1] == 19)) // "Used for benchmark testing of inter-network communications between two separate subnets."
}*/