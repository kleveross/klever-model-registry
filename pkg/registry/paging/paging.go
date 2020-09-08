package paging

import (
	"context"
	"fmt"
	"reflect"

	"github.com/caicloud/nirvana/definition"
	"github.com/caicloud/nirvana/errors"
	"github.com/caicloud/nirvana/operators/validator"
)

// ListOption is params in query which is used for paging.
type ListOption struct {
	Start int  `source:"Query,start"`
	Limit *int `source:"Query,limit"`
}

// ListMeta describes list of objects, i.e. holds information about pagination options set for the list.
type ListMeta struct {
	// Total number of items on the list. Used for pagination.
	TotalItems int `json:"totalItems"`
}

type Result struct {
	ListMeta `json:"metadata"`
	Items    []interface{} `json:"items"`
}

// ValidateListOption is validating the ListOption.
func ValidateListOption(opts *ListOption) error {
	if opts.Start < 0 {
		return fmt.Errorf("start should greater or equal 0")
	}
	if opts.Limit != nil {
		if *opts.Limit < 0 {
			return fmt.Errorf("limit should greater or equal 0")
		}
	}

	return nil
}

func toSlice(items interface{}) (out []interface{}) {
	value := reflect.ValueOf(items)
	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		len := value.Len()
		out = make([]interface{}, len)
		for i := 0; i < len; i++ {
			out[i] = value.Index(i).Interface()
		}
	} else {
		out = append(out, items)
	}
	return
}

// Page is page the input items datas.
func Page(items interface{}, opts *ListOption) Result {
	var result = Result{
		Items: make([]interface{}, 0),
	}
	if items == nil {
		return result
	}
	slice := toSlice(items)
	result.TotalItems = len(slice)

	if opts.Start >= result.TotalItems {
		return result
	}

	result.Items = slice[opts.Start:]
	if opts.Limit != nil && *opts.Limit < len(result.Items) {
		result.Items = result.Items[:*opts.Limit]
	}

	return result
}

// PageDefinitionParameter is used for apis/descriptors to define Definition's Parameter.
func PageDefinitionParameter() definition.Parameter {
	return definition.Parameter{
		Source: definition.Auto,
		Operators: []definition.Operator{
			validator.NewCustom(
				func(ctx context.Context, opts *ListOption) error {
					err := ValidateListOption(opts)
					if err != nil {
						return errors.BadRequest.Error(err.Error())
					}
					return nil
				}, "validate paging"),
		},
	}
}
