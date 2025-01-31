// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteExternalClusterParams creates a new DeleteExternalClusterParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteExternalClusterParams() *DeleteExternalClusterParams {
	return &DeleteExternalClusterParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteExternalClusterParamsWithTimeout creates a new DeleteExternalClusterParams object
// with the ability to set a timeout on a request.
func NewDeleteExternalClusterParamsWithTimeout(timeout time.Duration) *DeleteExternalClusterParams {
	return &DeleteExternalClusterParams{
		timeout: timeout,
	}
}

// NewDeleteExternalClusterParamsWithContext creates a new DeleteExternalClusterParams object
// with the ability to set a context for a request.
func NewDeleteExternalClusterParamsWithContext(ctx context.Context) *DeleteExternalClusterParams {
	return &DeleteExternalClusterParams{
		Context: ctx,
	}
}

// NewDeleteExternalClusterParamsWithHTTPClient creates a new DeleteExternalClusterParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteExternalClusterParamsWithHTTPClient(client *http.Client) *DeleteExternalClusterParams {
	return &DeleteExternalClusterParams{
		HTTPClient: client,
	}
}

/* DeleteExternalClusterParams contains all the parameters to send to the API endpoint
   for the delete external cluster operation.

   Typically these are written to a http.Request.
*/
type DeleteExternalClusterParams struct {

	/* Action.

	     The Action is used to check if to `Delete` the cluster:
	both the actual cluter from the provider
	and the respective KKP cluster object
	By default the cluster will `Disconnect` which means the KKP cluster object will be deleted,
	cluster still exists on the provider, but is no longer connected/imported in KKP
	*/
	Action *string

	// ClusterID.
	ClusterID string

	// ProjectID.
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete external cluster params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteExternalClusterParams) WithDefaults() *DeleteExternalClusterParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete external cluster params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteExternalClusterParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete external cluster params
func (o *DeleteExternalClusterParams) WithTimeout(timeout time.Duration) *DeleteExternalClusterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete external cluster params
func (o *DeleteExternalClusterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete external cluster params
func (o *DeleteExternalClusterParams) WithContext(ctx context.Context) *DeleteExternalClusterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete external cluster params
func (o *DeleteExternalClusterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete external cluster params
func (o *DeleteExternalClusterParams) WithHTTPClient(client *http.Client) *DeleteExternalClusterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete external cluster params
func (o *DeleteExternalClusterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAction adds the action to the delete external cluster params
func (o *DeleteExternalClusterParams) WithAction(action *string) *DeleteExternalClusterParams {
	o.SetAction(action)
	return o
}

// SetAction adds the action to the delete external cluster params
func (o *DeleteExternalClusterParams) SetAction(action *string) {
	o.Action = action
}

// WithClusterID adds the clusterID to the delete external cluster params
func (o *DeleteExternalClusterParams) WithClusterID(clusterID string) *DeleteExternalClusterParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the delete external cluster params
func (o *DeleteExternalClusterParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithProjectID adds the projectID to the delete external cluster params
func (o *DeleteExternalClusterParams) WithProjectID(projectID string) *DeleteExternalClusterParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the delete external cluster params
func (o *DeleteExternalClusterParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteExternalClusterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Action != nil {

		// header param action
		if err := r.SetHeaderParam("action", *o.Action); err != nil {
			return err
		}
	}

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
