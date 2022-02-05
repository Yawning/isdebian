// Copyright (c) 2022 Yawning Angel <yawning at schwanenlied dot me>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package isdebian // import "gitlab.com/yawning/isdebian.git"

import (
	"bytes"
	"os"
	"runtime"
	"syscall"
)

// IsDebian returns true iff the system is likely running Debian, ignoring
// derivatives.
func IsDebian() bool {
	// Skip if the kernel isn't something Debian targets.
	switch runtime.GOOS {
	case "linux", "freebsd", "netbsd", "hurd":
	default:
		return false
	}

	// Check uname(2).
	var uts syscall.Utsname
	if err := syscall.Uname(&uts); err == nil {
		toBytes := func(src []int8) []byte {
			dst := make([]byte, 0, len(src))
			for _, v := range src {
				dst = append(dst, byte(v))
			}
			return dst
		}
		if bytes.Contains(toBytes(uts.Version[:]), []byte("Debian")) {
			return true
		}
	}

	// Check `os-release`, in both places where it can exist.
	for _, fn := range []string{
		"/etc/os-release",
		"/usr/lib/os-release",
	} {
		if b, err := readFile(fn); err == nil {
			// Derivatives report `ID_LIKE=debian`, but ignore for now.
			if bytes.Contains(b, []byte("ID=debian")) {
				return true
			}
		}
	}

	return false
}

// IsWhonix returns true iff the system is likely running Whonix.
func IsWhonix() bool {
	// Skip if the kernel isn't something Debian targets.
	switch runtime.GOOS {
	case "linux", "freebsd", "netbsd", "hurd":
	default:
		return false
	}

	// Check for the various files indicating that this may be whonix.
	for _, fn := range []string{
		"/usr/share/whonix/marker",
		"/etc/whonix_version",
		"/usr/share/anon-gw-base-files/gateway",
		"/usr/share/anon-ws-base-files/workstation",
	} {
		if _, err := os.Stat(fn); err == nil {
			return true
		}
	}

	return false
}
