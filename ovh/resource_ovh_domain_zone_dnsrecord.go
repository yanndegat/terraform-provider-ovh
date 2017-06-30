package ovh

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func resourceDomainZoneDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainZoneDnsRecordCreate,
		Read:   resourceDomainZoneDnsRecordRead,
		Update: resourceDomainZoneDnsRecordUpdate,
		Delete: resourceDomainZoneDnsRecordDelete,

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"field_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceDomainZoneDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	zone := d.Get("zone_name").(string)

	// Create the new record
	record := &DomainZoneRecordOpts{
		FieldType: d.Get("field_type").(string),
		SubDomain: d.Get("sub_domain").(string),
		Target:    d.Get("target").(string),
		Ttl:       d.Get("ttl").(int),
	}

	log.Printf("[DEBUG] Record create configuration: %#v", record)

	response := &DomainZoneRecordResponse{}
	endpoint := fmt.Sprintf("/domain/zone/%s/record", zone)
	err := config.OVHClient.Post(endpoint, record, response)

	if err != nil {
		return fmt.Errorf("Failed to create Domain Record %s for zone %s: %s", record, zone, err)
	}

	log.Printf("[DEBUG] Domain Zone Record created with id: %d", response.Id)

	return readDomainZoneDnsRecord(d, response)
}

func resourceDomainZoneDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	zone := d.Get("zone_name").(string)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	response := &DomainZoneRecordResponse{}
	endpoint := fmt.Sprintf("/domain/zone/%s/record/%d", zone, recordId)
	err = config.OVHClient.Get(endpoint, response)
	if err != nil {
		return err
	}

	return readDomainZoneDnsRecord(d, response)
}

func resourceDomainZoneDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zone := d.Get("zone_name").(string)
	recordId, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	// Create the new record
	record := &DomainZoneRecordOpts{
		FieldType: d.Get("field_type").(string),
		SubDomain: d.Get("sub_domain").(string),
		Target:    d.Get("target").(string),
		Ttl:       d.Get("ttl").(int),
	}

	log.Printf("[DEBUG] Record create configuration: %#v", record)

	response := &DomainZoneRecordResponse{}
	endpoint := fmt.Sprintf("/domain/zone/%s/record/%d", zone, recordId)
	err = config.OVHClient.Put(endpoint, record, response)

	if err != nil {
		return fmt.Errorf("Failed to update Domain Record %s for zone %s: %s", record, zone, err)
	}

	log.Printf("[DEBUG] Domain Zone Record %d updated", response.Id)

	return readDomainZoneDnsRecord(d, response)
}

func resourceDomainZoneDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zone := d.Get("zone_name").(string)
	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	log.Printf("[DEBUG] Will delete Domain Record %d for zone %s", recordId, zone)
	endpoint := fmt.Sprintf("/domain/zone/%s/record/%d", zone, recordId)
	err = config.OVHClient.Delete(endpoint, nil)

	if err != nil {
		return fmt.Errorf("Error deleting Domain Record %d: %s", recordId, err)
	}

	d.SetId("")

	return nil
}

func domainZoneDnsRecordExists(zone string, id string, config *Config) (bool, error) {
	recordId, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("Error converting Record ID: %s", err)
	}

	response := &DomainZoneRecordResponse{}
	endpoint := fmt.Sprintf("/domain/zone/%s/record/%d", zone, recordId)
	err = config.OVHClient.Get(endpoint, response)
	if err != nil {
		switch typed := err.(type) {
		case *ovh.APIError:
			if typed.Code == 404 {
				return false, nil
			} else {
				return false, err
			}

		default:
			return false, err
		}
	}

	return true, nil
}

func readDomainZoneDnsRecord(d *schema.ResourceData, r *DomainZoneRecordResponse) error {
	d.Set("zone_name", r.ZoneName)
	d.Set("ttl", r.Ttl)
	d.Set("field_type", r.FieldType)
	d.Set("sub_domain", r.SubDomain)
	d.Set("target", r.Target)
	d.SetId(fmt.Sprintf("%d", r.Id))
	return nil
}
