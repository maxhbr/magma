// Code generated by go-swagger; DO NOT EDIT.

package subscribers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "magma/orc8r/cloud/go/obsidian/swagger/v1/models"
)

// PostLTENetworkIDSubscribersReader is a Reader for the PostLTENetworkIDSubscribers structure.
type PostLTENetworkIDSubscribersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostLTENetworkIDSubscribersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostLTENetworkIDSubscribersCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPostLTENetworkIDSubscribersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostLTENetworkIDSubscribersCreated creates a PostLTENetworkIDSubscribersCreated with default headers values
func NewPostLTENetworkIDSubscribersCreated() *PostLTENetworkIDSubscribersCreated {
	return &PostLTENetworkIDSubscribersCreated{}
}

/*PostLTENetworkIDSubscribersCreated handles this case with default header values.

Success
*/
type PostLTENetworkIDSubscribersCreated struct {
}

func (o *PostLTENetworkIDSubscribersCreated) Error() string {
	return fmt.Sprintf("[POST /lte/{network_id}/subscribers][%d] postLteNetworkIdSubscribersCreated ", 201)
}

func (o *PostLTENetworkIDSubscribersCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPostLTENetworkIDSubscribersDefault creates a PostLTENetworkIDSubscribersDefault with default headers values
func NewPostLTENetworkIDSubscribersDefault(code int) *PostLTENetworkIDSubscribersDefault {
	return &PostLTENetworkIDSubscribersDefault{
		_statusCode: code,
	}
}

/*PostLTENetworkIDSubscribersDefault handles this case with default header values.

Unexpected Error
*/
type PostLTENetworkIDSubscribersDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the post LTE network ID subscribers default response
func (o *PostLTENetworkIDSubscribersDefault) Code() int {
	return o._statusCode
}

func (o *PostLTENetworkIDSubscribersDefault) Error() string {
	return fmt.Sprintf("[POST /lte/{network_id}/subscribers][%d] PostLTENetworkIDSubscribers default  %+v", o._statusCode, o.Payload)
}

func (o *PostLTENetworkIDSubscribersDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostLTENetworkIDSubscribersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}