package pagerduty

import (
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyScheduleCreate,
		Read:   resourcePagerDutyScheduleRead,
		Update: resourcePagerDutyScheduleUpdate,
		Delete: resourcePagerDutyScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyScheduleImport,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"schedule_layer": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"start": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return false
								}
								return true
							},
						},
						"end": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotation_virtual_start": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"rotation_turn_length_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"users": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"restriction": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time_of_day": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"duration_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildScheduleLayers(d *schema.ResourceData, scheduleLayers *[]interface{}) *[]pagerduty.ScheduleLayer {

	pagerdutyLayers := make([]pagerduty.ScheduleLayer, len(*scheduleLayers))

	for i, l := range *scheduleLayers {
		layer := l.(map[string]interface{})

		scheduleLayer := pagerduty.ScheduleLayer{
			Name:                      layer["name"].(string),
			Start:                     layer["start"].(string),
			End:                       layer["end"].(string),
			RotationVirtualStart:      layer["rotation_virtual_start"].(string),
			RotationTurnLengthSeconds: uint(layer["rotation_turn_length_seconds"].(int)),
		}

		if layer["id"] != nil || layer["id"] != "" {
			scheduleLayer.ID = layer["id"].(string)
		}

		for _, u := range layer["users"].([]interface{}) {
			scheduleLayer.Users = append(
				scheduleLayer.Users,
				pagerduty.UserReference{
					User: pagerduty.APIObject{
						ID:   u.(string),
						Type: "user_reference"},
				},
			)
		}

		restrictions := layer["restriction"].([]interface{})

		if len(restrictions) > 0 {
			for _, r := range restrictions {
				restriction := r.(map[string]interface{})
				scheduleLayer.Restrictions = append(
					scheduleLayer.Restrictions,
					pagerduty.Restriction{
						Type:            restriction["type"].(string),
						StartTimeOfDay:  restriction["start_time_of_day"].(string),
						DurationSeconds: uint(restriction["duration_seconds"].(int)),
					},
				)
			}
		}

		pagerdutyLayers[i] = scheduleLayer

	}

	return &pagerdutyLayers
}

func buildScheduleStruct(d *schema.ResourceData) (*pagerduty.Schedule, error) {
	pagerdutyLayers := d.Get("schedule_layer").([]interface{})

	schedule := pagerduty.Schedule{
		Name:     d.Get("name").(string),
		TimeZone: d.Get("time_zone").(string),
	}

	schedule.ScheduleLayers = *buildScheduleLayers(d, &pagerdutyLayers)

	if attr, ok := d.GetOk("description"); ok {
		schedule.Description = attr.(string)
	}

	return &schedule, nil
}

func resourcePagerDutyScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	s, _ := buildScheduleStruct(d)

	log.Printf("[INFO] Creating PagerDuty schedule: %s", s.Name)

	e, err := client.CreateSchedule(*s)

	if err != nil {
		return err
	}

	d.SetId(e.ID)

	return resourcePagerDutyScheduleRead(d, meta)
}

func resourcePagerDutyScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty schedule: %s", d.Id())

	s, err := client.GetSchedule(d.Id(), pagerduty.GetScheduleOptions{})

	if err != nil {
		return err
	}

	d.Set("name", s.Name)
	d.Set("description", s.Description)

	scheduleLayers := make([]map[string]interface{}, 0, len(s.ScheduleLayers))

	for _, layer := range s.ScheduleLayers {
		restrictions := make([]map[string]interface{}, 0, len(layer.Restrictions))

		for _, r := range layer.Restrictions {
			restrictions = append(restrictions, map[string]interface{}{
				"duration_seconds":  r.DurationSeconds,
				"start_time_of_day": r.StartTimeOfDay,
				"type":              r.Type,
			})
		}

		users := make([]string, 0, len(layer.Users))

		for _, u := range layer.Users {
			users = append(users, u.User.ID)
		}

		scheduleLayers = append(scheduleLayers, map[string]interface{}{
			"id":    layer.ID,
			"name":  layer.Name,
			"start": layer.Start,
			"end":   layer.End,
			"users": users,
			"rotation_turn_length_seconds": layer.RotationTurnLengthSeconds,
			"rotation_virtual_start":       layer.RotationVirtualStart,
			"restriction":                  restrictions,
		})
	}

	d.Set("schedule_layer", scheduleLayers)

	return nil
}

func resourcePagerDutyScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	e, _ := buildScheduleStruct(d)

	d.MarkNewResource()

	log.Printf("[INFO] Updating PagerDuty schedule: %s", d.Id())

	e, err := client.UpdateSchedule(d.Id(), *e)

	if err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty schedule: %s", d.Id())

	err := client.DeleteSchedule(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourcePagerDutyScheduleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourcePagerDutyScheduleRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
