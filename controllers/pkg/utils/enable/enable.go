package fed_enable

import (
	"context"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/apis/core/typeconfig"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	ctlutil "sigs.k8s.io/kubefed/pkg/controller/util"
	"sigs.k8s.io/kubefed/pkg/kubefedctl/enable"
	"sigs.k8s.io/kubefed/pkg/kubefedctl/options"
	"sigs.k8s.io/kubefed/pkg/kubefedctl/util"
)

// ---------------------------------------------------------------
// Adding this module (from kubefed) is necessary to enable some types with `spec.statusCollection` true (FederatedTypeConfig)

type enableType struct {
	options.GlobalSubcommandOptions
	options.CommonEnableOptions
	enableTypeOptions
}

type enableTypeOptions struct {
	federatedVersion    string
	output              string
	outputYAML          bool
	filename            string
	enableTypeDirective *enable.EnableTypeDirective
}

func EnableKindInFederation(federationOperator *connectionhubv1alpha1.FederationOperator, args []string, cmdOut io.Writer, hostConfig *rest.Config) error {

	opts := &enableType{}
	opts.KubeFedNamespace = federationOperator.GetNamespaceMetadata().Name

	err := opts.Complete(args)
	if err != nil {
		return err
	}

	err = opts.Run(cmdOut, hostConfig)
	if err != nil {
		return err
	}

	return nil
}

func (j *enableType) Complete(args []string) error {
	j.enableTypeOptions.enableTypeDirective = enable.NewEnableTypeDirective()
	fd := j.enableTypeOptions.enableTypeDirective

	if j.output == "yaml" {
		j.enableTypeOptions.outputYAML = true
	} else if len(j.output) > 0 {
		return errors.Errorf("Invalid value for --output: %s", j.output)
	}

	if len(j.filename) > 0 {
		err := enable.DecodeYAMLFromFile(j.filename, fd)
		if err != nil {
			return errors.Wrapf(err, "Failed to load yaml from file %q", j.filename)
		}
		return nil
	}

	if err := j.SetName(args); err != nil {
		return err
	}

	fd.Name = j.TargetName

	if len(j.TargetVersion) > 0 {
		fd.Spec.TargetVersion = j.TargetVersion
	}
	if len(j.FederatedGroup) > 0 {
		fd.Spec.FederatedGroup = j.FederatedGroup
	}
	if len(j.federatedVersion) > 0 {
		fd.Spec.FederatedVersion = j.federatedVersion
	}

	return nil
}

// Run is the implementation of the `enable` command.
func (j *enableType) Run(cmdOut io.Writer, hostConfig *rest.Config) error {

	resourcesOld, err := enable.GetResources(hostConfig, j.enableTypeOptions.enableTypeDirective)
	if err != nil {
		return err
	}

	resources := typeResources(*resourcesOld)

	if j.enableTypeOptions.outputYAML {
		concreteTypeConfig := resources.TypeConfig.(*fedv1b1.FederatedTypeConfig)
		objects := []runtimeclient.Object{concreteTypeConfig, resources.CRD}
		err := writeObjectsToYAML(objects, cmdOut)
		if err != nil {
			return errors.Wrap(err, "Failed to write objects to YAML")
		}
		// -o yaml implies dry run
		return nil
	}

	return CreateResources(cmdOut, hostConfig, &resources, j.KubeFedNamespace, j.DryRun)
}

func writeObjectsToYAML(objects []runtimeclient.Object, w io.Writer) error {
	for _, obj := range objects {
		if _, err := w.Write([]byte("---\n")); err != nil {
			return errors.Wrap(err, "Error encoding object to yaml")
		}

		if err := writeObjectToYAML(obj, w); err != nil {
			return errors.Wrap(err, "Error encoding object to yaml")
		}
	}
	return nil
}

func writeObjectToYAML(obj runtimeclient.Object, w io.Writer) error {
	json, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return err
	}

	unstructuredObj := &unstructured.Unstructured{}
	if _, _, err := unstructured.UnstructuredJSONScheme.Decode(json, nil, unstructuredObj); err != nil {
		return err
	}

	return util.WriteUnstructuredToYaml(unstructuredObj, w)
}

