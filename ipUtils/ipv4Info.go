package ipUtils

import (
	"fmt"
	"regexp"
	"strconv"
)

type Ipv4Addr [5]int
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

// Parses an IPv4 address out of a string. Must not have any protocol or port.
func ParseIpv4(str string) (addr Ipv4Addr, err error) {
	if str == "" {
		return addr, fmt.Errorf("input string is empty")
	}

	// Regex didn't work when I tried to compress it.. so I guess we get to use the expanded version. Written by hand
	ipv4WithCidrRegex := regexp.MustCompile(`(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)(?:\/([1-3]?\d))?`)
	matches := ipv4WithCidrRegex.FindAllStringSubmatch(str, -1)[0][1:]
	if l := len(matches); l == 0 || l > 5 || l < 4 {
		return addr, fmt.Errorf("invalid format by regex")
	}
	if len(matches) == 4 {
		matches = append(matches, "-1")
	}

	for i, octet := range matches {
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
		} else if i == 3 && ip.IsCidrFormatted() {
			s += "/"
		}
	}

	return s
}

// Returns the ip address
func (ip Ipv4Addr) Print() (s string) {
	s = fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	if ip.IsCidrFormatted() {
		s += fmt.Sprintf("/%d", ip[4])
	}

	return s
}

func (ip Ipv4Addr) IsCidrFormatted() bool {
	return ip[4] != -1
}

// If there's no prefix and the address is a normal private IP address, this will return its network's class
// Otherwise, it prefers the prefix.
// For example, the address '112.17.100.45/16' is a class B; '240.23.18.1' is a class E
func (ip Ipv4Addr) GetPrivateClass() int {
	if ip[0] == 10 || ip[0] == 127 || (ip.IsCidrFormatted() && ip[4] == 8 && !(ip[0] >= 240)) {
		return CLASS_A
	} else if ip[0] == 172 && (ip[1] >= 16 && ip[1] <= 31) || (ip.IsCidrFormatted() && ip[4] == 16) {
		return CLASS_B
	} else if ip[0] == 192 && ip[1] == 168 || (ip.IsCidrFormatted() && ip[4] == 24) {
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

func (ip Ipv4Addr) GetNetmask() (ret [4]int) {
	if !ip.IsCidrFormatted() {
		return ret
	}

	prefix := ip[4]
	for i := 0; i < 4; i++ {
		j := 255
		for ; j > 0 && prefix > 0; j >>= 1 {
			prefix--
		}
		ret[i] = 255 - j
	}
	return ret
}

func (ip Ipv4Addr) PrintNetmask() (s string) {
	netmask := ip.GetNetmask()
	return fmt.Sprintf("%d.%d.%d.%d", netmask[0], netmask[1], netmask[2], netmask[3])
}

func (ip Ipv4Addr) PrintNetworkAddress() (s string) {
	netmask := ip.GetNetmask()

	for i := 0; i < 4; i++ {
		s += strconv.Itoa(ip[i] & netmask[i])
		if i < 3 {
			s += "."
		}
	}

	return s
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