package ipUtils

type IpAddr interface {
	PrintBinary() (s string)
	Print() (s string)
	IsCidrFormatted() bool
	GetPrefix() int
	GetClass() int
	GetType() int
	PrintNetworkAddress() (s string)
	Subnet(nets int) ([]IpAddr, error)
}

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

	// IPv6 specific
	LINK_LOCAL_UNICAST
	UNIQUE_LOCAL_UNICAST
	UNSPECIFIED_UNICAST
	EMBEDDED_IPV4
	WELL_KNOWN_MULTICAST
	TRANSIENT_MULTICAST
)