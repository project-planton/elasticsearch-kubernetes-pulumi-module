package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/elasticsearch-kubernetes-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *elasticsearchkubernetes.ElasticsearchKubernetesStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.Input)
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	createdNamespace, err := kubernetescorev1.NewNamespace(ctx, locals.ElasticsearchKubernetes.Metadata.Id,
		&kubernetescorev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(locals.ElasticsearchKubernetes.Metadata.Id),
					Labels: pulumi.ToStringMap(s.Labels),
				}),
		}, pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create namespace")
	}

	//export name of the namespace
	ctx.Export(outputs.Namespace, createdNamespace.Metadata.Name())

	if err := elasticsearch(ctx, locals, createdNamespace, s.Labels); err != nil {
		return errors.Wrap(err, "failed to create elastic search resources")
	}

	if locals.ElasticsearchKubernetes.Spec.Ingress.IsEnabled {
		if err := ingress(ctx, locals, createdNamespace, kubernetesProvider, s.Labels); err != nil {
			return errors.Wrap(err, "failed to create ingress resources")
		}
	}

	return nil
}
