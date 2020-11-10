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

package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kleveross/klever-model-registry/cmd/modeljob-operator/app"
	"github.com/kleveross/klever-model-registry/cmd/modeljob-operator/app/options"
	// +kubebuilder:scaffold:imports
)

func main() {
	s := options.NewServerOption()
	s.AddFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New())

	if err := app.Run(s); err != nil {
		log.Fatalf("%v\n", err)
	}
}
