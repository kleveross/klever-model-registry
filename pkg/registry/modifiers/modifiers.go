// +nirvana:api=modifiers:"Modifiers"

package modifiers

import "github.com/caicloud/nirvana/service"

// Modifiers returns a list of modifiers.
func Modifiers() []service.DefinitionModifier {
	return []service.DefinitionModifier{
		service.FirstContextParameter(),
		service.ConsumeNoneForHTTPGet(),
		service.ConsumeNoneForHTTPDelete(),
		service.ProduceNoneForHTTPDelete(),
		service.ConsumeNoneForHTTPHead(),
		service.ProduceNoneForHTTPHead(),
		service.ConsumeAllIfConsumesIsEmpty(),
		service.ProduceAllIfProducesIsEmpty(),
	}
}
