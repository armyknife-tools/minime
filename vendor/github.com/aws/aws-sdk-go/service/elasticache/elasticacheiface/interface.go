// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package elasticacheiface provides an interface for the Amazon ElastiCache.
package elasticacheiface

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

// ElastiCacheAPI is the interface type for elasticache.ElastiCache.
type ElastiCacheAPI interface {
	AddTagsToResourceRequest(*elasticache.AddTagsToResourceInput) (*request.Request, *elasticache.TagListMessage)

	AddTagsToResource(*elasticache.AddTagsToResourceInput) (*elasticache.TagListMessage, error)

	AuthorizeCacheSecurityGroupIngressRequest(*elasticache.AuthorizeCacheSecurityGroupIngressInput) (*request.Request, *elasticache.AuthorizeCacheSecurityGroupIngressOutput)

	AuthorizeCacheSecurityGroupIngress(*elasticache.AuthorizeCacheSecurityGroupIngressInput) (*elasticache.AuthorizeCacheSecurityGroupIngressOutput, error)

	CopySnapshotRequest(*elasticache.CopySnapshotInput) (*request.Request, *elasticache.CopySnapshotOutput)

	CopySnapshot(*elasticache.CopySnapshotInput) (*elasticache.CopySnapshotOutput, error)

	CreateCacheClusterRequest(*elasticache.CreateCacheClusterInput) (*request.Request, *elasticache.CreateCacheClusterOutput)

	CreateCacheCluster(*elasticache.CreateCacheClusterInput) (*elasticache.CreateCacheClusterOutput, error)

	CreateCacheParameterGroupRequest(*elasticache.CreateCacheParameterGroupInput) (*request.Request, *elasticache.CreateCacheParameterGroupOutput)

	CreateCacheParameterGroup(*elasticache.CreateCacheParameterGroupInput) (*elasticache.CreateCacheParameterGroupOutput, error)

	CreateCacheSecurityGroupRequest(*elasticache.CreateCacheSecurityGroupInput) (*request.Request, *elasticache.CreateCacheSecurityGroupOutput)

	CreateCacheSecurityGroup(*elasticache.CreateCacheSecurityGroupInput) (*elasticache.CreateCacheSecurityGroupOutput, error)

	CreateCacheSubnetGroupRequest(*elasticache.CreateCacheSubnetGroupInput) (*request.Request, *elasticache.CreateCacheSubnetGroupOutput)

	CreateCacheSubnetGroup(*elasticache.CreateCacheSubnetGroupInput) (*elasticache.CreateCacheSubnetGroupOutput, error)

	CreateReplicationGroupRequest(*elasticache.CreateReplicationGroupInput) (*request.Request, *elasticache.CreateReplicationGroupOutput)

	CreateReplicationGroup(*elasticache.CreateReplicationGroupInput) (*elasticache.CreateReplicationGroupOutput, error)

	CreateSnapshotRequest(*elasticache.CreateSnapshotInput) (*request.Request, *elasticache.CreateSnapshotOutput)

	CreateSnapshot(*elasticache.CreateSnapshotInput) (*elasticache.CreateSnapshotOutput, error)

	DeleteCacheClusterRequest(*elasticache.DeleteCacheClusterInput) (*request.Request, *elasticache.DeleteCacheClusterOutput)

	DeleteCacheCluster(*elasticache.DeleteCacheClusterInput) (*elasticache.DeleteCacheClusterOutput, error)

	DeleteCacheParameterGroupRequest(*elasticache.DeleteCacheParameterGroupInput) (*request.Request, *elasticache.DeleteCacheParameterGroupOutput)

	DeleteCacheParameterGroup(*elasticache.DeleteCacheParameterGroupInput) (*elasticache.DeleteCacheParameterGroupOutput, error)

	DeleteCacheSecurityGroupRequest(*elasticache.DeleteCacheSecurityGroupInput) (*request.Request, *elasticache.DeleteCacheSecurityGroupOutput)

	DeleteCacheSecurityGroup(*elasticache.DeleteCacheSecurityGroupInput) (*elasticache.DeleteCacheSecurityGroupOutput, error)

	DeleteCacheSubnetGroupRequest(*elasticache.DeleteCacheSubnetGroupInput) (*request.Request, *elasticache.DeleteCacheSubnetGroupOutput)

	DeleteCacheSubnetGroup(*elasticache.DeleteCacheSubnetGroupInput) (*elasticache.DeleteCacheSubnetGroupOutput, error)

	DeleteReplicationGroupRequest(*elasticache.DeleteReplicationGroupInput) (*request.Request, *elasticache.DeleteReplicationGroupOutput)

	DeleteReplicationGroup(*elasticache.DeleteReplicationGroupInput) (*elasticache.DeleteReplicationGroupOutput, error)

	DeleteSnapshotRequest(*elasticache.DeleteSnapshotInput) (*request.Request, *elasticache.DeleteSnapshotOutput)

	DeleteSnapshot(*elasticache.DeleteSnapshotInput) (*elasticache.DeleteSnapshotOutput, error)

	DescribeCacheClustersRequest(*elasticache.DescribeCacheClustersInput) (*request.Request, *elasticache.DescribeCacheClustersOutput)

	DescribeCacheClusters(*elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error)

	DescribeCacheClustersPages(*elasticache.DescribeCacheClustersInput, func(*elasticache.DescribeCacheClustersOutput, bool) bool) error

	DescribeCacheEngineVersionsRequest(*elasticache.DescribeCacheEngineVersionsInput) (*request.Request, *elasticache.DescribeCacheEngineVersionsOutput)

	DescribeCacheEngineVersions(*elasticache.DescribeCacheEngineVersionsInput) (*elasticache.DescribeCacheEngineVersionsOutput, error)

	DescribeCacheEngineVersionsPages(*elasticache.DescribeCacheEngineVersionsInput, func(*elasticache.DescribeCacheEngineVersionsOutput, bool) bool) error

	DescribeCacheParameterGroupsRequest(*elasticache.DescribeCacheParameterGroupsInput) (*request.Request, *elasticache.DescribeCacheParameterGroupsOutput)

	DescribeCacheParameterGroups(*elasticache.DescribeCacheParameterGroupsInput) (*elasticache.DescribeCacheParameterGroupsOutput, error)

	DescribeCacheParameterGroupsPages(*elasticache.DescribeCacheParameterGroupsInput, func(*elasticache.DescribeCacheParameterGroupsOutput, bool) bool) error

	DescribeCacheParametersRequest(*elasticache.DescribeCacheParametersInput) (*request.Request, *elasticache.DescribeCacheParametersOutput)

	DescribeCacheParameters(*elasticache.DescribeCacheParametersInput) (*elasticache.DescribeCacheParametersOutput, error)

	DescribeCacheParametersPages(*elasticache.DescribeCacheParametersInput, func(*elasticache.DescribeCacheParametersOutput, bool) bool) error

	DescribeCacheSecurityGroupsRequest(*elasticache.DescribeCacheSecurityGroupsInput) (*request.Request, *elasticache.DescribeCacheSecurityGroupsOutput)

	DescribeCacheSecurityGroups(*elasticache.DescribeCacheSecurityGroupsInput) (*elasticache.DescribeCacheSecurityGroupsOutput, error)

	DescribeCacheSecurityGroupsPages(*elasticache.DescribeCacheSecurityGroupsInput, func(*elasticache.DescribeCacheSecurityGroupsOutput, bool) bool) error

	DescribeCacheSubnetGroupsRequest(*elasticache.DescribeCacheSubnetGroupsInput) (*request.Request, *elasticache.DescribeCacheSubnetGroupsOutput)

	DescribeCacheSubnetGroups(*elasticache.DescribeCacheSubnetGroupsInput) (*elasticache.DescribeCacheSubnetGroupsOutput, error)

	DescribeCacheSubnetGroupsPages(*elasticache.DescribeCacheSubnetGroupsInput, func(*elasticache.DescribeCacheSubnetGroupsOutput, bool) bool) error

	DescribeEngineDefaultParametersRequest(*elasticache.DescribeEngineDefaultParametersInput) (*request.Request, *elasticache.DescribeEngineDefaultParametersOutput)

	DescribeEngineDefaultParameters(*elasticache.DescribeEngineDefaultParametersInput) (*elasticache.DescribeEngineDefaultParametersOutput, error)

	DescribeEngineDefaultParametersPages(*elasticache.DescribeEngineDefaultParametersInput, func(*elasticache.DescribeEngineDefaultParametersOutput, bool) bool) error

	DescribeEventsRequest(*elasticache.DescribeEventsInput) (*request.Request, *elasticache.DescribeEventsOutput)

	DescribeEvents(*elasticache.DescribeEventsInput) (*elasticache.DescribeEventsOutput, error)

	DescribeEventsPages(*elasticache.DescribeEventsInput, func(*elasticache.DescribeEventsOutput, bool) bool) error

	DescribeReplicationGroupsRequest(*elasticache.DescribeReplicationGroupsInput) (*request.Request, *elasticache.DescribeReplicationGroupsOutput)

	DescribeReplicationGroups(*elasticache.DescribeReplicationGroupsInput) (*elasticache.DescribeReplicationGroupsOutput, error)

	DescribeReplicationGroupsPages(*elasticache.DescribeReplicationGroupsInput, func(*elasticache.DescribeReplicationGroupsOutput, bool) bool) error

	DescribeReservedCacheNodesRequest(*elasticache.DescribeReservedCacheNodesInput) (*request.Request, *elasticache.DescribeReservedCacheNodesOutput)

	DescribeReservedCacheNodes(*elasticache.DescribeReservedCacheNodesInput) (*elasticache.DescribeReservedCacheNodesOutput, error)

	DescribeReservedCacheNodesPages(*elasticache.DescribeReservedCacheNodesInput, func(*elasticache.DescribeReservedCacheNodesOutput, bool) bool) error

	DescribeReservedCacheNodesOfferingsRequest(*elasticache.DescribeReservedCacheNodesOfferingsInput) (*request.Request, *elasticache.DescribeReservedCacheNodesOfferingsOutput)

	DescribeReservedCacheNodesOfferings(*elasticache.DescribeReservedCacheNodesOfferingsInput) (*elasticache.DescribeReservedCacheNodesOfferingsOutput, error)

	DescribeReservedCacheNodesOfferingsPages(*elasticache.DescribeReservedCacheNodesOfferingsInput, func(*elasticache.DescribeReservedCacheNodesOfferingsOutput, bool) bool) error

	DescribeSnapshotsRequest(*elasticache.DescribeSnapshotsInput) (*request.Request, *elasticache.DescribeSnapshotsOutput)

	DescribeSnapshots(*elasticache.DescribeSnapshotsInput) (*elasticache.DescribeSnapshotsOutput, error)

	DescribeSnapshotsPages(*elasticache.DescribeSnapshotsInput, func(*elasticache.DescribeSnapshotsOutput, bool) bool) error

	ListTagsForResourceRequest(*elasticache.ListTagsForResourceInput) (*request.Request, *elasticache.TagListMessage)

	ListTagsForResource(*elasticache.ListTagsForResourceInput) (*elasticache.TagListMessage, error)

	ModifyCacheClusterRequest(*elasticache.ModifyCacheClusterInput) (*request.Request, *elasticache.ModifyCacheClusterOutput)

	ModifyCacheCluster(*elasticache.ModifyCacheClusterInput) (*elasticache.ModifyCacheClusterOutput, error)

	ModifyCacheParameterGroupRequest(*elasticache.ModifyCacheParameterGroupInput) (*request.Request, *elasticache.CacheParameterGroupNameMessage)

	ModifyCacheParameterGroup(*elasticache.ModifyCacheParameterGroupInput) (*elasticache.CacheParameterGroupNameMessage, error)

	ModifyCacheSubnetGroupRequest(*elasticache.ModifyCacheSubnetGroupInput) (*request.Request, *elasticache.ModifyCacheSubnetGroupOutput)

	ModifyCacheSubnetGroup(*elasticache.ModifyCacheSubnetGroupInput) (*elasticache.ModifyCacheSubnetGroupOutput, error)

	ModifyReplicationGroupRequest(*elasticache.ModifyReplicationGroupInput) (*request.Request, *elasticache.ModifyReplicationGroupOutput)

	ModifyReplicationGroup(*elasticache.ModifyReplicationGroupInput) (*elasticache.ModifyReplicationGroupOutput, error)

	PurchaseReservedCacheNodesOfferingRequest(*elasticache.PurchaseReservedCacheNodesOfferingInput) (*request.Request, *elasticache.PurchaseReservedCacheNodesOfferingOutput)

	PurchaseReservedCacheNodesOffering(*elasticache.PurchaseReservedCacheNodesOfferingInput) (*elasticache.PurchaseReservedCacheNodesOfferingOutput, error)

	RebootCacheClusterRequest(*elasticache.RebootCacheClusterInput) (*request.Request, *elasticache.RebootCacheClusterOutput)

	RebootCacheCluster(*elasticache.RebootCacheClusterInput) (*elasticache.RebootCacheClusterOutput, error)

	RemoveTagsFromResourceRequest(*elasticache.RemoveTagsFromResourceInput) (*request.Request, *elasticache.TagListMessage)

	RemoveTagsFromResource(*elasticache.RemoveTagsFromResourceInput) (*elasticache.TagListMessage, error)

	ResetCacheParameterGroupRequest(*elasticache.ResetCacheParameterGroupInput) (*request.Request, *elasticache.CacheParameterGroupNameMessage)

	ResetCacheParameterGroup(*elasticache.ResetCacheParameterGroupInput) (*elasticache.CacheParameterGroupNameMessage, error)

	RevokeCacheSecurityGroupIngressRequest(*elasticache.RevokeCacheSecurityGroupIngressInput) (*request.Request, *elasticache.RevokeCacheSecurityGroupIngressOutput)

	RevokeCacheSecurityGroupIngress(*elasticache.RevokeCacheSecurityGroupIngressInput) (*elasticache.RevokeCacheSecurityGroupIngressOutput, error)
}

var _ ElastiCacheAPI = (*elasticache.ElastiCache)(nil)
