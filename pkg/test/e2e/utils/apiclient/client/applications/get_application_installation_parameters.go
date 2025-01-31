// Code generated by go-swagger; DO NOT EDIT.

package applications

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

// NewGetApplicationInstallationParams creates a new GetApplicationInstallationParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetApplicationInstallationParams() *GetApplicationInstallationParams {
	return &GetApplicationInstallationParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetApplicationInstallationParamsWithTimeout creates a new GetApplicationInstallationParams object
// with the ability to set a timeout on a request.
func NewGetApplicationInstallationParamsWithTimeout(timeout time.Duration) *GetApplicationInstallationParams {
	return &GetApplicationInstallationParams{
		timeout: timeout,
	}
}

// NewGetApplicationInstallationParamsWithContext creates a new GetApplicationInstallationParams object
// with the ability to set a context for a request.
func NewGetApplicationInstallationParamsWithContext(ctx context.Context) *GetApplicationInstallationParams {
	return &GetApplicationInstallationParams{
		Context: ctx,
	}
}

// NewGetApplicationInstallationParamsWithHTTPClient creates a new GetApplicationInstallationParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetApplicationInstallationParamsWithHTTPClient(client *http.Client) *GetApplicationInstallationParams {
	return &GetApplicationInstallationParams{
		HTTPClient: client,
	}
}

/* GetApplicationInstallationParams contains all the parameters to send to the API endpoint
   for the get application installation operation.

   Typically these are written to a http.Request.
*/
type GetApplicationInstallationParams struct {

	// AppinstallName.
	ApplicationInstallationName string

	// ClusterID.
	ClusterID string

	// Namespace.
	Namespace string

	// ProjectID.
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get application installation params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetApplicationInstallationParams) WithDefaults() *GetApplicationInstallationParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get application installation params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetApplicationInstallationParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get application installation params
func (o *GetApplicationInstallationParams) WithTimeout(timeout time.Duration) *GetApplicationInstallationParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get application installation params
func (o *GetApplicationInstallationParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get application installation params
func (o *GetApplicationInstallationParams) WithContext(ctx context.Context) *GetApplicationInstallationParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get application installation params
func (o *GetApplicationInstallationParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get application installation params
func (o *GetApplicationInstallationParams) WithHTTPClient(client *http.Client) *GetApplicationInstallationParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get application installation params
func (o *GetApplicationInstallationParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApplicationInstallationName adds the appinstallName to the get application installation params
func (o *GetApplicationInstallationParams) WithApplicationInstallationName(appinstallName string) *GetApplicationInstallationParams {
	o.SetApplicationInstallationName(appinstallName)
	return o
}

// SetApplicationInstallationName adds the appinstallName to the get application installation params
func (o *GetApplicationInstallationParams) SetApplicationInstallationName(appinstallName string) {
	o.ApplicationInstallationName = appinstallName
}

// WithClusterID adds the clusterID to the get application installation params
func (o *GetApplicationInstallationParams) WithClusterID(clusterID string) *GetApplicationInstallationParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the get application installation params
func (o *GetApplicationInstallationParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithNamespace adds the namespace to the get application installation params
func (o *GetApplicationInstallationParams) WithNamespace(namespace string) *GetApplicationInstallationParams {
	o.SetNamespace(namespace)
	return o
}

// SetNamespace adds the namespace to the get application installation params
func (o *GetApplicationInstallationParams) SetNamespace(namespace string) {
	o.Namespace = namespace
}

// WithProjectID adds the projectID to the get application installation params
func (o *GetApplicationInstallationParams) WithProjectID(projectID string) *GetApplicationInstallationParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the get application installation params
func (o *GetApplicationInstallationParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *GetApplicationInstallationParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param appinstall_name
	if err := r.SetPathParam("appinstall_name", o.ApplicationInstallationName); err != nil {
		return err
	}

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param namespace
	if err := r.SetPathParam("namespace", o.Namespace); err != nil {
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
