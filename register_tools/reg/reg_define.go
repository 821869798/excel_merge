package reg

type WindowsRegistryKey int

const (
	CLASSES_ROOT WindowsRegistryKey = iota
	CURRENT_USER
	LOCAL_MACHINE
	USERS
	CURRENT_CONFIG
	PERFORMANCE_DATA
)
