package dome9

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/dome9/dome9-sdk-go/services/cloudaccounts"
	"github.com/dome9/dome9-sdk-go/services/cloudaccounts/azure"
)

func resourceCloudAccountAzure() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountAzureCreate,
		Read:   resourceCloudAccountAzureRead,
		Update: resourceCloudAccountAzureUpdate,
		Delete: resourceCloudAccountAzureDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operation_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizational_unit_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudAccountAzureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	req := expandCloudAccountAzureRequest(d)
	log.Printf("[INFO] Creating Azure Cloud Account with request %+v\n", req)

	resp, _, err := client.cloudaccountAzure.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created Azure CloudAccount. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceCloudAccountAzureRead(d, meta)
}

func resourceCloudAccountAzureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	getCloudAccountQueryParams := cloudaccounts.QueryParameters{ID: d.Id()}
	resp, _, err := client.cloudaccountAzure.Get(&getCloudAccountQueryParams)
	if err != nil {
		return nil
	}

	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("subscription_id", resp.SubscriptionID)
	_ = d.Set("tenant_id", resp.TenantID)
	_ = d.Set("credentials", resp.Credentials)
	_ = d.Set("operation_mode", resp.OperationMode)
	_ = d.Set("vendor", resp.Vendor)
	_ = d.Set("creation_date", resp.CreationDate)
	_ = d.Set("organizational_unit_id", resp.OrganizationalUnitID)
	_ = d.Set("organizational_unit_path", resp.OrganizationalUnitPath)
	_ = d.Set("organizational_unit_name", resp.OrganizationalUnitName)

	return nil
}

func resourceCloudAccountAzureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	log.Printf("[INFO] Deleting Azure Cloud Account ID: %v\n", d.Id())
	if _, err := client.cloudaccountAzure.Delete(d.Id()); err != nil {
		return err
	}

	return nil
}

func resourceCloudAccountAzureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	log.Println("An updated occurred")

	if d.HasChange("name") {
		log.Println("The name has been changed")

		if resp, _, err := client.cloudaccountAzure.UpdateName(d.Id(), azure.CloudAccountUpdateNameRequest{
			Name: d.Get("name").(string),
		}); err != nil {
			return err
		} else {
			log.Printf("resourceCloudAccountAzureUpdate response is: %+v\n", resp)
		}
	}

	if d.HasChange("operation_mode") {
		log.Println("The operation mode has been changed")

		if resp, _, err := client.cloudaccountAzure.UpdateOperationMode(d.Id(), azure.CloudAccountUpdateOperationModeRequest{
			OperationMode: d.Get("operation_mode").(string),
		}); err != nil {
			return err
		} else {
			log.Printf("resourceCloudAccountAzureUpdate response is: %+v\n", resp)
		}
	}

	if d.HasChange("credentials") {
		log.Println("The credentials has been changed")

		if resp, _, err := client.cloudaccountAzure.UpdateCredentials(d.Id(), azure.CloudAccountUpdateCredentialsRequest{
			ApplicationID:  d.Get("credentials.client_id").(string),
			ApplicationKey: d.Get("credentials.client_password").(string),
		}); err != nil {
			return err
		} else {
			log.Printf("resourceCloudAccountAzureUpdate response is: %+v\n", resp)
		}
	}

	if d.HasChange("organizational_unit_id") {
		log.Println("The organizational unit id has been changed")

		if resp, _, err := client.cloudaccountAzure.UpdateOrganizationalID(d.Id(), azure.CloudAccountUpdateOrganizationalIDRequest{
			OrganizationalUnitID: d.Get("organizational_unit_id").(string),
		}); err != nil {
			return err
		} else {
			log.Printf("resourceCloudAccountAzureUpdate response is: %+v\n", resp)
		}
	}

	return nil
}

func expandCloudAccountAzureRequest(d *schema.ResourceData) azure.CloudAccountRequest {
	req := azure.CloudAccountRequest{
		Name:           d.Get("name").(string),
		SubscriptionID: d.Get("subscription_id").(string),
		TenantID:       d.Get("tenant_id").(string),
		OperationMode:  d.Get("operation_mode").(string),
		Credentials: azure.CloudAccountCredentials{
			ClientID:       d.Get("credentials.client_id").(string),
			ClientPassword: d.Get("credentials.client_password").(string),
		},
		OrganizationalUnitID:   d.Get("organizational_unit_id").(string),
		OrganizationalUnitPath: d.Get("organizational_unit_path").(string),
		OrganizationalUnitName: d.Get("organizational_unit_name").(string),
	}

	if r, ok := d.GetOk("creation_date"); ok {
		formatTemplate := "2006-01-02 15:04:05"
		creationDateTime, _ := time.Parse(formatTemplate, r.(string))
		req.CreationDate = creationDateTime
	}
	return req
}
