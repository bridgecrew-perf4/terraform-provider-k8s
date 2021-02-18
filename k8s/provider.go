package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/supplycom/k8s_client_go"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kubeconfig": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cluster_name": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"context_name": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"user_name": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"client_certificate_authority": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"client_certificate_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"client_key_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"token": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"k8s_ingresses": dataSourceIngresses(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	kubeconfig := d.Get("kubeconfig").(string)
	cluster_name := d.Get("cluster_name").(string)
	context_name := d.Get("context_name").(string)
	user_name := d.Get("user_name").(string)
	api_host := d.Get("api_host").(string)
	client_certificate_authority := d.Get("client_certificate_authority").(string)
	client_certificate_data := d.Get("client_certificate_data").(string)
	client_key_data := d.Get("client_key_data").(string)
	token := d.Get("token").(string)

	//errorMessage := fmt.Sprintf("kubeconfig: %s\n", kubeconfig) +
	//	fmt.Sprintf("cluster_name: %s\n", cluster_name) +
	//	fmt.Sprintf("context_name: %s\n", context_name) +
	//	fmt.Sprintf("user_name: %s\n", user_name) +
	//	fmt.Sprintf("api_host: %s\n", api_host) +
	//	fmt.Sprintf("client_certificate_authority: %s\n", client_certificate_authority) +
	//	fmt.Sprintf("client_certificate_data: %s\n", client_certificate_data) +
	//	fmt.Sprintf("client_key_data: %s\n", client_key_data) +
	//	fmt.Sprintf("token: %s\n", token)
	//
	//diags = append(diags, diag.FromErr(errors.New("DEBUG"))...)
	//diags = append(diags, diag.FromErr(errors.New(errorMessage))...)
	//
	//return nil, diags

	if kubeconfig != ""  {
		//diags = append(diags, diag.FromErr(errors.New("kubeconfig init"))...)
		//return nil, diags

		clientset, clientsetError := k8s_client_go.NewClientFromKubeFile(api_host, kubeconfig)
		if clientsetError != nil {
			diags = append(diags, diag.FromErr(clientsetError)...)
			diags = append(diags, diag.FromErr(errors.New("new error here"))...)
			return nil, diags
		}
		return clientset, diags
	} else if kubeconfig == "" &&
		cluster_name != "" &&
		context_name != "" &&
		user_name != "" &&
		api_host != "" &&
		client_certificate_authority != "" &&
		client_certificate_data != "" &&
		client_key_data != "" &&
		token != "" {

		clientset, clientsetError := k8s_client_go.NewClientFromKubeCreds(cluster_name, context_name, user_name, api_host, client_certificate_authority, client_certificate_data, client_key_data, token)

		if clientsetError != nil {
			diags = append(diags, diag.FromErr(clientsetError)...)
			diags = append(diags, diag.FromErr(errors.New("init error with cluster_name, context_name, user_name, api_host, client_certificate_authority, client_certificate_data, client_key_data, token"))...)

			return nil, diags
		}

		return clientset, diags
	} else if kubeconfig == "" &&
		api_host != "" &&
		client_certificate_authority != "" &&
		token != "" {

		clientset, clientsetError := k8s_client_go.NewClientFromToken(api_host, client_certificate_authority, token)

		if clientsetError != nil {
			diags = append(diags, diag.FromErr(clientsetError)...)
			diags = append(diags, diag.FromErr(errors.New("init with api_host, client_certificate_authority and token failed"))...)
			return nil, diags
		}

		return clientset, diags
	} else {
		errorMessage := fmt.Sprintf("kubeconfig: %s\n", kubeconfig) +
			fmt.Sprintf("cluster_name: %s\n", cluster_name) +
			fmt.Sprintf("context_name: %s\n", context_name) +
			fmt.Sprintf("user_name: %s\n", user_name) +
			fmt.Sprintf("api_host: %s\n", api_host) +
			fmt.Sprintf("client_certificate_authority: %s\n", client_certificate_authority) +
			fmt.Sprintf("client_certificate_data: %s\n", client_certificate_data) +
			fmt.Sprintf("client_key_data: %s\n", client_key_data) +
			fmt.Sprintf("token: %s\n", token)

		diags = append(diags, diag.FromErr(errors.New("no valid init parameters"))...)
		diags = append(diags, diag.FromErr(errors.New(errorMessage))...)
		return nil, diags
	}
}