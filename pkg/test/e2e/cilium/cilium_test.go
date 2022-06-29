//go:build e2e

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

package cilium

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cilium/cilium/api/v1/observer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clusterv1alpha1 "github.com/kubermatic/machine-controller/pkg/apis/cluster/v1alpha1"
	awstypes "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/aws/types"
	providerconfig "github.com/kubermatic/machine-controller/pkg/providerconfig/types"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/cluster/client"
	"k8c.io/kubermatic/v2/pkg/log"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/provider/kubernetes"
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/semver"
	"k8c.io/kubermatic/v2/pkg/test/e2e/utils"
	"k8c.io/kubermatic/v2/pkg/util/wait"
	yamlutil "k8c.io/kubermatic/v2/pkg/util/yaml"
	"k8c.io/kubermatic/v2/pkg/version/kubermatic"
	"k8c.io/operating-system-manager/pkg/providerconfig/ubuntu"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	userconfig      string
	accessKeyID     string
	secretAccessKey string
	logOptions      = log.NewDefaultOptions()
	namespace       = "kubermatic"
)

const (
	projectName  = "cilium-test-project"
	ciliumTestNs = "cilium-test"
)

func init() {
	flag.StringVar(&userconfig, "userconfig", "", "path to kubeconfig of usercluster")
	flag.StringVar(&namespace, "namespace", namespace, "namespace where KKP is installed into")
	logOptions.AddFlags(flag.CommandLine)
}

func TestInExistingCluster(t *testing.T) {
	if userconfig == "" {
		t.Logf("kubeconfig for usercluster not provided, test passes vacuously.")
		t.Logf("to run against an existing usercluster use following command:")
		t.Logf("go test ./pkg/test/e2e/cilium -v -tags e2e -timeout 30m -run TestInExistingCluster -userconfig <USERCLUSTER KUBECONFIG>")
		return
	}

	logger := log.NewFromOptions(logOptions).Sugar()

	config, err := clientcmd.BuildConfigFromFlags("", userconfig)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	client, err := ctrlruntimeclient.New(config, ctrlruntimeclient.Options{})
	if err != nil {
		t.Fatalf("failed to build ctrlruntime client: %v", err)
	}

	testUserCluster(context.Background(), t, logger, client)
}

func TestCiliumClusters(t *testing.T) {
	logger := log.NewFromOptions(logOptions).Sugar()

	accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		t.Fatalf("AWS_ACCESS_KEY_ID not set")
	}

	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		t.Fatalf("AWS_SECRET_ACCESS_KEY not set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	client, err := ctrlruntimeclient.New(config, ctrlruntimeclient.Options{})
	if err != nil {
		t.Fatalf("failed to build ctrlruntime client: %v", err)
	}

	tests := []struct {
		name      string
		proxyMode string
	}{
		{
			name:      "ebpf proxy mode test",
			proxyMode: resources.EBPFProxyMode,
		},
		{
			name:      "ipvs proxy mode test",
			proxyMode: resources.IPVSProxyMode,
		},
		{
			name:      "iptables proxy mode test",
			proxyMode: resources.IPTablesProxyMode,
		},
	}

	for _, test := range tests {
		proxyMode := test.proxyMode
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			client, cleanup, tLogger, err := createUserCluster(ctx, t, logger.With("proxymode", proxyMode), client, proxyMode)
			if cleanup != nil {
				defer cleanup()
			}

			if err != nil {
				t.Fatalf("failed to create user cluster: %v", err)
			}

			testUserCluster(ctx, t, tLogger, client)
		})
	}
}

