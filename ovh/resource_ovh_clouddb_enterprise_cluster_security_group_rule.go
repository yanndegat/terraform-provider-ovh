package ovh

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudDBEnterpriseClusterSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudDBEnterpriseClusterSecurityGroupRuleCreate,
		Read:   resourceCloudDBEnterpriseClusterSecurityGroupRuleRead,
		Delete: resourceCloudDBEnterpriseClusterSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudDBEnterpriseClusterSecurityGroupRuleImportState,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceCloudDBEnterpriseClusterSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sg := (&CloudDBEnterpriseClusterSecurityGroupRuleCreateUpdateOpts{}).FromResource(d)
	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)

	url := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl, clusterId, securityGroupId)
	var securityGroupRule CloudDBEnterpriseClusterSecurityGroupRule

	// Retry Post/Put action until accepted.
	// there are cases where a background operation may be ongoing at the time
	// of the request.
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		if err := config.OVHClient.Post(url, sg, &securityGroupRule); err != nil {
			apiError := err.(*ovh.APIError)
			// Cluster is pending
			if apiError.Code == http.StatusForbidden {
				return resource.RetryableError(err)
			}

			id := apiError.QueryID
			return resource.NonRetryableError(fmt.Errorf("Error calling DELETE (id: %s) %s:\n\t%q", id, url, err))
		}
		// Successful delete
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group Rule creation:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, securityGroupRule.Id)

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{
			string(CloudDBEnterpriseClusterSecurityGroupRuleStatusCreated),
		},
		Pending: []string{
			string(CloudDBEnterpriseClusterSecurityGroupRuleStatusCreating),
		},
		Refresh: func() (interface{}, string, error) {
			var resp CloudDBEnterpriseClusterSecurityGroupRule
			if err := config.OVHClient.Get(getUrl, &resp); err != nil {
				return nil, "", err
			}
			return d, string(resp.Status), nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group creation:\n\t %q", err)
	}

	d.SetId(securityGroupRule.Id)

	return resourceCloudDBEnterpriseClusterSecurityGroupRuleRead(d, meta)
}

func resourceCloudDBEnterpriseClusterSecurityGroupRuleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("import id is not cluster_id/security_group_id/security_group_rule formatted")
	}
	d.Set("cluster_id", splitId[0])
	d.Set("security_group_id", splitId[1])
	d.SetId(splitId[2])

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudDBEnterpriseClusterSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	url := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())
	resp := &CloudDBEnterpriseClusterSecurityGroupRule{}

	if err := config.OVHClient.Get(url, resp); err != nil {
		return helpers.CheckDeleted(d, err, url)
	}

	d.Set("source", resp.Source)
	return nil
}

func resourceCloudDBEnterpriseClusterSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	url := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())

	// Retry Delete action until accepted.
	// there are cases where a background operation may be ongoing at the time
	// of the delete request.
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		if err := config.OVHClient.Delete(url, nil); err != nil {
			apiError := err.(*ovh.APIError)
			// Cluster is pending
			if apiError.Code == http.StatusForbidden {
				return resource.RetryableError(err)
			}
			id := apiError.QueryID
			return resource.NonRetryableError(
				fmt.Errorf("Error calling DELETE (id: %s) %s:\n\t%q", id, url, err))
		}
		// Successful delete
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error execution DB Entreprise Security Group Rule Deletion:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusDeleted),
		},
		Pending: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusDeleting),
		},
		Refresh: func() (interface{}, string, error) {
			stateResp := CloudDBEnterpriseClusterSecurityGroupRule{}
			if err := config.OVHClient.Get(getUrl, &stateResp); err != nil {
				if err.(*ovh.APIError).Code == 404 {
					return d, string(CloudDBEnterpriseClusterSecurityGroupRuleStatusDeleted), nil
				}
				return nil, "", err
			}
			return d, string(stateResp.Status), nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