// -------------------------------------

type typeResources struct {
	TypeConfig typeconfig.Interface
	CRD        *apiextv1.CustomResourceDefinition
}

func qualifiedAPIResourceName(resource metav1.APIResource) string {
	if resource.Group == "" {
		return fmt.Sprintf("%s/%s", resource.Name, resource.Version)
	}
	return fmt.Sprintf("%s.%s/%s", resource.Name, resource.Group, resource.Version)
}

func CreateResources(cmdOut io.Writer, config *rest.Config, resources *typeResources, namespace string, dryRun bool) error {
	write := func(data string) {
		if cmdOut != nil {
			if _, err := cmdOut.Write([]byte(data)); err != nil {
				klog.Fatalf("Unexpected err: %v\n", err)
			}
		}
	}

	hostClientset, err := util.HostClientset(config)
	if err != nil {
		return errors.Wrap(err, "Failed to create host clientset")
	}
	_, err = hostClientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return errors.Wrapf(err, "KubeFed system namespace %q does not exist", namespace)
	} else if err != nil {
		return errors.Wrapf(err, "Error attempting to determine whether KubeFed system namespace %q exists", namespace)
	}

	client, err := genericclient.New(config)
	if err != nil {
		return errors.Wrap(err, "Failed to get kubefed clientset")
	}

	concreteTypeConfig := resources.TypeConfig.(*fedv1b1.FederatedTypeConfig)
	existingTypeConfig := &fedv1b1.FederatedTypeConfig{}
	err = client.Get(context.TODO(), existingTypeConfig, namespace, concreteTypeConfig.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return errors.Wrapf(err, "Error retrieving FederatedTypeConfig %q", concreteTypeConfig.Name)
	}
	if err == nil {
		fedType := existingTypeConfig.GetFederatedType()
		target := existingTypeConfig.GetTargetType()
		concreteType := concreteTypeConfig.GetFederatedType()
		if fedType.Name != concreteType.Name || fedType.Version != concreteType.Version || fedType.Group != concreteType.Group {
			return errors.Errorf("Federation is already enabled for %q with federated type %q. Changing the federated type to %q is not supported.",
				qualifiedAPIResourceName(target),
				qualifiedAPIResourceName(fedType),
				qualifiedAPIResourceName(concreteType))
		}
	}

	crdClient, err := apiextv1client.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "Failed to create crd clientset")
	}

	existingCRD, err := crdClient.CustomResourceDefinitions().Get(context.Background(), resources.CRD.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		if !dryRun {
			_, err = crdClient.CustomResourceDefinitions().Create(context.Background(), resources.CRD, metav1.CreateOptions{})
			if err != nil {
				return errors.Wrapf(err, "Error creating CRD %q", resources.CRD.Name)
			}
		}
		write(fmt.Sprintf("customresourcedefinition.apiextensions.k8s.io/%s created\n", resources.CRD.Name))
	case err != nil:
		return errors.Wrapf(err, "Error getting CRD %q", resources.CRD.Name)
	default:
		ftcs := &fedv1b1.FederatedTypeConfigList{}
		err := client.List(context.TODO(), ftcs, namespace)
		if err != nil {
			return errors.Wrap(err, "Error getting FederatedTypeConfig list")
		}

		for _, ftc := range ftcs.Items {
			targetAPI := concreteTypeConfig.Spec.TargetType
			existingAPI := ftc.Spec.TargetType
			if IsEquivalentAPI(&existingAPI, &targetAPI) {
				existingName := qualifiedAPIResourceName(ftc.GetTargetType())
				name := qualifiedAPIResourceName(concreteTypeConfig.GetTargetType())
				qualifiedFTCName := ctlutil.QualifiedName{
					Namespace: ftc.Namespace,
					Name:      ftc.Name,
				}

				return errors.Errorf("Failed to enable %q. Federation of this type is already enabled for equivalent type %q by FederatedTypeConfig %q",
					name, existingName, qualifiedFTCName)
			}

			if concreteTypeConfig.Name == ftc.Name {
				continue
			}

			fedType := ftc.Spec.FederatedType
			name := typeconfig.GroupQualifiedName(metav1.APIResource{Name: fedType.PluralName, Group: fedType.Group})
			if name == existingCRD.Name {
				return errors.Errorf("Failed to enable federation of %q due to the FederatedTypeConfig for %q already referencing a federated type CRD named %q. If these target types are distinct despite sharing the same kind, specifying a non-default --federated-group should allow %q to be enabled.",
					concreteTypeConfig.Name, ftc.Name, name, concreteTypeConfig.Name)
			}
		}

		existingCRD.Spec = resources.CRD.Spec
		if !dryRun {
			_, err = crdClient.CustomResourceDefinitions().Update(context.Background(), existingCRD, metav1.UpdateOptions{})
			if err != nil {
				return errors.Wrapf(err, "Error updating CRD %q", resources.CRD.Name)
			}
		}
		write(fmt.Sprintf("customresourcedefinition.apiextensions.k8s.io/%s updated\n", resources.CRD.Name))
	}

	concreteTypeConfig.Namespace = namespace
	statusCollectionEnabled := fedv1b1.StatusCollectionEnabled
	concreteTypeConfig.Spec.StatusCollection = &statusCollectionEnabled
	err = client.Get(context.TODO(), existingTypeConfig, namespace, concreteTypeConfig.Name)
	createdOrUpdated := "created"
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrapf(err, "Error retrieving FederatedTypeConfig %q", concreteTypeConfig.Name)
		}
		if !dryRun {
			err = client.Create(context.TODO(), concreteTypeConfig)
			if err != nil {
				return errors.Wrapf(err, "Error creating FederatedTypeConfig %q", concreteTypeConfig.Name)
			}
		}
	} else {
		patch := runtimeclient.MergeFrom(existingTypeConfig.DeepCopy())
		existingTypeConfig.Spec = concreteTypeConfig.Spec
		if !dryRun {
			err := client.Patch(context.TODO(), existingTypeConfig, patch)
			if err != nil {
				return errors.Wrapf(err, "Error updating FederatedTypeConfig %q", concreteTypeConfig.Name)
			}
		}
		createdOrUpdated = "updated"
	}
	write(fmt.Sprintf("federatedtypeconfig.core.kubefed.io/%s %s in namespace %s\n",
		concreteTypeConfig.Name, createdOrUpdated, namespace))
	return nil
}

