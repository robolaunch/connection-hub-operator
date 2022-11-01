//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerCredentials) DeepCopyInto(out *BrokerCredentials) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerCredentials.
func (in *BrokerCredentials) DeepCopy() *BrokerCredentials {
	if in == nil {
		return nil
	}
	out := new(BrokerCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerStatus) DeepCopyInto(out *BrokerStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerStatus.
func (in *BrokerStatus) DeepCopy() *BrokerStatus {
	if in == nil {
		return nil
	}
	out := new(BrokerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChartResourceStatus) DeepCopyInto(out *ChartResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChartResourceStatus.
func (in *ChartResourceStatus) DeepCopy() *ChartResourceStatus {
	if in == nil {
		return nil
	}
	out := new(ChartResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChartStatus) DeepCopyInto(out *ChartStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChartStatus.
func (in *ChartStatus) DeepCopy() *ChartStatus {
	if in == nil {
		return nil
	}
	out := new(ChartStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInstance) DeepCopyInto(out *CloudInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInstance.
func (in *CloudInstance) DeepCopy() *CloudInstance {
	if in == nil {
		return nil
	}
	out := new(CloudInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInstanceList) DeepCopyInto(out *CloudInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CloudInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInstanceList.
func (in *CloudInstanceList) DeepCopy() *CloudInstanceList {
	if in == nil {
		return nil
	}
	out := new(CloudInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInstanceSpec) DeepCopyInto(out *CloudInstanceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInstanceSpec.
func (in *CloudInstanceSpec) DeepCopy() *CloudInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(CloudInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInstanceStatus) DeepCopyInto(out *CloudInstanceStatus) {
	*out = *in
	out.DeployerStatus = in.DeployerStatus
	out.ConnectionResources = in.ConnectionResources
	out.GatewayConnection = in.GatewayConnection
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInstanceStatus.
func (in *CloudInstanceStatus) DeepCopy() *CloudInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(CloudInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionResourceStatus) DeepCopyInto(out *ConnectionResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionResourceStatus.
func (in *ConnectionResourceStatus) DeepCopy() *ConnectionResourceStatus {
	if in == nil {
		return nil
	}
	out := new(ConnectionResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionResourceStatuses) DeepCopyInto(out *ConnectionResourceStatuses) {
	*out = *in
	out.ClusterStatus = in.ClusterStatus
	out.EndpointStatus = in.EndpointStatus
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionResourceStatuses.
func (in *ConnectionResourceStatuses) DeepCopy() *ConnectionResourceStatuses {
	if in == nil {
		return nil
	}
	out := new(ConnectionResourceStatuses)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceStatus) DeepCopyInto(out *CustomResourceStatus) {
	*out = *in
	out.OwnedResourceStatus = in.OwnedResourceStatus
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceStatus.
func (in *CustomResourceStatus) DeepCopy() *CustomResourceStatus {
	if in == nil {
		return nil
	}
	out := new(CustomResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeployerStatus) DeepCopyInto(out *DeployerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeployerStatus.
func (in *DeployerStatus) DeepCopy() *DeployerStatus {
	if in == nil {
		return nil
	}
	out := new(DeployerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayConnection) DeepCopyInto(out *GatewayConnection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayConnection.
func (in *GatewayConnection) DeepCopy() *GatewayConnection {
	if in == nil {
		return nil
	}
	out := new(GatewayConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmChartProperties) DeepCopyInto(out *HelmChartProperties) {
	*out = *in
	out.Repository = in.Repository
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmChartProperties.
func (in *HelmChartProperties) DeepCopy() *HelmChartProperties {
	if in == nil {
		return nil
	}
	out := new(HelmChartProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmRepository) DeepCopyInto(out *HelmRepository) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmRepository.
func (in *HelmRepository) DeepCopy() *HelmRepository {
	if in == nil {
		return nil
	}
	out := new(HelmRepository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *K8sNodeInfo) DeepCopyInto(out *K8sNodeInfo) {
	*out = *in
	if in.Selectors != nil {
		in, out := &in.Selectors, &out.Selectors
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new K8sNodeInfo.
func (in *K8sNodeInfo) DeepCopy() *K8sNodeInfo {
	if in == nil {
		return nil
	}
	out := new(K8sNodeInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceStatus) DeepCopyInto(out *NamespaceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceStatus.
func (in *NamespaceStatus) DeepCopy() *NamespaceStatus {
	if in == nil {
		return nil
	}
	out := new(NamespaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorStatus) DeepCopyInto(out *OperatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorStatus.
func (in *OperatorStatus) DeepCopy() *OperatorStatus {
	if in == nil {
		return nil
	}
	out := new(OperatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnedResourceStatus) DeepCopyInto(out *OwnedResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnedResourceStatus.
func (in *OwnedResourceStatus) DeepCopy() *OwnedResourceStatus {
	if in == nil {
		return nil
	}
	out := new(OwnedResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhysicalInstance) DeepCopyInto(out *PhysicalInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhysicalInstance.
func (in *PhysicalInstance) DeepCopy() *PhysicalInstance {
	if in == nil {
		return nil
	}
	out := new(PhysicalInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PhysicalInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhysicalInstanceList) DeepCopyInto(out *PhysicalInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PhysicalInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhysicalInstanceList.
func (in *PhysicalInstanceList) DeepCopy() *PhysicalInstanceList {
	if in == nil {
		return nil
	}
	out := new(PhysicalInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PhysicalInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhysicalInstanceSpec) DeepCopyInto(out *PhysicalInstanceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhysicalInstanceSpec.
func (in *PhysicalInstanceSpec) DeepCopy() *PhysicalInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(PhysicalInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhysicalInstanceStatus) DeepCopyInto(out *PhysicalInstanceStatus) {
	*out = *in
	out.DeployerStatus = in.DeployerStatus
	out.ConnectionResources = in.ConnectionResources
	out.GatewayConnection = in.GatewayConnection
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhysicalInstanceStatus.
func (in *PhysicalInstanceStatus) DeepCopy() *PhysicalInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(PhysicalInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceItem) DeepCopyInto(out *ResourceItem) {
	*out = *in
	out.ObjectKey = in.ObjectKey
	out.GroupVersionKind = in.GroupVersionKind
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceItem.
func (in *ResourceItem) DeepCopy() *ResourceItem {
	if in == nil {
		return nil
	}
	out := new(ResourceItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Submariner) DeepCopyInto(out *Submariner) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Submariner.
func (in *Submariner) DeepCopy() *Submariner {
	if in == nil {
		return nil
	}
	out := new(Submariner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Submariner) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerBroker) DeepCopyInto(out *SubmarinerBroker) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerBroker.
func (in *SubmarinerBroker) DeepCopy() *SubmarinerBroker {
	if in == nil {
		return nil
	}
	out := new(SubmarinerBroker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubmarinerBroker) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerBrokerList) DeepCopyInto(out *SubmarinerBrokerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SubmarinerBroker, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerBrokerList.
func (in *SubmarinerBrokerList) DeepCopy() *SubmarinerBrokerList {
	if in == nil {
		return nil
	}
	out := new(SubmarinerBrokerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubmarinerBrokerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerBrokerSpec) DeepCopyInto(out *SubmarinerBrokerSpec) {
	*out = *in
	out.Helm = in.Helm
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerBrokerSpec.
func (in *SubmarinerBrokerSpec) DeepCopy() *SubmarinerBrokerSpec {
	if in == nil {
		return nil
	}
	out := new(SubmarinerBrokerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerBrokerStatus) DeepCopyInto(out *SubmarinerBrokerStatus) {
	*out = *in
	out.NamespaceStatus = in.NamespaceStatus
	out.ChartStatus = in.ChartStatus
	out.ChartResourceStatus = in.ChartResourceStatus
	in.NodeInfo.DeepCopyInto(&out.NodeInfo)
	out.BrokerCredentials = in.BrokerCredentials
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerBrokerStatus.
func (in *SubmarinerBrokerStatus) DeepCopy() *SubmarinerBrokerStatus {
	if in == nil {
		return nil
	}
	out := new(SubmarinerBrokerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerList) DeepCopyInto(out *SubmarinerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Submariner, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerList.
func (in *SubmarinerList) DeepCopy() *SubmarinerList {
	if in == nil {
		return nil
	}
	out := new(SubmarinerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubmarinerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerOperator) DeepCopyInto(out *SubmarinerOperator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerOperator.
func (in *SubmarinerOperator) DeepCopy() *SubmarinerOperator {
	if in == nil {
		return nil
	}
	out := new(SubmarinerOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubmarinerOperator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerOperatorList) DeepCopyInto(out *SubmarinerOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SubmarinerOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerOperatorList.
func (in *SubmarinerOperatorList) DeepCopy() *SubmarinerOperatorList {
	if in == nil {
		return nil
	}
	out := new(SubmarinerOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubmarinerOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerOperatorSpec) DeepCopyInto(out *SubmarinerOperatorSpec) {
	*out = *in
	out.BrokerCredentials = in.BrokerCredentials
	out.Helm = in.Helm
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerOperatorSpec.
func (in *SubmarinerOperatorSpec) DeepCopy() *SubmarinerOperatorSpec {
	if in == nil {
		return nil
	}
	out := new(SubmarinerOperatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerOperatorStatus) DeepCopyInto(out *SubmarinerOperatorStatus) {
	*out = *in
	out.NamespaceStatus = in.NamespaceStatus
	out.ChartStatus = in.ChartStatus
	out.ChartResourceStatus = in.ChartResourceStatus
	in.NodeInfo.DeepCopyInto(&out.NodeInfo)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerOperatorStatus.
func (in *SubmarinerOperatorStatus) DeepCopy() *SubmarinerOperatorStatus {
	if in == nil {
		return nil
	}
	out := new(SubmarinerOperatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerSpec) DeepCopyInto(out *SubmarinerSpec) {
	*out = *in
	out.BrokerCredentials = in.BrokerCredentials
	out.BrokerHelmChart = in.BrokerHelmChart
	out.OperatorHelmChart = in.OperatorHelmChart
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerSpec.
func (in *SubmarinerSpec) DeepCopy() *SubmarinerSpec {
	if in == nil {
		return nil
	}
	out := new(SubmarinerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubmarinerStatus) DeepCopyInto(out *SubmarinerStatus) {
	*out = *in
	in.NodeInfo.DeepCopyInto(&out.NodeInfo)
	in.BrokerStatus.DeepCopyInto(&out.BrokerStatus)
	out.OperatorStatus = in.OperatorStatus
	out.CustomResourceStatus = in.CustomResourceStatus
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubmarinerStatus.
func (in *SubmarinerStatus) DeepCopy() *SubmarinerStatus {
	if in == nil {
		return nil
	}
	out := new(SubmarinerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tenancy) DeepCopyInto(out *Tenancy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tenancy.
func (in *Tenancy) DeepCopy() *Tenancy {
	if in == nil {
		return nil
	}
	out := new(Tenancy)
	in.DeepCopyInto(out)
	return out
}