//gocyclo:ignore
func testUserCluster(ctx context.Context, t *testing.T, log *zap.SugaredLogger, client ctrlruntimeclient.Client) {
	log.Info("Waiting for nodes to come up...")
	if err := checkNodeReadiness(ctx, t, log, client); err != nil {
		t.Fatalf("nodes never became ready: %v", err)
	}

	log.Info("Waiting for pods to get ready...")
	err := waitForPods(ctx, t, log, client, "kube-system", "k8s-app", []string{
		"cilium-operator",
		"cilium",
	})
	if err != nil {
		t.Fatalf("pods never became ready: %v", err)
	}

	log.Info("Running Cilium connectivity tests...")
	ns := corev1.Namespace{}
	ns.Name = ciliumTestNs
	err = client.Create(ctx, &ns)
	if err != nil {
		t.Fatalf("failed to create %q namespace: %v", ciliumTestNs, err)
	}
	defer func() {
		err := client.Delete(ctx, &ns)
		if err != nil {
			t.Fatalf("failed to delete %q namespace: %v", ciliumTestNs, err)
		}
	}()

	log = log.With("namespace", ciliumTestNs)
	log.Debug("Namespace created")

	installCiliumConnectivityTests(ctx, t, log, client)

	log.Info("Deploying hubble-relay-nodeport and hubble-ui-nodeport services...")
	cleanup := deployHubbleServices(ctx, t, log, client)
	defer cleanup()

	log.Info("Waiting for Cilium connectivity pods to get ready...")
	err = waitForPods(ctx, t, log, client, ciliumTestNs, "name", []string{
		"echo-a",
		"echo-b",
		"echo-b-headless",
		"echo-b-host",
		"echo-b-host-headless",
		"host-to-b-multi-node-clusterip",
		"host-to-b-multi-node-headless",
		"pod-to-a",
		"pod-to-a-allowed-cnp",
		"pod-to-a-denied-cnp",
		"pod-to-b-intra-node-nodeport",
		"pod-to-b-multi-node-clusterip",
		"pod-to-b-multi-node-headless",
		"pod-to-b-multi-node-nodeport",
		"pod-to-external-1111",
		"pod-to-external-fqdn-allow-google-cnp",
	})
	if err != nil {
		t.Fatalf("pods never became ready: %v", err)
	}

	log.Info("Checking for Hubble pods...")
	err = waitForPods(ctx, t, log, client, "kube-system", "k8s-app", []string{
		"hubble-relay",
		"hubble-ui",
	})
	if err != nil {
		t.Fatalf("pods never became ready: %v", err)
	}

	nodeIP, err := getAnyNodeIP(ctx, client)
	if err != nil {
		t.Fatalf("Nodes are ready, but could not get an IP: %v", err)
	}

	log.Info("Testing Hubble relay observe...")
	err = wait.PollLog(log, 2*time.Second, 5*time.Minute, func() (error, error) {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", nodeIP, 30077), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to dial to Hubble relay: %w", err), nil
		}
		defer conn.Close()

		nFlows := 20
		flowsClient, err := observer.NewObserverClient(conn).
			GetFlows(ctx, &observer.GetFlowsRequest{Number: uint64(nFlows)})
		if err != nil {
			return fmt.Errorf("failed to get flow client: %w", err), nil
		}

		for c := 0; c < nFlows; c++ {
			_, err := flowsClient.Recv()
			if err != nil {
				return fmt.Errorf("failed to get flow: %w", err), nil
			}
			// fmt.Println(flow)
		}

		return nil, nil
	})
	if err != nil {
		t.Fatalf("Hubble relay observe test failed: %v", err)
	}

	log.Info("Testing Hubble UI observe...")
	err = wait.PollLog(log, 2*time.Second, 5*time.Minute, func() (error, error) {
		uiURL := fmt.Sprintf("http://%s", net.JoinHostPort(nodeIP, "30007"))
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uiURL, nil)
		if err != nil {
			return fmt.Errorf("failed to construct request to Hubble UI: %w", err), nil
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get response from Hubble UI: %w", err), nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected HTTP 200 OK, got HTTP %d", resp.StatusCode), nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err), nil
		}

		if !strings.Contains(string(body), "Hubble") {
			return errors.New("failed to find Hubble in the body"), nil
		}

		return nil, nil
	})
	if err != nil {
		t.Fatalf("Hubble UI observe test failed: %v", err)
	}
}

