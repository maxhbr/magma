// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NetworkProbeTask Network Probe Task
//
// swagger:model network_probe_task
type NetworkProbeTask struct {

	// task details
	// Required: true
	TaskDetails *NetworkProbeTaskDetails `json:"task_details"`

	// task id
	// Required: true
	TaskID NetworkProbeTaskID `json:"task_id"`
}

// Validate validates this network probe task
func (m *NetworkProbeTask) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTaskDetails(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTaskID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkProbeTask) validateTaskDetails(formats strfmt.Registry) error {

	if err := validate.Required("task_details", "body", m.TaskDetails); err != nil {
		return err
	}

	if m.TaskDetails != nil {
		if err := m.TaskDetails.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task_details")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task_details")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkProbeTask) validateTaskID(formats strfmt.Registry) error {

	if err := validate.Required("task_id", "body", NetworkProbeTaskID(m.TaskID)); err != nil {
		return err
	}

	if err := m.TaskID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("task_id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("task_id")
		}
		return err
	}

	return nil
}

// ContextValidate validate this network probe task based on the context it is used
func (m *NetworkProbeTask) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTaskDetails(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTaskID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkProbeTask) contextValidateTaskDetails(ctx context.Context, formats strfmt.Registry) error {

	if m.TaskDetails != nil {
		if err := m.TaskDetails.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task_details")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task_details")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkProbeTask) contextValidateTaskID(ctx context.Context, formats strfmt.Registry) error {

	if err := m.TaskID.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("task_id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("task_id")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *NetworkProbeTask) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetworkProbeTask) UnmarshalBinary(b []byte) error {
	var res NetworkProbeTask
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
