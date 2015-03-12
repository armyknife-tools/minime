package aws

import (
	"strings"

	"github.com/hashicorp/aws-sdk-go/aws"
	awsEC2 "github.com/hashicorp/aws-sdk-go/gen/ec2"
	"github.com/hashicorp/aws-sdk-go/gen/elb"
	"github.com/hashicorp/aws-sdk-go/gen/rds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/goamz/ec2"
)

// Takes the result of flatmap.Expand for an array of listeners and
// returns ELB API compatible objects
func expandListeners(configured []interface{}) ([]elb.Listener, error) {
	listeners := make([]elb.Listener, 0, len(configured))

	// Loop over our configured listeners and create
	// an array of goamz compatabile objects
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := elb.Listener{
			InstancePort:     aws.Integer(data["instance_port"].(int)),
			InstanceProtocol: aws.String(data["instance_protocol"].(string)),
			LoadBalancerPort: aws.Integer(data["lb_port"].(int)),
			Protocol:         aws.String(data["lb_protocol"].(string)),
		}

		if v, ok := data["ssl_certificate_id"]; ok {
			l.SSLCertificateID = aws.String(v.(string))
		}

		listeners = append(listeners, l)
	}

	return listeners, nil
}

// Takes the result of flatmap.Expand for an array of ingress/egress
// security group rules and returns EC2 API compatible objects
func expandIPPerms(id string, configured []interface{}) []awsEC2.IPPermission {
	perms := make([]awsEC2.IPPermission, len(configured))
	for i, mRaw := range configured {
		var perm awsEC2.IPPermission
		m := mRaw.(map[string]interface{})

		perm.FromPort = aws.Integer(m["from_port"].(int))
		perm.ToPort = aws.Integer(m["to_port"].(int))
		perm.IPProtocol = aws.String(m["protocol"].(string))

		var groups []string
		if raw, ok := m["security_groups"]; ok {
			list := raw.(*schema.Set).List()
			for _, v := range list {
				groups = append(groups, v.(string))
			}
		}
		if v, ok := m["self"]; ok && v.(bool) {
			groups = append(groups, id)
		}

		if len(groups) > 0 {
			perm.UserIDGroupPairs = make([]awsEC2.UserIDGroupPair, len(groups))
			for i, name := range groups {
				ownerId, id := "", name
				if items := strings.Split(id, "/"); len(items) > 1 {
					ownerId, id = items[0], items[1]
				}

				perm.UserIDGroupPairs[i] = awsEC2.UserIDGroupPair{
					GroupID: aws.String(id),
					UserID:  aws.String(ownerId),
				}
			}
		}

		if raw, ok := m["cidr_blocks"]; ok {
			list := raw.([]interface{})
			perm.IPRanges = make([]awsEC2.IPRange, len(list))
			for i, v := range list {
				perm.IPRanges[i] = awsEC2.IPRange{aws.String(v.(string))}
			}
		}

		perms[i] = perm
	}

	return perms
}

// Takes the result of flatmap.Expand for an array of parameters and
// returns Parameter API compatible objects
func expandParameters(configured []interface{}) ([]rds.Parameter, error) {
	parameters := make([]rds.Parameter, 0, len(configured))

	// Loop over our configured parameters and create
	// an array of goamz compatabile objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		p := rds.Parameter{
			ApplyMethod:    aws.String(data["apply_method"].(string)),
			ParameterName:  aws.String(data["name"].(string)),
			ParameterValue: aws.String(data["value"].(string)),
		}

		parameters = append(parameters, p)
	}

	return parameters, nil
}

// Flattens a health check into something that flatmap.Flatten()
// can handle
func flattenHealthCheck(check *elb.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	chk := make(map[string]interface{})
	chk["unhealthy_threshold"] = *check.UnhealthyThreshold
	chk["healthy_threshold"] = *check.HealthyThreshold
	chk["target"] = *check.Target
	chk["timeout"] = *check.Timeout
	chk["interval"] = *check.Interval

	result = append(result, chk)

	return result
}

// Flattens an array of UserSecurityGroups into a []string
func flattenSecurityGroups(list []ec2.UserSecurityGroup) []string {
	result := make([]string, 0, len(list))
	for _, g := range list {
		result = append(result, g.Id)
	}
	return result
}

// Flattens an array of UserSecurityGroups into a []string
func flattenSecurityGroupsSDK(list []awsEC2.UserIDGroupPair) []string {
	result := make([]string, 0, len(list))
	for _, g := range list {
		result = append(result, *g.GroupID)
	}
	return result
}

// Flattens an array of Instances into a []string
func flattenInstances(list []elb.Instance) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		result = append(result, *i.InstanceID)
	}
	return result
}

// Expands an array of String Instance IDs into a []Instances
func expandInstanceString(list []interface{}) []elb.Instance {
	result := make([]elb.Instance, 0, len(list))
	for _, i := range list {
		result = append(result, elb.Instance{aws.String(i.(string))})
	}
	return result
}

// Flattens an array of Listeners into a []map[string]interface{}
func flattenListeners(list []elb.ListenerDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"instance_port":     *i.Listener.InstancePort,
			"instance_protocol": strings.ToLower(*i.Listener.InstanceProtocol),
			"lb_port":           *i.Listener.LoadBalancerPort,
			"lb_protocol":       strings.ToLower(*i.Listener.Protocol),
		}
		// SSLCertificateID is optional, and may be nil
		if i.Listener.SSLCertificateID != nil {
			l["ssl_certificate_id"] = *i.Listener.SSLCertificateID
		}
		result = append(result, l)
	}
	return result
}

// Flattens an array of Parameters into a []map[string]interface{}
func flattenParameters(list []rds.Parameter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		result = append(result, map[string]interface{}{
			"name":  strings.ToLower(*i.ParameterName),
			"value": strings.ToLower(*i.ParameterValue),
		})
	}
	return result
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}
