package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
	"github.com/hashicorp/aws-sdk-go/gen/rds"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsDbInstanceCreate,
		Read:   resourceAwsDbInstanceRead,
		Update: resourceAwsDbInstanceUpdate,
		Delete: resourceAwsDbInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"engine": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_encrypted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"allocated_storage": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"backup_retention_period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"backup_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"iops": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"maintenance_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"multi_az": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"publicly_accessible": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"vpc_security_group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},

			"security_group_names": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},

			"final_snapshot_identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"db_subnet_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"parameter_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// apply_immediately is used to determine when the update modifications
			// take place.
			// See http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.DBInstance.Modifying.html
			"apply_immediately": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn
	tags := tagsFromMapRDS(d.Get("tags").(map[string]interface{}))
	opts := rds.CreateDBInstanceMessage{
		AllocatedStorage:     aws.Integer(d.Get("allocated_storage").(int)),
		DBInstanceClass:      aws.String(d.Get("instance_class").(string)),
		DBInstanceIdentifier: aws.String(d.Get("identifier").(string)),
		DBName:               aws.String(d.Get("name").(string)),
		MasterUsername:       aws.String(d.Get("username").(string)),
		MasterUserPassword:   aws.String(d.Get("password").(string)),
		Engine:               aws.String(d.Get("engine").(string)),
		EngineVersion:        aws.String(d.Get("engine_version").(string)),
		StorageEncrypted:     aws.Boolean(d.Get("storage_encrypted").(bool)),
		Tags:                 tags,
	}

	if attr, ok := d.GetOk("storage_type"); ok {
		opts.StorageType = aws.String(attr.(string))
	}

	attr := d.Get("backup_retention_period")
	opts.BackupRetentionPeriod = aws.Integer(attr.(int))

	if attr, ok := d.GetOk("iops"); ok {
		opts.IOPS = aws.Integer(attr.(int))
	}

	if attr, ok := d.GetOk("port"); ok {
		opts.Port = aws.Integer(attr.(int))
	}

	if attr, ok := d.GetOk("multi_az"); ok {
		opts.MultiAZ = aws.Boolean(attr.(bool))
	}

	if attr, ok := d.GetOk("availability_zone"); ok {
		opts.AvailabilityZone = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("maintenance_window"); ok {
		opts.PreferredMaintenanceWindow = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("backup_window"); ok {
		opts.PreferredBackupWindow = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("publicly_accessible"); ok {
		opts.PubliclyAccessible = aws.Boolean(attr.(bool))
	}

	if attr, ok := d.GetOk("db_subnet_group_name"); ok {
		opts.DBSubnetGroupName = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("parameter_group_name"); ok {
		opts.DBParameterGroupName = aws.String(attr.(string))
	}

	if attr := d.Get("vpc_security_group_ids").(*schema.Set); attr.Len() > 0 {
		var s []string
		for _, v := range attr.List() {
			s = append(s, v.(string))
		}
		opts.VPCSecurityGroupIDs = s
	}

	if attr := d.Get("security_group_names").(*schema.Set); attr.Len() > 0 {
		var s []string
		for _, v := range attr.List() {
			s = append(s, v.(string))
		}
		opts.DBSecurityGroups = s
	}

	log.Printf("[DEBUG] DB Instance create configuration: %#v", opts)
	_, err := conn.CreateDBInstance(&opts)
	if err != nil {
		return fmt.Errorf("Error creating DB Instance: %s", err)
	}

	d.SetId(d.Get("identifier").(string))

	log.Printf("[INFO] DB Instance ID: %s", d.Id())

	log.Println(
		"[INFO] Waiting for DB Instance to be available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "backing-up", "modifying"},
		Target:     "available",
		Refresh:    resourceAwsDbInstanceStateRefreshFunc(d, meta),
		Timeout:    40 * time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceAwsDbInstanceRead(d, meta)
}

func resourceAwsDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	v, err := resourceAwsBbInstanceRetrieve(d, meta)

	if err != nil {
		return err
	}
	if v == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", v.DBName)
	d.Set("username", v.MasterUsername)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("allocated_storage", v.AllocatedStorage)
	d.Set("storage_type", v.StorageType)
	d.Set("instance_class", v.DBInstanceClass)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("backup_retention_period", v.BackupRetentionPeriod)
	d.Set("backup_window", v.PreferredBackupWindow)
	d.Set("maintenance_window", v.PreferredMaintenanceWindow)
	d.Set("multi_az", v.MultiAZ)
	if v.DBSubnetGroup != nil {
		d.Set("db_subnet_group_name", v.DBSubnetGroup.DBSubnetGroupName)
	}

	if len(v.DBParameterGroups) > 0 {
		d.Set("parameter_group_name", v.DBParameterGroups[0].DBParameterGroupName)
	}

	if v.Endpoint != nil {
		d.Set("port", v.Endpoint.Port)
		d.Set("address", v.Endpoint.Address)

		if v.Endpoint.Address != nil && v.Endpoint.Port != nil {
			d.Set("endpoint",
				fmt.Sprintf("%s:%d", *v.Endpoint.Address, *v.Endpoint.Port))
		}
	}

	d.Set("status", v.DBInstanceStatus)
	d.Set("storage_encrypted", v.StorageEncrypted)

	// list tags for resource
	// set tags
	conn := meta.(*AWSClient).rdsconn
	arn, err := buildRDSARN(d, meta)
	if err != nil {
		name := "<empty>"
		if v.DBName != "" {
			name = *v.DBName
		}

		log.Printf("[DEBUG] Error building ARN for DB Instance, not setting Tags for DB %s", name)
	} else {
		resp, err := conn.ListTagsForResource(&rds.ListTagsForResourceMessage{
			ResourceName: aws.String(arn),
		})

		if err != nil {
			log.Printf("[DEBUG] Error retreiving tags for ARN: %s", arn)
		}

		var dt []rds.Tag
		if len(resp.TagList) > 0 {
			dt = resp.TagList
		}
		d.Set("tags", tagsToMapRDS(dt))
	}

	// Create an empty schema.Set to hold all vpc security group ids
	ids := &schema.Set{
		F: func(v interface{}) int {
			return hashcode.String(v.(string))
		},
	}
	for _, v := range v.VPCSecurityGroups {
		ids.Add(*v.VPCSecurityGroupID)
	}
	d.Set("vpc_security_group_ids", ids)

	// Create an empty schema.Set to hold all security group names
	sgn := &schema.Set{
		F: func(v interface{}) int {
			return hashcode.String(v.(string))
		},
	}
	for _, v := range v.DBSecurityGroups {
		sgn.Add(*v.DBSecurityGroupName)
	}
	d.Set("security_group_names", sgn)

	return nil
}

func resourceAwsDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn

	log.Printf("[DEBUG] DB Instance destroy: %v", d.Id())

	opts := rds.DeleteDBInstanceMessage{DBInstanceIdentifier: aws.String(d.Id())}

	finalSnapshot := d.Get("final_snapshot_identifier").(string)
	if finalSnapshot == "" {
		opts.SkipFinalSnapshot = aws.Boolean(true)
	} else {
		opts.FinalDBSnapshotIdentifier = aws.String(finalSnapshot)
	}

	log.Printf("[DEBUG] DB Instance destroy configuration: %v", opts)
	if _, err := conn.DeleteDBInstance(&opts); err != nil {
		return err
	}

	log.Println(
		"[INFO] Waiting for DB Instance to be destroyed")
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "backing-up",
			"modifying", "deleting", "available"},
		Target:     "",
		Refresh:    resourceAwsDbInstanceStateRefreshFunc(d, meta),
		Timeout:    40 * time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

func resourceAwsDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn

	d.Partial(true)

	req := &rds.ModifyDBInstanceMessage{
		ApplyImmediately:     aws.Boolean(d.Get("apply_immediately").(bool)),
		DBInstanceIdentifier: aws.String(d.Id()),
	}
	d.SetPartial("apply_immediately")

	if d.HasChange("allocated_storage") {
		d.SetPartial("allocated_storage")
		req.AllocatedStorage = aws.Integer(d.Get("allocated_storage").(int))
	}
	if d.HasChange("backup_retention_period") {
		d.SetPartial("backup_retention_period")
		req.BackupRetentionPeriod = aws.Integer(d.Get("backup_retention_period").(int))
	}
	if d.HasChange("instance_class") {
		d.SetPartial("instance_class")
		req.DBInstanceClass = aws.String(d.Get("instance_class").(string))
	}
	if d.HasChange("parameter_group_name") {
		d.SetPartial("parameter_group_name")
		req.DBParameterGroupName = aws.String(d.Get("parameter_group_name").(string))
	}
	if d.HasChange("engine_version") {
		d.SetPartial("engine_version")
		req.EngineVersion = aws.String(d.Get("engine_version").(string))
	}
	if d.HasChange("iops") {
		d.SetPartial("iops")
		req.IOPS = aws.Integer(d.Get("iops").(int))
	}
	if d.HasChange("backup_window") {
		d.SetPartial("backup_window")
		req.PreferredBackupWindow = aws.String(d.Get("backup_window").(string))
	}
	if d.HasChange("maintenance_window") {
		d.SetPartial("maintenance_window")
		req.PreferredMaintenanceWindow = aws.String(d.Get("maintenance_window").(string))
	}
	if d.HasChange("password") {
		d.SetPartial("password")
		req.MasterUserPassword = aws.String(d.Get("password").(string))
	}
	if d.HasChange("multi_az") {
		d.SetPartial("multi_az")
		req.MultiAZ = aws.Boolean(d.Get("multi_az").(bool))
	}
	if d.HasChange("storage_type") {
		d.SetPartial("storage_type")
		req.StorageType = aws.String(d.Get("storage_type").(string))
	}

	if d.HasChange("vpc_security_group_ids") {
		if attr := d.Get("vpc_security_group_ids").(*schema.Set); attr.Len() > 0 {
			var s []string
			for _, v := range attr.List() {
				s = append(s, v.(string))
			}
			req.VPCSecurityGroupIDs = s
		}
	}

	if d.HasChange("vpc_security_group_ids") {
		if attr := d.Get("security_group_names").(*schema.Set); attr.Len() > 0 {
			var s []string
			for _, v := range attr.List() {
				s = append(s, v.(string))
			}
			req.DBSecurityGroups = s
		}
	}

	log.Printf("[DEBUG] DB Instance Modification request: %#v", req)
	_, err := conn.ModifyDBInstance(req)
	if err != nil {
		return fmt.Errorf("Error modifying DB Instance %s: %s", d.Id(), err)
	}

	if arn, err := buildRDSARN(d, meta); err == nil {
		if err := setTagsRDS(conn, d, arn); err != nil {
			return err
		} else {
			d.SetPartial("tags")
		}
	}
	d.Partial(false)
	return resourceAwsDbInstanceRead(d, meta)
}

func resourceAwsBbInstanceRetrieve(
	d *schema.ResourceData, meta interface{}) (*rds.DBInstance, error) {
	conn := meta.(*AWSClient).rdsconn

	opts := rds.DescribeDBInstancesMessage{
		DBInstanceIdentifier: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] DB Instance describe configuration: %#v", opts)

	resp, err := conn.DescribeDBInstances(&opts)

	if err != nil {
		dbinstanceerr, ok := err.(aws.APIError)
		if ok && dbinstanceerr.Code == "DBInstanceNotFound" {
			return nil, nil
		}
		return nil, fmt.Errorf("Error retrieving DB Instances: %s", err)
	}

	if len(resp.DBInstances) != 1 ||
		*resp.DBInstances[0].DBInstanceIdentifier != d.Id() {
		if err != nil {
			return nil, nil
		}
	}

	v := resp.DBInstances[0]

	return &v, nil
}

func resourceAwsDbInstanceStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceAwsBbInstanceRetrieve(d, meta)

		if err != nil {
			log.Printf("Error on retrieving DB Instance when waiting: %s", err)
			return nil, "", err
		}

		if v == nil {
			return nil, "", nil
		}

		return v, *v.DBInstanceStatus, nil
	}
}

func buildRDSARN(d *schema.ResourceData, meta interface{}) (string, error) {
	iamconn := meta.(*AWSClient).iamconn
	region := meta.(*AWSClient).region
	// An zero value GetUserRequest{} defers to the currently logged in user
	resp, err := iamconn.GetUser(&iam.GetUserRequest{})
	if err != nil {
		return "", err
	}
	user := resp.User
	arn := fmt.Sprintf("arn:aws:rds:%s:%s:db:%s", region, *user.UserID, d.Id())
	return arn, nil
}
