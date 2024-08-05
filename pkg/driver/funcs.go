package driver

import (
	"github.com/friendly-fhir/fhenix/pkg/model"
	"github.com/friendly-fhir/fhenix/pkg/transform"
)

var Funcs = transform.Funcs{
	"commonbase": CommonBase,
}

func commonBase(t1, t2 *model.Type) *model.Type {
	if t1 == t2 {
		return t1
	}
	if t1.Base == nil && t2.Base == nil {
		return nil
	}
	if t1.Base != nil {
		t1 = t1.Base
	}
	if t2.Base != nil {
		t2 = t2.Base
	}
	return commonBase(t1, t2)
}

func CommonBase(types []*model.Type) *model.Type {
	if len(types) == 0 {
		return nil
	}
	base := types[0]
	for _, t := range types[1:] {
		base = commonBase(base, t)
	}
	return base
}
