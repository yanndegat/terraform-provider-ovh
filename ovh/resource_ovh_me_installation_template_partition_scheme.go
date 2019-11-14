package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

func resourceMeInstallationTemplatePartitionScheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeInstallationTemplatePartitionSchemeCreate,
		Read:   resourceMeInstallationTemplatePartitionSchemeRead,
		Update: resourceMeInstallationTemplatePartitionSchemeUpdate,
		Delete: resourceMeInstallationTemplatePartitionSchemeDelete,

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This template name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of this partitioning scheme",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "on a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications)",
			},
		},
	}
}

func resourceMeInstallationTemplatePartitionSchemeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)

	opts := (&PartitionSchemeCreateOrUpdateOpts{}).FromResource(d)
	endpoint := fmt.Sprintf("/me/installationTemplate/%s/partitionScheme", templateName)

	if err := config.OVHClient.Post(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", templateName, opts.Name))

	return resourceMeInstallationTemplatePartitionSchemeRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)

	opts := (&PartitionSchemeCreateOrUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(templateName),
		url.PathEscape(opts.Name),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Calling PUT %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	return resourceMeInstallationTemplatePartitionSchemeRead(d, meta)
}

func resourceMeInstallationTemplatePartitionSchemeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	templateName := d.Get("template_name").(string)
	name := d.Get("name").(string)

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(templateName),
		url.PathEscape(name),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("Calling DELETE %s: %s \n", endpoint, err.Error())
	}

	return nil
}

func resourceMeInstallationTemplatePartitionSchemeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	templateName := d.Get("template_name").(string)
	name := d.Get("name").(string)

	scheme, err := getPartitionScheme(templateName, name, config.OVHClient)
	if err != nil {
		return err
	}

	// set resource attributes
	for k, v := range scheme.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s-%s", templateName, name))
	return nil
}

func getPartitionScheme(template, scheme string, client *ovh.Client) (*PartitionScheme, error) {
	r := &PartitionScheme{}

	endpoint := fmt.Sprintf(
		"/me/installationTemplate/%s/partitionScheme/%s",
		url.PathEscape(template),
		url.PathEscape(scheme),
	)

	if err := client.Get(endpoint, &r); err != nil {
		return nil, fmt.Errorf("Error calling %s: %s \n", endpoint, err.Error())
	}

	return r, nil
}
