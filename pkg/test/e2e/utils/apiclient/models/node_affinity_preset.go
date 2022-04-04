// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NodeAffinityPreset node affinity preset
//
// swagger:model NodeAffinityPreset
type NodeAffinityPreset struct {

	// key
	Key string `json:"Key,omitempty"`

	// type
	Type string `json:"Type,omitempty"`

	// values
	Values []string `json:"Values"`
}

// Validate validates this node affinity preset
func (m *NodeAffinityPreset) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this node affinity preset based on context it is used
func (m *NodeAffinityPreset) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NodeAffinityPreset) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NodeAffinityPreset) UnmarshalBinary(b []byte) error {
	var res NodeAffinityPreset
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
