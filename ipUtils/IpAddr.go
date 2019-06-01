package ipUtils

type IpAddr interface {
	PrintBinary() (s string)
	Print() (s string)
	IsCidrFormatted() bool
	GetPrefix() int
	GetPrivateClass() int
	GetType() int
	GetNetmask() (ret [4]int)
	PrintNetmask() (s string)
	PrintNetworkAddress() (s string)
}