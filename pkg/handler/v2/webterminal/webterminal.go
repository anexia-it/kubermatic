/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webterminal

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"k8c.io/kubermatic/v2/pkg/handler/auth"
	handlercommon "k8c.io/kubermatic/v2/pkg/handler/common"
	"k8c.io/kubermatic/v2/pkg/handler/v1/common"
	"k8c.io/kubermatic/v2/pkg/provider"
)

func CreateOIDCKubeconfigSecretEndpoint(projectProvider provider.ProjectProvider, privilegedProjectProvider provider.PrivilegedProjectProvider, oidcIssuerVerifier auth.OIDCIssuerVerifier, oidcCfg common.OIDCConfiguration) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(handlercommon.CreateOIDCKubeconfigReq)
		return handlercommon.CreateOIDCKubeconfigSecretEndpoint(ctx, projectProvider, privilegedProjectProvider, oidcIssuerVerifier, oidcCfg, req)
	}
}

func DecodeCreateOIDCKubeconfig(c context.Context, r *http.Request) (interface{}, error) {
	return handlercommon.DecodeCreateOIDCKubeconfig(c, r)
}

func EncodeOIDCKubeconfig(c context.Context, w http.ResponseWriter, response interface{}) (err error) {
	return handlercommon.EncodeOIDCKubeconfigSecret(c, w, response)
}