func IsEquivalentAPI(existingAPI, newAPI *fedv1b1.APIResource) bool {
	if existingAPI.PluralName != newAPI.PluralName {
		return false
	}

	apis, ok := equivalentAPIs[existingAPI.PluralName]
	if !ok {
		return false
	}

	for _, gv := range apis {
		if gv.Group == existingAPI.Group && gv.Version == existingAPI.Version {
			// skip exactly matched API from equivalent API list
			continue
		}

		if gv.Group == newAPI.Group && gv.Version == newAPI.Version {
			return true
		}
	}

	return false
}

var equivalentAPIs = map[string][]schema.GroupVersion{
	"deployments": {
		{
			Group:   "apps",
			Version: "v1",
		},
		{
			Group:   "apps",
			Version: "v1beta1",
		},
		{
			Group:   "apps",
			Version: "v1beta2",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
	"daemonsets": {
		{
			Group:   "apps",
			Version: "v1",
		},
		{
			Group:   "apps",
			Version: "v1beta1",
		},
		{
			Group:   "apps",
			Version: "v1beta2",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
	"statefulsets": {
		{
			Group:   "apps",
			Version: "v1",
		},
		{
			Group:   "apps",
			Version: "v1beta1",
		},
		{
			Group:   "apps",
			Version: "v1beta2",
		},
	},
	"replicasets": {
		{
			Group:   "apps",
			Version: "v1",
		},
		{
			Group:   "apps",
			Version: "v1beta1",
		},
		{
			Group:   "apps",
			Version: "v1beta2",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
	"networkpolicies": {
		{
			Group:   "networking.k8s.io",
			Version: "v1",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
	"podsecuritypolicies": {
		{
			Group:   "policy",
			Version: "v1beta1",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
	"ingresses": {
		{
			Group:   "networking.k8s.io",
			Version: "v1beta1",
		},
		{
			Group:   "extensions",
			Version: "v1beta1",
		},
	},
}
