package version

// Build information. Populated at build-time.
var (
	Version   = "unknown"
	Revision  = "unknown"
	Branch    = "unknown"
	BuildUser = "unknown"
	BuildDate = "unknown"
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildUser": BuildUser,
	"buildDate": BuildDate,
}

func IsValidVersion() bool {
	return Version != "unknown"
}
