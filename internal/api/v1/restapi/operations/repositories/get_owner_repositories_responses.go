// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/gitpods/gitpods/internal/api/v1/models"
)

// GetOwnerRepositoriesOKCode is the HTTP code returned for type GetOwnerRepositoriesOK
const GetOwnerRepositoriesOKCode int = 200

/*GetOwnerRepositoriesOK The repositories found by its owner name

swagger:response getOwnerRepositoriesOK
*/
type GetOwnerRepositoriesOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Repository `json:"body,omitempty"`
}

// NewGetOwnerRepositoriesOK creates GetOwnerRepositoriesOK with default headers values
func NewGetOwnerRepositoriesOK() *GetOwnerRepositoriesOK {

	return &GetOwnerRepositoriesOK{}
}

// WithPayload adds the payload to the get owner repositories o k response
func (o *GetOwnerRepositoriesOK) WithPayload(payload []*models.Repository) *GetOwnerRepositoriesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get owner repositories o k response
func (o *GetOwnerRepositoriesOK) SetPayload(payload []*models.Repository) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOwnerRepositoriesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*models.Repository, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetOwnerRepositoriesNotFoundCode is the HTTP code returned for type GetOwnerRepositoriesNotFound
const GetOwnerRepositoriesNotFoundCode int = 404

/*GetOwnerRepositoriesNotFound The owner could not be found by this username

swagger:response getOwnerRepositoriesNotFound
*/
type GetOwnerRepositoriesNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOwnerRepositoriesNotFound creates GetOwnerRepositoriesNotFound with default headers values
func NewGetOwnerRepositoriesNotFound() *GetOwnerRepositoriesNotFound {

	return &GetOwnerRepositoriesNotFound{}
}

// WithPayload adds the payload to the get owner repositories not found response
func (o *GetOwnerRepositoriesNotFound) WithPayload(payload *models.Error) *GetOwnerRepositoriesNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get owner repositories not found response
func (o *GetOwnerRepositoriesNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOwnerRepositoriesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetOwnerRepositoriesDefault unexpected error

swagger:response getOwnerRepositoriesDefault
*/
type GetOwnerRepositoriesDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOwnerRepositoriesDefault creates GetOwnerRepositoriesDefault with default headers values
func NewGetOwnerRepositoriesDefault(code int) *GetOwnerRepositoriesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetOwnerRepositoriesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get owner repositories default response
func (o *GetOwnerRepositoriesDefault) WithStatusCode(code int) *GetOwnerRepositoriesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get owner repositories default response
func (o *GetOwnerRepositoriesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get owner repositories default response
func (o *GetOwnerRepositoriesDefault) WithPayload(payload *models.Error) *GetOwnerRepositoriesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get owner repositories default response
func (o *GetOwnerRepositoriesDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOwnerRepositoriesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
