// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Anexia anexia
//
// swagger:model Anexia
type Anexia struct {

	// datacenter
	Datacenter string `json:"datacenter,omitempty"`

	// enabled
	Enabled bool `json:"enabled,omitempty"`

	// Token is used to authenticate with the Anexia API.
	Token string `json:"token,omitempty"`
}

// Validate validates this anexia
func (m *Anexia) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this anexia based on context it is used
func (m *Anexia) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Anexia) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Anexia) UnmarshalBinary(b []byte) error {
	var res Anexia
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
