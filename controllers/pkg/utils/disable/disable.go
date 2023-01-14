package fed_disable

import (
	"errors"
	"io"
	"strings"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/kubefed/pkg/apis/core/typeconfig"
	ctlutil "sigs.k8s.io/kubefed/pkg/controller/util"
	"sigs.k8s.io/kubefed/pkg/kubefedctl"
	"sigs.k8s.io/kubefed/pkg/kubefedctl/enable"
	"sigs.k8s.io/kubefed/pkg/kubefedctl/options"
)

type disableType struct {
	options.GlobalSubcommandOptions
	options.CommonEnableOptions
	disableTypeOptions
}

type disableTypeOptions struct {
	deleteCRD           bool
	enableTypeDirective *enable.EnableTypeDirective
}

func DisableKindInFederation(federationOperator *connectionhubv1alpha1.FederationOperator, args []string, cmdOut io.Writer, hostConfig *rest.Config) error {

	opts := &disableType{}
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

// Complete ensures that options are valid and marshals them if necessary.
func (j *disableType) Complete(args []string) error {
	j.disableTypeOptions.enableTypeDirective = enable.NewEnableTypeDirective()
	j.FederatedGroup = options.DefaultFederatedGroup
	directive := j.disableTypeOptions.enableTypeDirective

	if err := j.SetName(args); err != nil {
		return err
	}

	if !j.deleteCRD {
		if len(j.TargetVersion) > 0 {
			return errors.New("--version flag valid only with --delete-crd")
		} else if j.FederatedGroup != options.DefaultFederatedGroup {
			return errors.New("--kubefed-group flag valid only with --delete-crd")
		}
	}

	if len(j.TargetVersion) > 0 {
		directive.Spec.TargetVersion = j.TargetVersion
	}
	if len(j.FederatedGroup) > 0 {
		directive.Spec.FederatedGroup = j.FederatedGroup
	}

	return nil
}

// Run is the implementation of the `disable` command.
func (j *disableType) Run(cmdOut io.Writer, hostConfig *rest.Config) error {

	name := j.TargetName
	if !strings.Contains(j.TargetName, ".") {
		apiResource, err := enable.LookupAPIResource(hostConfig, j.TargetName, "")
		if err != nil {
			return err
		}
		name = typeconfig.GroupQualifiedName(*apiResource)
	}

	typeConfigName := ctlutil.QualifiedName{
		Namespace: j.KubeFedNamespace,
		Name:      name,
	}
	j.disableTypeOptions.enableTypeDirective.Name = typeConfigName.Name
	return kubefedctl.DisableFederation(cmdOut, hostConfig, j.disableTypeOptions.enableTypeDirective, typeConfigName, j.deleteCRD, j.DryRun, true)
}
