package version

// Version is the application version.
//
// It defaults to the value below and is overridden at build time with:
//
//	go build -ldflags "-X github.com/sthbryan/ftm/internal/version.Version=0.10.1"
//
// Keep this a var, not a const: ldflags cannot patch constants.
var Version = "0.10.0"
