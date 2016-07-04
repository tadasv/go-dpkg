package main

import (
	"fmt"
	"github.com/tadasv/go-dpkg"
	"os"
)

func main() {
	fileName := "/var/lib/dpkg/status"

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	}

	packages, err := dpkg.ReadPackagesFromFile(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, pkg := range packages {
		fmt.Printf("%-32s%-32s %s\n", pkg.Package, pkg.Version, pkg.Status)
	}

	fmt.Printf("Read %d packages\n", len(packages))
}
