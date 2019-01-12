// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/gitpods/gitpods/internal/api/v1/models"
)

// CreateRepositoryOKCode is the HTTP code returned for type CreateRepositoryOK
const CreateRepositoryOKCode int = 200

/*CreateRepositoryOK The repository has been created and is returned to you

swagger:response createRepositoryOK
*/
type CreateRepositoryOK struct {

	/*
	  In: Body
	*/
	Payload *models.Repository `json:"body,omitempty"`
}

// NewCreateRepositoryOK creates CreateRepositoryOK with default headers values
func NewCreateRepositoryOK() *CreateRepositoryOK {

	return &CreateRepositoryOK{}
}

// WithPayload adds the payload to the create repository o k response
func (o *CreateRepositoryOK) WithPayload(payload *models.Repository) *CreateRepositoryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create repository o k response
func (o *CreateRepositoryOK) SetPayload(payload *models.Repository) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRepositoryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateRepositoryUnprocessableEntityCode is the HTTP code returned for type CreateRepositoryUnprocessableEntity
const CreateRepositoryUnprocessableEntityCode int = 422

/*CreateRepositoryUnprocessableEntity The new repository has not been created due to invalid input

swagger:response createRepositoryUnprocessableEntity
*/
type CreateRepositoryUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.ValidationError `json:"body,omitempty"`
}

// NewCreateRepositoryUnprocessableEntity creates CreateRepositoryUnprocessableEntity with default headers values
func NewCreateRepositoryUnprocessableEntity() *CreateRepositoryUnprocessableEntity {

	return &CreateRepositoryUnprocessableEntity{}
}

// WithPayload adds the payload to the create repository unprocessable entity response
func (o *CreateRepositoryUnprocessableEntity) WithPayload(payload *models.ValidationError) *CreateRepositoryUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create repository unprocessable entity response
func (o *CreateRepositoryUnprocessableEntity) SetPayload(payload *models.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRepositoryUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateRepositoryDefault unexpected error

swagger:response createRepositoryDefault
*/
type CreateRepositoryDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateRepositoryDefault creates CreateRepositoryDefault with default headers values
func NewCreateRepositoryDefault(code int) *CreateRepositoryDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateRepositoryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create repository default response
func (o *CreateRepositoryDefault) WithStatusCode(code int) *CreateRepositoryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create repository default response
func (o *CreateRepositoryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create repository default response
func (o *CreateRepositoryDefault) WithPayload(payload *models.Error) *CreateRepositoryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create repository default response
func (o *CreateRepositoryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRepositoryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
