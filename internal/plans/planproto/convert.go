// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package planproto

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/plans"
)

func NewPath(src cty.Path) (*Path, error) {
	ret := &Path{
		Steps: make([]*Path_Step, len(src)),
	}
	for i, srcStep := range src {
		step, err := NewPathStep(srcStep)
		if err != nil {
			return nil, fmt.Errorf("step %d: %w", i, err)
		}
		ret.Steps[i] = step
	}
	return ret, nil
}

func NewPathStep(step cty.PathStep) (*Path_Step, error) {
	switch s := step.(type) {
	case cty.IndexStep:
		value, err := plans.NewDynamicValue(s.Key, s.Key.Type())
		if err != nil {
			return nil, err
		}
		return &Path_Step{
			Selector: &Path_Step_ElementKey{
				ElementKey: NewPlanDynamicValue(value),
			},
		}, nil
	case cty.GetAttrStep:
		return &Path_Step{
			Selector: &Path_Step_AttributeName{
				AttributeName: s.Name,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported step type %t", step)
	}
}

func NewPlanDynamicValue(dv plans.DynamicValue) *DynamicValue {
	if dv == nil {
		// protobuf can't represent nil, so we'll represent it as a
		// DynamicValue that has no serializations at all.
		return &DynamicValue{}
	}
	return &DynamicValue{
		Msgpack: []byte(dv),
	}
}
