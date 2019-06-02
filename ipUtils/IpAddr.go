package ipUtils

type IpAddr interface {
	PrintBinary() (s string)
	Print() (s string)
	IsCidrFormatted() bool
	GetPrefix() int
	GetPrivateClass() int
	GetType() int
	PrintNetworkAddress() (s string)
}
