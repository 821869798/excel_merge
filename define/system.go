package define

type SystemType int

const (
	SystemTypeNone SystemType = iota
	SystemTypeWindows
	SystemTypeLinux
	SystemTypeMac
)
