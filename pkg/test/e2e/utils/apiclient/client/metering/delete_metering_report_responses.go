// Code generated by go-swagger; DO NOT EDIT.

package metering

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// DeleteMeteringReportReader is a Reader for the DeleteMeteringReport structure.
type DeleteMeteringReportReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteMeteringReportReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteMeteringReportOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewDeleteMeteringReportUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewDeleteMeteringReportForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteMeteringReportDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteMeteringReportOK creates a DeleteMeteringReportOK with default headers values
func NewDeleteMeteringReportOK() *DeleteMeteringReportOK {
	return &DeleteMeteringReportOK{}
}

/* DeleteMeteringReportOK describes a response with status code 200, with default header values.

EmptyResponse is a empty response
*/
type DeleteMeteringReportOK struct {
}

func (o *DeleteMeteringReportOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/admin/metering/reports/{report_name}][%d] deleteMeteringReportOK ", 200)
}

func (o *DeleteMeteringReportOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteMeteringReportUnauthorized creates a DeleteMeteringReportUnauthorized with default headers values
func NewDeleteMeteringReportUnauthorized() *DeleteMeteringReportUnauthorized {
	return &DeleteMeteringReportUnauthorized{}
}

/* DeleteMeteringReportUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type DeleteMeteringReportUnauthorized struct {
}

func (o *DeleteMeteringReportUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/admin/metering/reports/{report_name}][%d] deleteMeteringReportUnauthorized ", 401)
}

func (o *DeleteMeteringReportUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteMeteringReportForbidden creates a DeleteMeteringReportForbidden with default headers values
func NewDeleteMeteringReportForbidden() *DeleteMeteringReportForbidden {
	return &DeleteMeteringReportForbidden{}
}

/* DeleteMeteringReportForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type DeleteMeteringReportForbidden struct {
}

func (o *DeleteMeteringReportForbidden) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/admin/metering/reports/{report_name}][%d] deleteMeteringReportForbidden ", 403)
}

func (o *DeleteMeteringReportForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteMeteringReportDefault creates a DeleteMeteringReportDefault with default headers values
func NewDeleteMeteringReportDefault(code int) *DeleteMeteringReportDefault {
	return &DeleteMeteringReportDefault{
		_statusCode: code,
	}
}

/* DeleteMeteringReportDefault describes a response with status code -1, with default header values.

errorResponse
*/
type DeleteMeteringReportDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete metering report default response
func (o *DeleteMeteringReportDefault) Code() int {
	return o._statusCode
}

func (o *DeleteMeteringReportDefault) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/admin/metering/reports/{report_name}][%d] deleteMeteringReport default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteMeteringReportDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteMeteringReportDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
