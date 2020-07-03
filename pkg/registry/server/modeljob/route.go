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
package modeljob

import (
	"github.com/astaxie/beego"
)

// RegisterRoutes for modeljob APIs
func RegisterRoutes() {
	beego.Router("/api/modeljob/v1alpha1", &modeljobController{}, "post:Create")
	beego.Router("/api/modeljob/v1alpha1/:modeljob_id", &modeljobController{}, "get:Get")
	beego.Router("/api/modeljob/v1alpha1/:modeljob_id", &modeljobController{}, "delete:Delete")
	beego.Router("/api/modeljob/v1alpha1", &modeljobController{}, "get:List")
}
