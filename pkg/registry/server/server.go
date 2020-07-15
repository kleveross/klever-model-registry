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
package server

import (
	"fmt"

	"github.com/astaxie/beego"

	eventcontroller "github.com/kleveross/klever-model-registry/pkg/registry/server/event"
	"github.com/kleveross/klever-model-registry/pkg/registry/server/harbor"
	logcontroller "github.com/kleveross/klever-model-registry/pkg/registry/server/log"
	modeljobcontroller "github.com/kleveross/klever-model-registry/pkg/registry/server/modeljob"
)

func routePath(url string) string {
	return fmt.Sprintf("/apis/v1alpha1%v", url)
}

// RegisterRoutes register all routes
func RegisterRoutes() {

	harbor.RegisterRoutes() // Register harbor route

	// Modeljob route
	beego.Router(routePath("/modeljobs"), &modeljobcontroller.ModeljobController{}, "post:Create")
	beego.Router(routePath("/modeljobs/:modeljob_id"), &modeljobcontroller.ModeljobController{}, "get:Get")
	beego.Router(routePath("/modeljobs/:modeljob_id"), &modeljobcontroller.ModeljobController{}, "delete:Delete")
	beego.Router(routePath("/modeljobs"), &modeljobcontroller.ModeljobController{}, "get:List")

	// Event route
	beego.Router(routePath("/namespaces/:namespace/modeljobs/:modeljob_id/events"),
		&eventcontroller.EventController{}, "get:GetModelJobEvents")

	//Log route
	beego.Router(routePath("/namespaces/:namespace/pods/:pod_id/log"), &logcontroller.LogController{}, "get:Get")
}
