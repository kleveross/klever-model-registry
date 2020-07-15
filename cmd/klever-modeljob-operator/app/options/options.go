/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import "flag"

// ServerOption is the main context object for the controller manager.
type ServerOption struct {
	PrintVersion bool

	MetricsAddr string

	EnableLeaderElection bool
}

// NewServerOption creates a new CMServer with a default config.
func NewServerOption() *ServerOption {
	s := ServerOption{}
	return &s
}

// AddFlags adds flags for a specific CMServer to the specified FlagSet.
func (s *ServerOption) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&s.PrintVersion, "version", false, "Show version and quit")

	flag.StringVar(&s.MetricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")

	flag.BoolVar(&s.EnableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
}
