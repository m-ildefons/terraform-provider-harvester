package loadbalancer

import (
	"context"
	"time"

	loadbalancerv1 "github.com/harvester/harvester-load-balancer/pkg/apis/loadbalancer.harvesterhci.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/harvester/terraform-provider-harvester/internal/util"
	"github.com/harvester/terraform-provider-harvester/pkg/client"
	"github.com/harvester/terraform-provider-harvester/pkg/constants"
	"github.com/harvester/terraform-provider-harvester/pkg/helper"
	"github.com/harvester/terraform-provider-harvester/pkg/importer"
)

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: Schema(),
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(2 * time.Minute),
			Read:    schema.DefaultTimeout(2 * time.Minute),
			Update:  schema.DefaultTimeout(2 * time.Minute),
			Delete:  schema.DefaultTimeout(2 * time.Minute),
			Default: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceLoadBalancerCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace := data.Get(constants.FieldCommonNamespace).(string)
	name := data.Get(constants.FieldCommonName).(string)
	toCreate, err := util.ResourceConstruct(data, Creator(namespace, name))
	if err != nil {
		return diag.FromErr(err)
	}
	lb, err := c.HarvesterLoadbalancerClient.
		LoadbalancerV1beta1().
		LoadBalancers(namespace).
		Create(ctx, toCreate.(*loadbalancerv1.LoadBalancer), metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(helper.BuildID(namespace, name))
	return diag.FromErr(resourceLoadBalancerImport(data, lb))
}

func resourceLoadBalancerRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.FromErr(nil)
}

func resourceLoadBalancerUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.FromErr(nil)
}

func resourceLoadBalancerDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.FromErr(nil)
}

func resourceLoadBalancerImport(data *schema.ResourceData, obj *loadbalancerv1.LoadBalancer) error {
	stateGetter, err := importer.ResourceLoadBalancerStateGetter(obj)
	if err != nil {
		return err
	}
	return util.ResourceStatesSet(data, stateGetter)
}
