// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetRepositoryHandlerFunc turns a function with the right signature into a get repository handler
type GetRepositoryHandlerFunc func(GetRepositoryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetRepositoryHandlerFunc) Handle(params GetRepositoryParams) middleware.Responder {
	return fn(params)
}

// GetRepositoryHandler interface for that can handle valid get repository params
type GetRepositoryHandler interface {
	Handle(GetRepositoryParams) middleware.Responder
}

// NewGetRepository creates a new http.Handler for the get repository operation
func NewGetRepository(ctx *middleware.Context, handler GetRepositoryHandler) *GetRepository {
	return &GetRepository{Context: ctx, Handler: handler}
}

/*GetRepository swagger:route GET /repositories/{owner}/{name} repositories getRepository

Get a repository by owner name and its name

*/
type GetRepository struct {
	Context *middleware.Context
	Handler GetRepositoryHandler
}

func (o *GetRepository) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetRepositoryParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
