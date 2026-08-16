/*
Copyright 2026 The Kubermatic Kubernetes Platform contributors.

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

package cloudcontroller

import (
	"testing"

	kubermaticv1 "k8c.io/kubermatic/sdk/v2/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/resources"

	appsv1 "k8s.io/api/apps/v1"
)

func TestAnexiaDeploymentCCMVersion(t *testing.T) {
	tests := []struct {
		name           string
		defaultVersion string
		clusterVersion string
		expected       string
	}{
		{
			name:     "bundled default",
			expected: anexiaCCMVersion,
		},
		{
			name:           "configured default",
			defaultVersion: "1.6.0",
			expected:       "1.6.0",
		},
		{
			name:           "cluster override",
			defaultVersion: "1.6.0",
			clusterVersion: "1.7.0",
			expected:       "1.7.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &kubermaticv1.KubermaticConfiguration{}
			config.Spec.UserCluster.Anexia.CCMVersion = tt.defaultVersion

			cluster := &kubermaticv1.Cluster{}
			cluster.Spec.Cloud.Anexia = &kubermaticv1.AnexiaCloudSpec{CCMVersion: tt.clusterVersion}

			data := resources.NewTemplateDataBuilder().
				WithKubermaticConfiguration(config).
				WithCluster(cluster).
				Build()

			_, reconcile := anexiaDeploymentReconciler(data)()
			deployment, err := reconcile(&appsv1.Deployment{})
			if err != nil {
				t.Fatalf("failed to reconcile deployment: %v", err)
			}

			expectedImage := "anx-cr.io/anexia/anx-cloud-controller-manager:" + tt.expected
			if got := deployment.Spec.Template.Spec.Containers[0].Image; got != expectedImage {
				t.Fatalf("expected CCM image %q, got %q", expectedImage, got)
			}
		})
	}
}
