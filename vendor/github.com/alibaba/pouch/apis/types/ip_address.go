// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// IPAddress Address represents an IPv4 or IPv6 IP address.
// swagger:model IPAddress
type IPAddress struct {

	// IP address.
	Addr string `json:"Addr,omitempty"`

	// Mask length of the IP address.
	PrefixLen int64 `json:"PrefixLen,omitempty"`
}

// Validate validates this IP address
func (m *IPAddress) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IPAddress) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IPAddress) UnmarshalBinary(b []byte) error {
	var res IPAddress
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
