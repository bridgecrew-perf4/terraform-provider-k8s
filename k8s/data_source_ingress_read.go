package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/supplycom/k8s_client_go"
	"k8s.io/api/extensions/v1beta1"
	"strconv"
	"time"
)

func dataSourceIngressesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	namespace, namespaceSet := d.GetOk("namespace")
	if namespaceSet == false {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "debug",
			Detail:   fmt.Sprintf("namespaceSet: %t, namespace: %s", namespaceSet, namespace),
		})
	}

	k8sClient := m.(*k8s_client_go.Client)

	ingresses, ingressesError := k8sClient.ListIngresses(namespace.(string))
	if ingressesError != nil {
		diags = append(diags, diag.FromErr(ingressesError)...)
		return diags
	}

	serviceName, serviceNameSet := d.GetOk("servicename")
	serviceNameFilteredIngresses := make([]v1beta1.Ingress, 0)
	if serviceNameSet == true {
		for _, ingress := range ingresses {
			ingressServiceName, ingressServiceNameError := GetServiceName(ingress)
			if ingressServiceNameError != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "can't read service name",
					Detail:   fmt.Sprintf("error reading service name from ingress: %s",ingress.Name),
				})
			}
			if ingressServiceName == serviceName {
				serviceNameFilteredIngresses = append(serviceNameFilteredIngresses, ingress)
			}
		}
	} else {
		serviceNameFilteredIngresses = ingresses
	}

	appName, appNameSet := d.GetOk("appname")
	appNameFilteredIngresses := make([]v1beta1.Ingress, 0)
	if appNameSet == true {
		for _, ingress := range serviceNameFilteredIngresses {
			if ingress.ObjectMeta.Labels["appname"] == appName {

				appNameFilteredIngresses = append(appNameFilteredIngresses, ingress)
			}
		}
	} else {
		appNameFilteredIngresses = serviceNameFilteredIngresses
	}

	resultIngresses := make([]interface{}, len(appNameFilteredIngresses))
	for i, ingress := range appNameFilteredIngresses {
		it := make(map[string]interface{})
		it["id"] = ingress.GetUID()
		it["name"] = ingress.GetName()
		it["namespace"] = ingress.GetNamespace()
		it["servicename"], _ = GetServiceName(ingress)
		it["appname"] = ingress.Labels["appname"]
		it["branchname"] = ingress.Labels["branchname"]

		resultIngresses[i] = it
	}

	if err := d.Set("ingresses", resultIngresses); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	//diags = append(diags, diag.Diagnostic{
	//	Severity: diag.Error,
	//	Summary:  "nil check",
	//	Detail: fmt.Sprintln("resultIngresses count: %d", len(resultIngresses)) +
	//		fmt.Sprintln("serviceNameFilteredIngresses count: %d", len(serviceNameFilteredIngresses)) +
	//		fmt.Sprintln("ingresses count: %d", len(ingresses)) +
	//		fmt.Sprintln("k8sClient hostUrl: %s", k8sClient.HostURL) +
	//		fmt.Sprintln("k8sClientSet: %s", k8sClient.K8sClientSet),
	//})
	return diags
}

func GetServiceName(ingress v1beta1.Ingress) (string, error) {
	serviceName := ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.ServiceName
	if recover() != nil {
		err := errors.New("problem referencing service name on ingress")
		return "", err
	}
	return serviceName, nil
}
