package dpkg

import (
	"bytes"
	"testing"
)

func assertEqual(t *testing.T, expected string, current string) {
	if expected != current {
		t.Errorf("Invalid value. Expected '%v', got '%v'", expected, current)
	}
}

func TestMapToPackage(t *testing.T) {
	parser := NewParser(nil)
	m := map[string]string{
		"Package":        "package",
		"Version":        "version",
		"Section":        "section",
		"Installed-Size": "123",
		"Maintainer":     "tadas",
		"Status":         "status",
		"Source":         "source",
		"Architecture":   "amd64",
		"Multi-Arch":     "same",
		"Depends":        "depends",
		"Pre-Depends":    "predepends",
		"Description":    "desc",
		"Homepage":       "home",
		"Priority":       "priority",
	}

	pkg, err := parser.mapToPackage(m)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assertEqual(t, "package", pkg.Package)
	assertEqual(t, "version", pkg.Version)
	assertEqual(t, "section", pkg.Section)
	assertEqual(t, "tadas", pkg.Maintainer)
	assertEqual(t, "status", pkg.Status)
	assertEqual(t, "source", pkg.Source)
	assertEqual(t, "amd64", pkg.Architecture)
	assertEqual(t, "same", pkg.MultiArch)
	assertEqual(t, "depends", pkg.Depends)
	assertEqual(t, "predepends", pkg.PreDepends)
	assertEqual(t, "desc", pkg.Description)
	assertEqual(t, "home", pkg.Homepage)
	assertEqual(t, "priority", pkg.Priority)

	if pkg.InstalledSize != 123 {
		t.Errorf("Invalid size: %v", pkg.InstalledSize)
	}

}

func TestParseLineHandlesEmptyString(t *testing.T) {
	parser := NewParser(nil)
	key, value := parser.parseLine("")

	assertEqual(t, "", key)
	assertEqual(t, "", value)
}

func TestParseLineHandlesNewLine(t *testing.T) {
	parser := NewParser(nil)
	key, value := parser.parseLine("\n")

	assertEqual(t, "", key)
	assertEqual(t, "", value)
}

func TestParseLineHandlesMultilineValue(t *testing.T) {
	parser := NewParser(nil)
	key, value := parser.parseLine(" some: value\n")

	assertEqual(t, "", key)
	assertEqual(t, " some: value", value)
}

func TestParseLineHandlesKeyValue(t *testing.T) {
	parser := NewParser(nil)
	key, value := parser.parseLine("Key: value\n")

	assertEqual(t, "Key", key)
	assertEqual(t, " value", value)
}

func TestParseAllBlankLines(t *testing.T) {
	data := "\n\n\n\n\n"
	reader := bytes.NewBufferString(data)
	parser := NewParser(reader)
	packages := parser.Parse()
	if len(packages) != 0 {
		t.Errorf("Expected 0 packages, got: %v", len(packages))
	}
}

func TestParseValidData(t *testing.T) {
	data := `Package: libquadmath0
Status: install ok installed
Priority: optional
Section: libs
Installed-Size: 275
Maintainer: Debian GCC Maintainers <debian-gcc@lists.debian.org>
Architecture: amd64
Multi-Arch: same
Source: gcc-4.9
Version: 4.9.2-10
Depends: gcc-4.9-base (= 4.9.2-10), libc6 (>= 2.14)
Pre-Depends: multiarch-support
Description: GCC Quad-Precision Math Library
 A library, which provides quad-precision mathematical functions on targets
 supporting the __float128 datatype. The library is used to provide on such
 targets the REAL(16) type in the GNU Fortran compiler.
Homepage: http://gcc.gnu.org/

Package: libedit2
Status: install ok installed
Priority: standard
Section: libs
Installed-Size: 277
Maintainer: LLVM Packaging Team <pkg-llvm-team@lists.alioth.debian.org>
Architecture: amd64
Multi-Arch: same
Source: libedit
Version: 3.1-20140620-2
Depends: libbsd0 (>= 0.0), libc6 (>= 2.17), libtinfo5
Pre-Depends: multiarch-support
Description: BSD editline and history libraries
 Command line editor library provides generic line editing,
 history, and tokenization functions.
 .
 It slightly resembles GNU readline.
Homepage: http://www.thrysoee.dk/editline/`

	reader := bytes.NewBufferString(data)
	parser := NewParser(reader)
	packages := parser.Parse()

	if len(packages) != 2 {
		t.Errorf("Expected 2 packages, got: %v", len(packages))
	}

	pkg := packages[0]
	assertEqual(t, "libquadmath0", pkg.Package)
	assertEqual(t, "4.9.2-10", pkg.Version)
	assertEqual(t, "libs", pkg.Section)
	assertEqual(t, "Debian GCC Maintainers <debian-gcc@lists.debian.org>", pkg.Maintainer)
	assertEqual(t, "install ok installed", pkg.Status)
	assertEqual(t, "gcc-4.9", pkg.Source)
	assertEqual(t, "amd64", pkg.Architecture)
	assertEqual(t, "same", pkg.MultiArch)
	assertEqual(t, "gcc-4.9-base (= 4.9.2-10), libc6 (>= 2.14)", pkg.Depends)
	assertEqual(t, "multiarch-support", pkg.PreDepends)
	assertEqual(t, "http://gcc.gnu.org/", pkg.Homepage)
	assertEqual(t, "optional", pkg.Priority)
	if pkg.InstalledSize != 275 {
		t.Errorf("Incorrect size: %v", pkg.InstalledSize)
	}

	pkg = packages[1]
	assertEqual(t, "libedit2", pkg.Package)
	assertEqual(t, "3.1-20140620-2", pkg.Version)
	assertEqual(t, "libs", pkg.Section)
	assertEqual(t, "LLVM Packaging Team <pkg-llvm-team@lists.alioth.debian.org>", pkg.Maintainer)
	assertEqual(t, "install ok installed", pkg.Status)
	assertEqual(t, "libedit", pkg.Source)
	assertEqual(t, "amd64", pkg.Architecture)
	assertEqual(t, "same", pkg.MultiArch)
	assertEqual(t, "libbsd0 (>= 0.0), libc6 (>= 2.17), libtinfo5", pkg.Depends)
	assertEqual(t, "multiarch-support", pkg.PreDepends)
	assertEqual(t, "http://www.thrysoee.dk/editline/", pkg.Homepage)
	assertEqual(t, "standard", pkg.Priority)
	if pkg.InstalledSize != 277 {
		t.Errorf("Incorrect size: %v", pkg.InstalledSize)
	}
}
