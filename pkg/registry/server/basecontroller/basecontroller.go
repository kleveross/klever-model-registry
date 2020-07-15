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
package basecontroller

import (
	"github.com/astaxie/beego"

	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) RendorError(code int, message string) {
	errors.SendError(b.Ctx.ResponseWriter, code, message)
}

func (b *BaseController) RendorCreateSuccess() {
	errors.SendError(b.Ctx.ResponseWriter, errors.CreateSuccessCode, "")
}

func (b *BaseController) RendorDeleteSuccess() {
	errors.SendError(b.Ctx.ResponseWriter, errors.DeleteSuccessCode, "")
}

func (b *BaseController) RendorSuccessData(data interface{}) {
	b.Data["json"] = data
	b.ServeJSON()
}
