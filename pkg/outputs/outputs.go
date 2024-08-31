package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	NamespaceOutputName               = "namespace"
	ServiceOutputName                 = "service"
	PortForwardCommandOutputName      = "port-forward-command"
	KubeEndpointOutputName            = "kube-endpoint"
	IngressExternalHostnameOutputName = "ingress-external-hostname"
	IngressInternalHostnameOutputName = "ingress-internal-hostname"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *elasticsearchkubernetes.ElasticsearchKubernetesStackInput) *elasticsearchkubernetes.ElasticsearchKubernetesStackOutputs {
	return &elasticsearchkubernetes.ElasticsearchKubernetesStackOutputs{
		Namespace:          autoapistackoutput.GetVal(pulumiOutputs, NamespaceOutputName),
		Service:            autoapistackoutput.GetVal(pulumiOutputs, ServiceOutputName),
		PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, PortForwardCommandOutputName),
		KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, KubeEndpointOutputName),
		ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostnameOutputName),
		InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostnameOutputName),
	}
}
