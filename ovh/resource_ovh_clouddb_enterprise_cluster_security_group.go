package ovh

import (
	"context"
	"fmt"
	"github.com/ovh/go-ovh/ovh"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudDBEnterpriseClusterSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudDBEnterpriseClusterSecurityGroupCreateOrUpdate,
		Read:   resourceCloudDBEnterpriseClusterSecurityGroupRead,
		Delete: resourceCloudDBEnterpriseClusterSecurityGroupDelete,
		Update: resourceCloudDBEnterpriseClusterSecurityGroupCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudDBEnterpriseClusterSecurityGroupImportState,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
		},
	}
}

func resourceCloudDBEnterpriseClusterSecurityGroupCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sg := (&CloudDBEnterpriseClusterSecurityGroupCreateUpdateOpts{}).FromResource(d)

	clusterId := d.Get("cluster_id").(string)

	getUrl := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl, clusterId)

	var securityGroup CloudDBEnterpriseClusterSecurityGroup

	if d.Id() != "" {
		getUrl = fmt.Sprintf("%s/%s", getUrl, d.Id())
		securityGroup.Id = d.Id()
	}

	// Retry Post/Put action until accepted.
	// there are cases where a background operation may be ongoing at the time
	// of the request.
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		if d.Id() != "" {
			err = config.OVHClient.Put(getUrl, sg, nil)
		} else {
			err = config.OVHClient.Post(getUrl, sg, &securityGroup)
		}
		if err != nil {
			apiError := err.(*ovh.APIError)
			// Cluster is pending
			if apiError.Code == http.StatusForbidden {
				return resource.RetryableError(err)
			}

			id := apiError.QueryID
			return resource.NonRetryableError(
				fmt.Errorf("Error calling DELETE (id: %s) %s:\n\t%q", id, getUrl, err))
		}
		// Successful delete
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error execution DB Entreprise Security Group Create/Update:\n\t %q", err)
	}

	getUrl = fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", clusterId, securityGroup.Id)

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusCreated),
		},
		Pending: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusCreating),
		},
		Refresh: func() (interface{}, string, error) {
			stateResp := CloudDBEnterpriseClusterSecurityGroup{}
			if err := config.OVHClient.Get(getUrl, &stateResp); err != nil {
				return nil, "", err
			}
			return d, string(stateResp.Status), nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group creation:\n\t %q", err)
	}

	d.SetId(securityGroup.Id)

	return resourceCloudDBEnterpriseClusterSecurityGroupRead(d, meta)
}

func resourceCloudDBEnterpriseClusterSecurityGroupImportState(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import id is not cluster_id/security_group_id formatted")
	}
	id := splitId[0]
	groupId := splitId[1]
	d.SetId(groupId)
	d.Set("cluster_id", id)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudDBEnterpriseClusterSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)
	url := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", clusterId, d.Id())
	resp := &CloudDBEnterpriseClusterSecurityGroup{}

	if err := config.OVHClient.Get(url, resp); err != nil {
		return helpers.CheckDeleted(d, err, url)
	}
	d.Set("name", resp.Name)
	return nil
}

func resourceCloudDBEnterpriseClusterSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)

	url := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", clusterId, d.Id())

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
		return fmt.Errorf("Error execution DB Entreprise Security Group Deletion:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", clusterId, d.Id())
	// monitor task	execution
	stateConf := &resource.StateChangeConf{
		Target: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusDeleted),
		},
		Pending: []string{
			string(CloudDBEnterpriseClusterSecurityGroupStatusDeleting),
		},
		Refresh: func() (interface{}, string, error) {
			stateResp := CloudDBEnterpriseClusterSecurityGroup{}
			if err := config.OVHClient.Get(getUrl, &stateResp); err != nil {
				if err.(*ovh.APIError).Code == 404 {
					return d, string(CloudDBEnterpriseClusterSecurityGroupStatusDeleted), nil
				}
				return nil, "", err
			}
			return d, string(stateResp.Status), nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
