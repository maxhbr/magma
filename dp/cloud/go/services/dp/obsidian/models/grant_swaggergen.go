// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Grant grant
//
// swagger:model grant
type Grant struct {

	// bandwidth mhz
	// Example: 20
	// Required: true
	BandwidthMhz int64 `json:"bandwidth_mhz"`

	// frequency mhz
	// Example: 3600
	// Required: true
	// Maximum: 3700
	// Minimum: 3550
	FrequencyMhz int64 `json:"frequency_mhz"`

	// grant expire time
	// Format: date-time
	GrantExpireTime strfmt.DateTime `json:"grant_expire_time,omitempty"`

	// max eirp
	// Required: true
	// Maximum: 37
	// Minimum: -137
	MaxEirp float64 `json:"max_eirp"`

	// state
	// Required: true
	// Enum: [granted guthorized]
	State string `json:"state"`

	// transmit expire time
	// Format: date-time
	TransmitExpireTime strfmt.DateTime `json:"transmit_expire_time,omitempty"`
}

// Validate validates this grant
func (m *Grant) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBandwidthMhz(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFrequencyMhz(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGrantExpireTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMaxEirp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransmitExpireTime(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Grant) validateBandwidthMhz(formats strfmt.Registry) error {

	if err := validate.Required("bandwidth_mhz", "body", int64(m.BandwidthMhz)); err != nil {
		return err
	}

	return nil
}

func (m *Grant) validateFrequencyMhz(formats strfmt.Registry) error {

	if err := validate.Required("frequency_mhz", "body", int64(m.FrequencyMhz)); err != nil {
		return err
	}

	if err := validate.MinimumInt("frequency_mhz", "body", m.FrequencyMhz, 3550, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("frequency_mhz", "body", m.FrequencyMhz, 3700, false); err != nil {
		return err
	}

	return nil
}

func (m *Grant) validateGrantExpireTime(formats strfmt.Registry) error {
	if swag.IsZero(m.GrantExpireTime) { // not required
		return nil
	}

	if err := validate.FormatOf("grant_expire_time", "body", "date-time", m.GrantExpireTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Grant) validateMaxEirp(formats strfmt.Registry) error {

	if err := validate.Required("max_eirp", "body", float64(m.MaxEirp)); err != nil {
		return err
	}

	if err := validate.Minimum("max_eirp", "body", m.MaxEirp, -137, false); err != nil {
		return err
	}

	if err := validate.Maximum("max_eirp", "body", m.MaxEirp, 37, false); err != nil {
		return err
	}

	return nil
}

var grantTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["granted","guthorized"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		grantTypeStatePropEnum = append(grantTypeStatePropEnum, v)
	}
}

const (

	// GrantStateGranted captures enum value "granted"
	GrantStateGranted string = "granted"

	// GrantStateGuthorized captures enum value "guthorized"
	GrantStateGuthorized string = "guthorized"
)

// prop value enum
func (m *Grant) validateStateEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, grantTypeStatePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Grant) validateState(formats strfmt.Registry) error {

	if err := validate.RequiredString("state", "body", m.State); err != nil {
		return err
	}

	// value enum
	if err := m.validateStateEnum("state", "body", m.State); err != nil {
		return err
	}

	return nil
}

func (m *Grant) validateTransmitExpireTime(formats strfmt.Registry) error {
	if swag.IsZero(m.TransmitExpireTime) { // not required
		return nil
	}

	if err := validate.FormatOf("transmit_expire_time", "body", "date-time", m.TransmitExpireTime.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this grant based on context it is used
func (m *Grant) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Grant) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Grant) UnmarshalBinary(b []byte) error {
	var res Grant
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
