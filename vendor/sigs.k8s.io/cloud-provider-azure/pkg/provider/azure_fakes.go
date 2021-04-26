/*
Copyright 2020 The Kubernetes Authors.

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

package provider

import (
	"fmt"

	"github.com/golang/mock/gomock"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/cloud-provider-azure/pkg/auth"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/diskclient/mockdiskclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/interfaceclient/mockinterfaceclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/loadbalancerclient/mockloadbalancerclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/publicipclient/mockpublicipclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/routeclient/mockrouteclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/routetableclient/mockroutetableclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/securitygroupclient/mocksecuritygroupclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/subnetclient/mocksubnetclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/vmclient/mockvmclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/vmssclient/mockvmssclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/vmssvmclient/mockvmssvmclient"
)

var (
	errPreconditionFailedEtagMismatch = fmt.Errorf("PreconditionFailedEtagMismatch")
)

// GetTestCloud returns a fake azure cloud for unit tests in Azure related CSI drivers
func GetTestCloud(ctrl *gomock.Controller) (az *Cloud) {
	az = &Cloud{
		Config: Config{
			AzureAuthConfig: auth.AzureAuthConfig{
				TenantID:       "tenant",
				SubscriptionID: "subscription",
			},
			ResourceGroup:                "rg",
			VnetResourceGroup:            "rg",
			RouteTableResourceGroup:      "rg",
			SecurityGroupResourceGroup:   "rg",
			Location:                     "westus",
			VnetName:                     "vnet",
			SubnetName:                   "subnet",
			SecurityGroupName:            "nsg",
			RouteTableName:               "rt",
			PrimaryAvailabilitySetName:   "as",
			PrimaryScaleSetName:          "vmss",
			MaximumLoadBalancerRuleCount: 250,
			VMType:                       vmTypeStandard,
		},
		nodeZones:          map[string]sets.String{},
		nodeInformerSynced: func() bool { return true },
		nodeResourceGroups: map[string]string{},
		unmanagedNodes:     sets.NewString(),
		routeCIDRs:         map[string]string{},
		eventRecorder:      &record.FakeRecorder{},
	}
	az.DisksClient = mockdiskclient.NewMockInterface(ctrl)
	az.InterfacesClient = mockinterfaceclient.NewMockInterface(ctrl)
	az.LoadBalancerClient = mockloadbalancerclient.NewMockInterface(ctrl)
	az.PublicIPAddressesClient = mockpublicipclient.NewMockInterface(ctrl)
	az.RoutesClient = mockrouteclient.NewMockInterface(ctrl)
	az.RouteTablesClient = mockroutetableclient.NewMockInterface(ctrl)
	az.SecurityGroupsClient = mocksecuritygroupclient.NewMockInterface(ctrl)
	az.SubnetsClient = mocksubnetclient.NewMockInterface(ctrl)
	az.VirtualMachineScaleSetsClient = mockvmssclient.NewMockInterface(ctrl)
	az.VirtualMachineScaleSetVMsClient = mockvmssvmclient.NewMockInterface(ctrl)
	az.VirtualMachinesClient = mockvmclient.NewMockInterface(ctrl)
	az.VMSet = newAvailabilitySet(az)
	az.vmCache, _ = az.newVMCache()
	az.lbCache, _ = az.newLBCache()
	az.nsgCache, _ = az.newNSGCache()
	az.rtCache, _ = az.newRouteTableCache()

	common := &controllerCommon{cloud: az, resourceGroup: "rg", location: "westus"}
	az.controllerCommon = common
	az.ManagedDiskController = &ManagedDiskController{common: common}

	return az
}

// GetTestCloudWithExtendedLocation returns a fake azure cloud for unit tests in Azure related CSI drivers with extended location.
func GetTestCloudWithExtendedLocation(ctrl *gomock.Controller) (az *Cloud) {
	az = GetTestCloud(ctrl)
	az.Config.ExtendedLocationName = "microsoftlosangeles1"
	az.Config.ExtendedLocationType = "EdgeZone"
	return az
}
