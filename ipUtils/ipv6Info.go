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

	// Test:
	// 0    1   2 3 4 5    6    7 [i]
	// 2001:0db8::    ff00:0042:8329
	i := 0
	for si, hextet := range spl { // Don't you quibble with me!
		//fmt.Printf("i: %d\tsi: %d\thextet: `%s`\tellipsis detected: %v\n", i, si, hextet, hextet == "")
		if hextet == "" {
			//fmt.Printf("ellipsis detected. position: %d\tloop until: %d\n", i, i+(8-len(spl)))
			// Position of "ellipsis" (this is what Go's net package calls the :: as far as I can tell)
			for j := i; j < i + (8 - len(spl) + 1); j++ {
				//fmt.Printf("\tj: %d\n", j)
				addr[j] = 0
			}
			i += 8 - len(spl) + 1
			//fmt.Printf("i -> %d\tsi: %d\taddr: `%v`\n", i, si, addr)
			continue
		}

		parsed, err := strconv.ParseInt(hextet, 16, 64)
		if err != nil {
			return addr, fmt.Errorf("hextet #%d is NaN: %s", si + 1, hextet)
		}

		addr[i] = int(parsed)
		i++
	}

	// TODO
	return addr, err
}

func (ip Ipv6Addr) PrintBinary() (s string) {
	return s
}

func (ip Ipv6Addr) Print() (s string) {
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

func (ip Ipv6Addr) GetNetmask() (ret [4]int) {
	return [4]int{0, 0, 0, 0}
}

func (ip Ipv6Addr) PrintNetmask() (s string) {
	return s
}

func (ip Ipv6Addr) PrintNetworkAddress() (s string) {
	return s
}