package alicloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/validation"
)

// common
func validateInstancePort(v interface{}, k string) (ws []string, errors []error) {
	return validation.IntBetween(1, 65535)(v, k)
}

func validateInstanceProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocal := v.(string)
	if !isProtocalValid(protocal) {
		errors = append(errors, fmt.Errorf(
			"%q is an invalid value. Valid values are either http, https, tcp or udp",
			k))
		return
	}
	return
}

// ecs
func validateDiskCategory(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(ecs.DiskCategoryCloud),
		string(ecs.DiskCategoryCloudEfficiency),
		string(ecs.DiskCategoryCloudSSD),
	}, false)(v, k)
}

func validateInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateInstanceDescription(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringLenBetween(2, 256)(v, k)
}

func validateDiskName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateDiskDescription(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringLenBetween(2, 128)(v, k)
}

//security group
func validateSecurityGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateSecurityGroupDescription(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringLenBetween(2, 256)(v, k)
}

func validateSecurityRuleType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(GroupRuleIngress),
		string(GroupRuleEgress),
	}, false)(v, k)
}

func validateSecurityRuleIpProtocol(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(GroupRuleTcp),
		string(GroupRuleUdp),
		string(GroupRuleIcmp),
		string(GroupRuleGre),
		string(GroupRuleAll),
	}, false)(v, k)
}

func validateSecurityRuleNicType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(GroupRuleInternet),
		string(GroupRuleIntranet),
	}, false)(v, k)
}

func validateSecurityRulePolicy(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(GroupRulePolicyAccept),
		string(GroupRulePolicyDrop),
	}, false)(v, k)
}

func validateSecurityPriority(v interface{}, k string) (ws []string, errors []error) {
	return validation.IntBetween(1, 100)(v, k)
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	return validation.CIDRNetwork(0, 32)(v, k)
}

func validateRouteEntryNextHopType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		string(ecs.NextHopIntance),
		string(ecs.NextHopTunnel),
	}, false)(v, k)
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	return validation.CIDRNetwork(16, 29)(v, k)
}

// validateIoOptimized ensures that the string value is a valid IoOptimized that
// represents a IoOptimized - it adds an error otherwise
func validateIoOptimized(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"",
		string(ecs.IoOptimizedNone),
		string(ecs.IoOptimizedOptimized),
	}, false)(v, k)
}

// validateInstanceNetworkType ensures that the string value is a classic or vpc
func validateInstanceNetworkType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"",
		string(ClassicNet),
		string(VpcNet),
	}, false)(v, k)
}

func validateInstanceChargeType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"",
		string(common.PrePaid),
		string(common.PostPaid),
	}, false)(v, k)
}

func validateInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"",
		string(common.PayByBandwidth),
		string(common.PayByTraffic),
	}, false)(v, k)
}

func validateInternetMaxBandWidthOut(v interface{}, k string) (ws []string, errors []error) {
	return validation.IntBetween(1, 100)(v, k)
}

// SLB
func validateSlbName(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringLenBetween(0, 80)(v, k)
}

func validateSlbInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"paybybandwidth",
		"paybytraffic",
	}, false)(v, k)
}

func validateSlbBandwidth(v interface{}, k string) (ws []string, errors []error) {
	return validation.IntBetween(1, 1000)(v, k)
}

func validateSlbListenerBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value < 1 || value > 1000) && value != -1 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000 or -1",
			k))
		return
	}
	return
}

func validateSlbListenerScheduler(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{"wrr", "wlc"}, false)(v, k)
}

func validateSlbListenerStickySession(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{"", "on", "off"}, false)(v, k)
}

func validateSlbListenerStickySessionType(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{"", "insert", "server"}, false)(v, k)
}

func validateSlbListenerCookie(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{"", "insert", "server"}, false)(v, k)
}

func validateSlbListenerPersistenceTimeout(v interface{}, k string) (ws []string, errors []error) {
	return validation.IntBetween(0, 86400)(v, k)
}

//data source validate func
//data_source_alicloud_image
func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}

func validateImageOwners(v interface{}, k string) (ws []string, errors []error) {
	return validation.StringInSlice([]string{
		"",
		string(ecs.ImageOwnerSystem),
		string(ecs.ImageOwnerSelf),
		string(ecs.ImageOwnerOthers),
		string(ecs.ImageOwnerMarketplace),
		string(ecs.ImageOwnerDefault),
	}, false)(v, k)
}

func validateRegion(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		region := common.Region(value)
		var valid string
		for _, re := range common.ValidRegions {
			if region == re {
				return
			}
			valid = valid + ", " + string(re)
		}
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Region ID , expected %#v, got %q",
			k, valid, value))

	}
	return
}
