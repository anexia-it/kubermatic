/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package main

import (
	apimodels "k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

func getOSNameFromSpec(spec apimodels.OperatingSystemSpec) string {
	if spec.Centos != nil {
		return "centos"
	}
	if spec.Ubuntu != nil {
		return "ubuntu"
	}
	if spec.Sles != nil {
		return "sles"
	}
	if spec.Rhel != nil {
		return "rhel"
	}
	if spec.Flatcar != nil {
		return "flatcar"
	}

	return ""
}
