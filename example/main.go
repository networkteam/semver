package main

import (
	"fmt"

	"github.com/networkteam/semver"
)

func main() {
	v, err := semver.ParseVersion("1.0.0-alpha.1+001")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed Version", v.Major, v.Minor, v.Patch, v.PreRelease, v.Build)
		fmt.Println("String Representation:", v.String())
	}
}