func waitForPods(ctx context.Context, t *testing.T, log *zap.SugaredLogger, client ctrlruntimeclient.Client, namespace string, key string, names []string) error {
	log = log.With("namespace", namespace)

	r, err := labels.NewRequirement(key, selection.In, names)
	if err != nil {
		return fmt.Errorf("failed to build requirement: %w", err)
	}
	l := labels.NewSelector().Add(*r)

	return wait.PollLog(log, 5*time.Second, 5*time.Minute, func() (error, error) {
		pods := corev1.PodList{}
		err = client.List(ctx, &pods, ctrlruntimeclient.InNamespace(namespace), ctrlruntimeclient.MatchingLabelsSelector{Selector: l})
		if err != nil {
			return fmt.Errorf("failed to list Pods: %w", err), nil
		}

		if len(pods.Items) == 0 {
			return errors.New("no Pods found"), nil
		}

		unready := sets.NewString()
		for _, pod := range pods.Items {
			ready := false
			for _, c := range pod.Status.Conditions {
				if c.Type == corev1.ContainersReady {
					ready = c.Status == corev1.ConditionTrue
				}
			}

			if !ready {
				unready.Insert(pod.Name)
			}
		}

		if unready.Len() > 0 {
			return fmt.Errorf("not all Pods are ready: %v", unready.List()), nil
		}

		return nil, nil
	})
}

func deployHubbleServices(ctx context.Context, t *testing.T, log *zap.SugaredLogger, client ctrlruntimeclient.Client) func() {
	hubbleRelaySvc, err := resourcesFromYaml("./testdata/hubble-relay-svc.yaml")
	if err != nil {
		t.Fatalf("failed to read objects from yaml: %v", err)
	}

	hubbleUISvc, err := resourcesFromYaml("./testdata/hubble-ui-svc.yaml")
	if err != nil {
		t.Fatalf("failed to read objects from yaml: %v", err)
	}

	var cleanups []func()

	objects := []ctrlruntimeclient.Object{}
	objects = append(objects, hubbleRelaySvc...)
	objects = append(objects, hubbleUISvc...)

	cleanups = append(cleanups, func() {
		for _, object := range objects {
			err := client.Delete(ctx, object)
			if err != nil {
				log.Errorw("Failed to delete resource", zap.Error(err))
			}
		}
	})

	for _, object := range objects {
		err := client.Create(ctx, object)
		if err != nil {
			t.Fatalf("Failed to apply resource: %v", err)
		}

		log.Debugw("Created object", "kind", object.GetObjectKind(), "name", object.GetName())
	}

	return func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}

func installCiliumConnectivityTests(ctx context.Context, t *testing.T, log *zap.SugaredLogger, client ctrlruntimeclient.Client) {
	objs, err := resourcesFromYaml("./testdata/connectivity-check.yaml")
	if err != nil {
		t.Fatalf("failed to read objects from yaml: %v", err)
	}

	for _, obj := range objs {
		obj.SetNamespace(ciliumTestNs)
		if err := client.Create(ctx, obj); err != nil {
			t.Fatalf("failed to apply resource: %v", err)
		}

		log.Debugw("Created object", "kind", obj.GetObjectKind(), "name", obj.GetName())
	}
}

func getAnyNodeIP(ctx context.Context, client ctrlruntimeclient.Client) (string, error) {
	nodeList := corev1.NodeList{}
	if err := client.List(ctx, &nodeList); err != nil {
		return "", fmt.Errorf("failed to get nodes list: %w", err)
	}

	if len(nodeList.Items) == 0 {
		return "", errors.New("cluster has no nodes")
	}

	for _, node := range nodeList.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeExternalIP {
				return addr.Address, nil
			}
		}
	}

	return "", errors.New("no node has an ExternalIP")
}

func checkNodeReadiness(ctx context.Context, t *testing.T, log *zap.SugaredLogger, client ctrlruntimeclient.Client) error {
	expectedNodes := 2

	return wait.PollLog(log, 10*time.Second, 15*time.Minute, func() (error, error) {
		nodeList := corev1.NodeList{}
		err := client.List(ctx, &nodeList)
		if err != nil {
			return fmt.Errorf("failed to list nodes: %w", err), nil
		}

		if len(nodeList.Items) != expectedNodes {
			return fmt.Errorf("cluster has %d of %d nodes", len(nodeList.Items), expectedNodes), nil
		}

		readyNodeCount := 0
		for _, node := range nodeList.Items {
			for _, c := range node.Status.Conditions {
				if c.Type == corev1.NodeReady {
					readyNodeCount++
				}
			}
		}

		if readyNodeCount != expectedNodes {
			return fmt.Errorf("%d of %d nodes are ready", readyNodeCount, expectedNodes), nil
		}

		return nil, nil
	})
}

