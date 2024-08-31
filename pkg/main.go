package pkg

import (
	"github.com/pkg/errors"
	elasticsearchv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/elasticsearch/elasticsearch/v1"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/kuberneteslabelkeys"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
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

	createdOperatorNamespace, err := kubernetescorev1.NewNamespace(ctx, "elastic-system",
		&kubernetescorev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name: pulumi.String("elastic-system"),
					//Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		}, pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create namespace")
	}

	//create helm-release
	createdEckOperatorHelmRelease, err := helm.NewRelease(ctx, "eck-operator",
		&helm.ReleaseArgs{
			Name:            pulumi.String("eck-operator"),
			Namespace:       createdOperatorNamespace.Metadata.Name(),
			Chart:           pulumi.String("eck-operator"),
			Version:         pulumi.String("2.14.0"),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"configKubernetes": pulumi.Map{
					"inherited_labels": pulumi.ToStringArray(
						[]string{
							kuberneteslabelkeys.Resource,
							kuberneteslabelkeys.Organization,
							kuberneteslabelkeys.Environment,
							kuberneteslabelkeys.ResourceKind,
							kuberneteslabelkeys.ResourceId,
						},
					),
				},
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://helm.elastic.co"),
			},
		}, pulumi.Parent(createdOperatorNamespace),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create helm release")
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

	nodeSets := elasticsearchv1.ElasticsearchSpecNodeSetsArray{}
	nodeSets = append(nodeSets, elasticsearchv1.ElasticsearchSpecNodeSetsArgs{
		Config: pulumi.Map{
			"node.roles":            pulumi.ToStringArray([]string{"master"}),
			"node.store.allow_mmap": pulumi.Bool(false),
		},
		Count: pulumi.Int(1),
		Name:  pulumi.String("master-nodes"),
		PodTemplate: pulumi.Map{
			"metadata": pulumi.Map{
				"labels": pulumi.StringMap{
					"role": pulumi.String("master"),
				},
			},
			"spec": pulumi.Map{
				"containers": pulumi.Array{
					pulumi.Map{
						"name": pulumi.String("elasticsearch"),
						"resources": pulumi.Map{
							"requests": pulumi.Map{
								"memory": pulumi.String("2Gi"),
								"cpu":    pulumi.String("1"),
							},
							"limits": pulumi.Map{
								"memory": pulumi.String("4Gi"),
								"cpu":    pulumi.String("2"),
							},
						},
					},
				},
			},
		},
	})

	nodeSets = append(nodeSets, elasticsearchv1.ElasticsearchSpecNodeSetsArgs{
		Name:  pulumi.String("data-nodes"),
		Count: pulumi.Int(1),
		Config: pulumi.Map{
			"node.roles":            pulumi.ToStringArray([]string{"data"}),
			"node.store.allow_mmap": pulumi.Bool(false),
		},
		PodTemplate: pulumi.Map{
			"metadata": pulumi.Map{
				"labels": pulumi.StringMap{
					"role": pulumi.String("master"),
				},
			},
			"spec": pulumi.Map{
				"containers": pulumi.Array{
					pulumi.Map{
						"name": pulumi.String("elasticsearch"),
						"resources": pulumi.Map{
							"requests": pulumi.Map{
								"memory": pulumi.String("4Gi"),
								"cpu":    pulumi.String("2"),
							},
							"limits": pulumi.Map{
								"memory": pulumi.String("8Gi"),
								"cpu":    pulumi.String("4"),
							},
						},
					},
				},
			},
		},
		VolumeClaimTemplates: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesArray{
			elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesArgs{
				Metadata: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesMetadataArgs{
					Name: pulumi.String("elasticsearch-data"),
				},
				Spec: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesSpecArgs{
					AccessModes: pulumi.StringArray{
						pulumi.String("ReadWriteOnce"),
					},
					Resources: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesSpecResourcesArgs{
						Requests: pulumi.Map{
							"storage": pulumi.String("2Gi"),
						},
					},
				},
			},
		},
	})

	elasticsearchv1.NewElasticsearch(ctx, locals.ElasticsearchKubernetes.Metadata.Id, &elasticsearchv1.ElasticsearchArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.ElasticsearchKubernetes.Metadata.Id),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.Labels),
		},
		Spec: &elasticsearchv1.ElasticsearchSpecArgs{
			NodeSets: nodeSets,
			Version:  pulumi.String("8.13.0"),
		},
	}, pulumi.Parent(createdNamespace), pulumi.DependsOn([]pulumi.Resource{createdEckOperatorHelmRelease}))

	return nil
}
