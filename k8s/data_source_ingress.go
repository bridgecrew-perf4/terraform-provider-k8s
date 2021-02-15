package k8s

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIngresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIngressesRead,
		Schema: map[string]*schema.Schema{
			"namespace": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"servicename": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"appname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ingresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servicename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"appname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branchname": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
