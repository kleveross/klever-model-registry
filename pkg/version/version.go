// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"fmt"
	"runtime"
)

// Following values should be substituted with a real value during build.
var (
	// VERSION is the app-global version string.
	VERSION = "UNKNOWN"

	// COMMIT is the app-global git sha string.
	COMMIT = "UNKNOWN"

	// REPOROOT is the app-global repository path string.
	REPOROOT = "UNKNOWN"
)

// PrintVersion prints versions from the array returned by Info().
func PrintVersion() {
	for _, i := range Info() {
		fmt.Printf("%v\n", i)
	}
}

// Info returns an array of various service versions
func Info() []string {
	return []string{
		fmt.Sprintf("Version: %s", VERSION),
		fmt.Sprintf("Git SHA: %s", COMMIT),
		fmt.Sprintf("Repo Root: %s", REPOROOT),
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