func resourcesFromYaml(filename string) ([]ctrlruntimeclient.Object, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	manifests, err := yamlutil.ParseMultipleDocuments(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var objs []ctrlruntimeclient.Object
	for _, m := range manifests {
		obj := &unstructured.Unstructured{}
		if err := kyaml.NewYAMLOrJSONDecoder(bytes.NewReader(m.Raw), 1024).Decode(obj); err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

// creates a usercluster on aws.
func createUserCluster(
	ctx context.Context,
	t *testing.T,
	log *zap.SugaredLogger,
	masterClient ctrlruntimeclient.Client,
	proxyMode string,
) (ctrlruntimeclient.Client, func(), *zap.SugaredLogger, error) {
	var teardowns []func()
	cleanup := func() {
		n := len(teardowns)
		for i := range teardowns {
			teardowns[n-1-i]()
		}
	}

	configGetter, err := provider.DynamicKubermaticConfigurationGetterFactory(masterClient, namespace)
	if err != nil {
		return nil, nil, log, fmt.Errorf("failed to create configGetter: %w", err)
	}

	// prepare helpers
	projectProvider, _ := kubernetes.NewProjectProvider(nil, masterClient)
	addonProvider := kubernetes.NewAddonProvider(masterClient, nil, configGetter)

	userClusterConnectionProvider, err := client.NewExternal(masterClient)
	if err != nil {
		return nil, nil, log, fmt.Errorf("failed to create userClusterConnectionProvider: %w", err)
	}

	clusterProvider := kubernetes.NewClusterProvider(
		nil,
		nil,
		userClusterConnectionProvider,
		"",
		nil,
		masterClient,
		nil,
		false,
		kubermatic.Versions{},
		nil,
	)

	log.Info("Creating project...")
	project, err := projectProvider.New(ctx, projectName, nil)
	if err != nil {
		return nil, nil, log, err
	}
	log = log.With("project", project.Name)
	teardowns = append(teardowns, func() {
		log.Info("Deleting project...")
		if err := masterClient.Delete(ctx, project); err != nil {
			t.Errorf("failed to delete project: %v", err)
		}
	})

	version := utils.KubernetesVersion()

	// create a usercluster on AWS
	cluster := &kubermaticv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("cilium-e2e-%s-", proxyMode),
			Labels: map[string]string{
				kubermaticv1.ProjectIDLabelKey: project.Name,
			},
		},
		Spec: kubermaticv1.ClusterSpec{
			HumanReadableName: fmt.Sprintf("Cilium %s e2e test cluster", proxyMode),
			Version:           *semver.NewSemverOrDie(version),
			KubernetesDashboard: kubermaticv1.KubernetesDashboard{
				Enabled: pointer.Bool(false),
			},
			EnableUserSSHKeyAgent: pointer.Bool(false),
			Cloud: kubermaticv1.CloudSpec{
				DatacenterName: "aws-eu-central-1a",
				AWS: &kubermaticv1.AWSCloudSpec{
					SecretAccessKey: secretAccessKey,
					AccessKeyID:     accessKeyID,
				},
			},
			CNIPlugin: &kubermaticv1.CNIPluginSettings{
				Type:    kubermaticv1.CNIPluginTypeCilium,
				Version: "v1.11",
			},
			ClusterNetwork: kubermaticv1.ClusterNetworkingConfig{
				ProxyMode:           proxyMode,
				KonnectivityEnabled: pointer.Bool(true),
			},
		},
	}

	log.Info("Creating cluster...")
	cluster, err = clusterProvider.NewUnsecured(ctx, project, cluster, "cilium@e2e.test")
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to create cluster: %w", err)
	}
	log = log.With("cluster", cluster.Name)
	teardowns = append(teardowns, func() {
		// This deletion will happen in the background, i.e. we are not waiting
		// for its completion. This is fine in e2e tests, where the surrounding
		// bash script will (as part of its normal cleanup) delete (and wait) all
		// userclusters anyway.
		log.Info("Deleting cluster...")
		if err := masterClient.Delete(ctx, cluster); err != nil {
			t.Errorf("failed to delete cluster: %v", err)
		}
	})

	// wait for cluster to be up and running
	log.Info("Waiting for cluster to become healthy...")
	err = wait.Poll(2*time.Second, 10*time.Minute, func() (error, error) {
		curCluster := kubermaticv1.Cluster{}
		if err := masterClient.Get(ctx, ctrlruntimeclient.ObjectKeyFromObject(cluster), &curCluster); err != nil {
			return fmt.Errorf("failed to retrieve cluster: %w", err), nil
		}

		if !curCluster.Status.ExtendedHealth.AllHealthy() {
			return errors.New("cluster is not all healthy"), nil
		}

		return nil, nil
	})
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("cluster did not become healthy: %w", err)
	}

	// update our local cluster variable with the newly reconciled address values
	if err := masterClient.Get(ctx, ctrlruntimeclient.ObjectKeyFromObject(cluster), cluster); err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to retrieve cluster: %w", err)
	}

	// create hubble addon
	log.Info("Installing hubble addon...")
	if _, err = addonProvider.NewUnsecured(ctx, cluster, "hubble", nil, nil); err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, cleanup, log, fmt.Errorf("failed to create addon: %w", err)
	}

	// retrieve usercluster kubeconfig, this can fail a couple of times until
	// the exposing mechanism is ready
	log.Info("Retrieving cluster client...")
	var clusterClient ctrlruntimeclient.Client
	err = wait.Poll(1*time.Second, 30*time.Second, func() (transient error, terminal error) {
		clusterClient, transient = clusterProvider.GetAdminClientForCustomerCluster(ctx, cluster)
		return transient, nil
	})
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("cluster did not become available: %w", err)
	}

	utilruntime.Must(clusterv1alpha1.AddToScheme(clusterClient.Scheme()))

	// prepare MachineDeployment
	encodedOSSpec, err := json.Marshal(ubuntu.Config{})
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to encode osspec: %w", err)
	}

	encodedCloudProviderSpec, err := json.Marshal(awstypes.RawConfig{
		InstanceType:     providerconfig.ConfigVarString{Value: "t3.small"},
		DiskType:         providerconfig.ConfigVarString{Value: "standard"},
		DiskSize:         int64(25),
		VpcID:            providerconfig.ConfigVarString{Value: cluster.Spec.Cloud.AWS.VPCID},
		InstanceProfile:  providerconfig.ConfigVarString{Value: cluster.Spec.Cloud.AWS.InstanceProfileName},
		Region:           providerconfig.ConfigVarString{Value: "eu-central-1"},
		AvailabilityZone: providerconfig.ConfigVarString{Value: "eu-central-1a"},
		SecurityGroupIDs: []providerconfig.ConfigVarString{{
			Value: cluster.Spec.Cloud.AWS.SecurityGroupID,
		}},
		Tags: map[string]string{
			"kubernetes.io/cluster/" + cluster.Name: "",
			"system/cluster":                        cluster.Name,
		},
	})
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to encode providerspec: %w", err)
	}

	cfg := providerconfig.Config{
		CloudProvider: providerconfig.CloudProviderAWS,
		CloudProviderSpec: runtime.RawExtension{
			Raw: encodedCloudProviderSpec,
		},
		OperatingSystem: providerconfig.OperatingSystemUbuntu,
		OperatingSystemSpec: runtime.RawExtension{
			Raw: encodedOSSpec,
		},
	}

	encodedConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to encode providerconfig: %w", err)
	}

	labels := map[string]string{
		"type": "worker",
	}

	md := clusterv1alpha1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-nodes",
			Namespace: "kube-system",
		},
		Spec: clusterv1alpha1.MachineDeploymentSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: pointer.Int32(2),
			Template: clusterv1alpha1.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: clusterv1alpha1.MachineSpec{
					Versions: clusterv1alpha1.MachineVersionInfo{
						Kubelet: version,
					},
					ProviderSpec: clusterv1alpha1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: encodedConfig,
						},
					},
				},
			},
		},
	}

	// create MachineDeployment
	log.Info("Creating MachineDeployment...")
	err = wait.PollImmediate(1*time.Second, 30*time.Second, func() (error, error) {
		return clusterClient.Create(ctx, &md), nil
	})
	if err != nil {
		return nil, cleanup, log, fmt.Errorf("failed to create MachineDeployment: %w", err)
	}

	return clusterClient, cleanup, log, nil
}
