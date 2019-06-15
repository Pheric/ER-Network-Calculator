package ipUtils

import (
	"fmt"
	"math"
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
	for i, hextet := range ip[:8] {
		s += fmt.Sprintf("%016b", hextet)
		if i < 7 {
			s += ":"
		} else if i == 7 && ip.IsCidrFormatted() {
			s += fmt.Sprintf("/%08b", ip.GetPrefix())
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

// IPv6 doesn't use classes
func (ip Ipv6Addr) GetClass() int {
	return NOT_CLASSFUL
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
	} // TODO Solicited-Node multicast (FF02:0:0:0:0:1:FF00::/104
	return -1
}

// Takes a hextet and a prefix and returns true if the IP address is of that type
// For example: [2001:db8::ff00:42:8329].isOfType(0x2000, 3) -> true
// See for usage: https://ptgmedia.pearsoncmg.com/images/chap4_9781587144776/elementLinks/04fig06_alt.jpg
func (ip Ipv6Addr) isOfType(mask int, prefix int) bool {
	m := getMask(prefix)
	ok := true

	for i := 0; ok && i < 8; i++ {
		if m[i]&ip[i] != m[i]&mask {
			ok = false
		}
	}

	return ok
}

func (ip Ipv6Addr) PrintNetworkAddress() (s string) {
	netmask := getMask(ip[8])

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

// WIP function that auto-subnets an IPv6 address
// @param `nets`: number of desired subnets. Must be <= 65536
// Returns the resultant list of subnets or an error
func (ip Ipv6Addr) Subnet(nets int) ([]IpAddr, error) {
	// f(x) gets the number of bits required to represent x
	f := func(x int) float64 {
		return math.Ceil(math.Log(float64(x)) / math.Log(2)) // change of base formula for log base 2 of x. Go uses Log() as *natural* log (why?!)
	}

	// smask(x) gets a mask of length x
	smask := func(x int) int {
		return int(math.Pow(2, float64(x))) - 1
	}

	nbits := int(f(nets))

	if !ip.IsCidrFormatted() {
		return nil, fmt.Errorf("address has no indicated network")
	} else if nbits > 16 {
		return nil, fmt.Errorf("requiring over 16 bits is not supported")
	} else if ip.GetPrefix()+nbits > 64 {
		return nil, fmt.Errorf("network is too small to subnet properly")
	} else if nets <= 0 {
		return nil, fmt.Errorf("invalid number of subnets")
	}

	// Convert to a network address //
	mask := getMask(ip.GetPrefix())
	for i := 0; i < 8; i++ {
		ip[i] = ip[i] & mask[i]
	}

	// Subnet //

	// Whether the subnet ID will need to be split into two pieces
	split := true
	// Subnet ID may span multiple fields. Here, we get the first of two (this function supports up to 16 bits, so we never have multiple overlaps).
	fieldIndex := int(math.Ceil(float64(ip.GetPrefix())/16)) - 1
	// Find the number of bits available in the (first) field that will be used by the subnet ID
	bitsAvail := 16 - (ip.GetPrefix() % 16)
	if ip.GetPrefix()%16 == 0 || bitsAvail == nbits {
		fieldIndex++ // Address ends right on a delimiter, so we get the next field
		split = false
	}

	// checking fieldIndex: fmt.Printf("fieldIndex: %d / %d = %f -> ceil = %f - 1 = %f [+1? split: %v]\tbitsAvail: %d\tnbits: %d\n", ip.GetPrefix(), 16, float64(ip.GetPrefix()) / 16, math.Ceil(float64(ip.GetPrefix())/16), math.Ceil(float64(ip.GetPrefix())/16)-1, split, bitsAvail, nbits)

	var ret []IpAddr
	max := int(math.Pow(2, float64(nbits)))
	for i := 0; i < max; i += int(math.Floor(float64(max) / float64(nets))) {
		// We now have a different subnet on each iteration. Now, to put that mask into the IP address... //
		addr := ip
		addr[8] = ip.GetPrefix() + nbits

		var fSubnetId, lSubnetId int
		if !split { // The subnet ID will not extend into the next field, so we don't need to split it.
			fSubnetId = i << uint(16-nbits)
		} else { // The subnet ID spills over into a second field, so we must split it
			fSubnetId = i & (smask(nbits) ^ smask(bitsAvail)) >> uint(bitsAvail) // Get the bits in the subnet ID that won't overhang, and right-justify them
			lSubnetId = i & smask(nbits-bitsAvail) << uint(16-(nbits-bitsAvail)) // Do the same for the second half, but left-justified
		}
		addr[fieldIndex] ^= fSubnetId   // Put the first half (if applicable) of this ID into the appropriate field of the address
		addr[fieldIndex+1] ^= lSubnetId // Do the same for the second (overhanging) half, even if unset (not required)

		ret = append(ret, addr)
	}

	return ret, nil
}

// Returns this address incremented by `amt`, starting from the rightmost hextet and working left
func (ip Ipv6Addr) increment(amt float64) (addr Ipv6Addr) {
	addr = ip
	for i := 7; i >= 0; i-- {
		if amt >= 0xFFFF {
			amt -= float64(0xFFFF - addr[i])
			addr[i] = 0xFFFF
		} else {
			addr[i] += int(amt) // cast OK because we're under 65k
			break
		}
	}

	return
}

// See IPv4's implementation
// This is just a helper for IPv6, since it doesn't actually have a subnet mask like IPv4 does
func getMask(prefix int) (ret Ipv6Addr) {
	for i := 0; i < 8; i++ {
		j := 0xFFFF
		for ; j > 0 && prefix > 0; j >>= 1 {
			prefix--
		}
		ret[i] = 0xFFFF - j
	}
	return ret
}
