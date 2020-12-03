package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudDBEnterpriseClusterStatus string
type CloudDBEnterpriseClusterSecurityGroupStatus string
type CloudDBEnterpriseClusterSecurityGroupRuleStatus string

const (
	CloudDBEnterpriseClusterBaseUrl                  = "/cloudDB/enterprise/cluster/%s"
	CloudDBEnterpriseClusterSecurityGroupBaseUrl     = CloudDBEnterpriseClusterBaseUrl + "/securityGroup"
	CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl = CloudDBEnterpriseClusterSecurityGroupBaseUrl + "/%s/rule"

	CloudDBEnterpriseClusterStatusCreated    CloudDBEnterpriseClusterStatus = "created"
	CloudDBEnterpriseClusterStatusCreating   CloudDBEnterpriseClusterStatus = "creating"
	CloudDBEnterpriseClusterStatusDeleting   CloudDBEnterpriseClusterStatus = "deleting"
	CloudDBEnterpriseClusterStatusReopening  CloudDBEnterpriseClusterStatus = "reopening"
	CloudDBEnterpriseClusterStatusRestarting CloudDBEnterpriseClusterStatus = "restarting"
	CloudDBEnterpriseClusterStatusScaling    CloudDBEnterpriseClusterStatus = "scaling"
	CloudDBEnterpriseClusterStatusSuspended  CloudDBEnterpriseClusterStatus = "suspended"
	CloudDBEnterpriseClusterStatusSuspending CloudDBEnterpriseClusterStatus = "suspending"
	CloudDBEnterpriseClusterStatusUpdating   CloudDBEnterpriseClusterStatus = "updating"

	CloudDBEnterpriseClusterSecurityGroupStatusCreated  CloudDBEnterpriseClusterSecurityGroupStatus = "created"
	CloudDBEnterpriseClusterSecurityGroupStatusCreating CloudDBEnterpriseClusterSecurityGroupStatus = "creating"
	CloudDBEnterpriseClusterSecurityGroupStatusDeleting CloudDBEnterpriseClusterSecurityGroupStatus = "deleting"
	CloudDBEnterpriseClusterSecurityGroupStatusDeleted  CloudDBEnterpriseClusterSecurityGroupStatus = "deleted"
	CloudDBEnterpriseClusterSecurityGroupStatusUpdated  CloudDBEnterpriseClusterSecurityGroupStatus = "updated"
	CloudDBEnterpriseClusterSecurityGroupStatusUpdating CloudDBEnterpriseClusterSecurityGroupStatus = "updating"

	CloudDBEnterpriseClusterSecurityGroupRuleStatusCreated  CloudDBEnterpriseClusterSecurityGroupRuleStatus = "created"
	CloudDBEnterpriseClusterSecurityGroupRuleStatusCreating CloudDBEnterpriseClusterSecurityGroupRuleStatus = "creating"
	CloudDBEnterpriseClusterSecurityGroupRuleStatusDeleting CloudDBEnterpriseClusterSecurityGroupRuleStatus = "deleting"
	CloudDBEnterpriseClusterSecurityGroupRuleStatusDeleted  CloudDBEnterpriseClusterSecurityGroupRuleStatus = "deleted"
	CloudDBEnterpriseClusterSecurityGroupRuleStatusUpdated  CloudDBEnterpriseClusterSecurityGroupRuleStatus = "updated"
	CloudDBEnterpriseClusterSecurityGroupRuleStatusUpdating CloudDBEnterpriseClusterSecurityGroupRuleStatus = "updating"
)

type CloudDBEnterpriseCluster struct {
	Id         string                         `json:"id"`
	Status     CloudDBEnterpriseClusterStatus `json:"status"`
	RegionName string                         `json:"regionName"`
}

type CloudDBEnterpriseClusterSecurityGroupCreateUpdateOpts struct {
	Name      string `json:"name"`
	ClusterId string `json:"clusterId"`
}

type CloudDBEnterpriseClusterSecurityGroup struct {
	Id     string                         `json:"id"`
	Name   string                         `json:"name"`
	Status CloudDBEnterpriseClusterStatus `json:"status"`
	TaskId string                         `json:"taskId"`
}

func (opts *CloudDBEnterpriseClusterSecurityGroupCreateUpdateOpts) FromResource(d *schema.ResourceData) *CloudDBEnterpriseClusterSecurityGroupCreateUpdateOpts {
	name := helpers.GetNilStringPointerFromData(d, "name")
	opts.Name = *name
	clusterId := helpers.GetNilStringPointerFromData(d, "cluster_id")
	opts.ClusterId = *clusterId
	return opts
}

type CloudDBEnterpriseClusterSecurityGroupRuleCreateUpdateOpts struct {
	Source string `json:"source"`
}

type CloudDBEnterpriseClusterSecurityGroupRule struct {
	Id     string                                          `json:"id"`
	Source string                                          `json:"source"`
	Status CloudDBEnterpriseClusterSecurityGroupRuleStatus `json:"status"`
	TaskId string                                          `json:"taskId"`
}

func (opts *CloudDBEnterpriseClusterSecurityGroupRuleCreateUpdateOpts) FromResource(d *schema.ResourceData) *CloudDBEnterpriseClusterSecurityGroupRuleCreateUpdateOpts {
	source := helpers.GetNilStringPointerFromData(d, "source")
	opts.Source = *source
	return opts
}
