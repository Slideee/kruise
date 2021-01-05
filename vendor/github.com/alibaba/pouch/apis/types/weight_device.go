// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// WeightDevice Weight for BlockIO Device
// swagger:model WeightDevice
type WeightDevice struct {

	// Weight Device
	Path string `json:"Path,omitempty"`

	// weight
	// Minimum: 0
	Weight uint16 `json:"Weight,omitempty"`
}

// Validate validates this weight device
func (m *WeightDevice) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateWeight(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WeightDevice) validateWeight(formats strfmt.Registry) error {

	if swag.IsZero(m.Weight) { // not required
		return nil
	}

	if err := validate.MinimumInt("Weight", "body", int64(m.Weight), 0, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *WeightDevice) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WeightDevice) UnmarshalBinary(b []byte) error {
	var res WeightDevice
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
