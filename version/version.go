package version

import "fmt"

var (
	BinaryVersion string
	GoVersion     string
	GitLastLog    string
)

func Version() {
	fmt.Println("BinaryVersion:", BinaryVersion)
	fmt.Println("GoVersion:", GoVersion)
	fmt.Println("GitLastLog:", GitLastLog)
}

func Help() {
	fmt.Println("The commands are:")
	fmt.Println("version       see all versions")
}
