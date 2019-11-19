package models

import (
	"time"

	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
)

type AlertDetail struct {
	AlertId             string    `gorm:"column:alert_id" json:"alert_id"`
	AlertName           string    `gorm:"column:alert_name" json:"alert_name"`
	Disabled            bool      `gorm:"column:disabled" json:"disabled"`
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`
	RunningStatus       string    `gorm:"column:running_status" json:"running_status"`
	AlertStatus         string    `gorm:"column:alert_status" json:"alert_status"`
	PolicyId            string    `gorm:"column:policy_id" json:"policy_id"`
	RsFilterName        string    `gorm:"column:rs_filter_name" json:"rs_filter_name"`
	RsFilterParam       string    `gorm:"column:rs_filter_param" json:"rs_filter_param"`
	RsTypeName          string    `gorm:"column:rs_type_name" json:"rs_type_name"`
	ExecutorId          string    `gorm:"column:executor_id" json:"executor_id"`
	PolicyName          string    `gorm:"column:policy_name" json:"policy_name"`
	PolicyDescription   string    `gorm:"column:policy_description" json:"policy_description"`
	PolicyConfig        string    `gorm:"column:policy_config" json:"policy_config"`
	Creator             string    `gorm:"column:creator" json:"creator"`
	AvailableStartTime  string    `gorm:"column:available_start_time" json:"available_start_time"`
	AvailableEndTime    string    `gorm:"column:available_end_time" json:"available_end_time"`
	Language            string    `gorm:"column:language" json:"language"`
	Metrics             []string  `gorm:"column:metrics" json:"metrics"`
	RulesCount          uint32    `json:"rules_count"`
	PositivesCount      uint32    `json:"positives_count"`
	MostRecentAlertTime string    `json:"most_recent_alert_time"`
	NfAddressListId     string    `gorm:"column:nf_address_list_id" json:"nf_address_list_id"`
}

func AlertDetailToPb(alertDetail *AlertDetail) *pb.AlertDetail {
	pbAlertDetail := pb.AlertDetail{}
	pbAlertDetail.AlertId = alertDetail.AlertId
	pbAlertDetail.AlertName = alertDetail.AlertName
	pbAlertDetail.Disabled = alertDetail.Disabled
	pbAlertDetail.CreateTime = pbutil.ToProtoTimestamp(alertDetail.CreateTime)
	pbAlertDetail.RunningStatus = alertDetail.RunningStatus
	pbAlertDetail.AlertStatus = alertDetail.AlertStatus
	pbAlertDetail.PolicyId = alertDetail.PolicyId
	pbAlertDetail.RsFilterName = alertDetail.RsFilterName
	pbAlertDetail.RsFilterParam = alertDetail.RsFilterParam
	pbAlertDetail.RsTypeName = alertDetail.RsTypeName
	pbAlertDetail.ExecutorId = alertDetail.ExecutorId
	pbAlertDetail.PolicyName = alertDetail.PolicyName
	pbAlertDetail.PolicyDescription = alertDetail.PolicyDescription
	pbAlertDetail.PolicyConfig = alertDetail.PolicyConfig
	pbAlertDetail.Creator = alertDetail.Creator
	pbAlertDetail.AvailableStartTime = alertDetail.AvailableStartTime
	pbAlertDetail.AvailableEndTime = alertDetail.AvailableEndTime
	pbAlertDetail.Language = alertDetail.Language
	pbAlertDetail.Metrics = alertDetail.Metrics
	pbAlertDetail.RulesCount = alertDetail.RulesCount
	pbAlertDetail.PositivesCount = alertDetail.PositivesCount
	pbAlertDetail.MostRecentAlertTime = alertDetail.MostRecentAlertTime
	pbAlertDetail.NfAddressListId = alertDetail.NfAddressListId
	return &pbAlertDetail
}

func ParseAldSet2PbSet(inAlds []*AlertDetail) []*pb.AlertDetail {
	var pbAlds []*pb.AlertDetail
	for _, inAld := range inAlds {
		pbAld := AlertDetailToPb(inAld)
		pbAlds = append(pbAlds, pbAld)
	}
	return pbAlds
}

type ResourceStatus struct {
	ResourceName       string `gorm:"column:resource_name" json:"resource_name"`
	CurrentLevel       string `json:"current_level"`
	PositiveCount      uint32 `json:"positive_count"`
	CumulatedSendCount uint32 `json:"cumulated_send_count"`
	NextResendInterval uint32 `json:"next_resend_interval"`
	NextSendableTime   string `json:"next_sendable_time"`
	AggregatedAlerts   string `json:"aggregated_alerts"`
}

type AlertStatus struct {
	RuleId           string           `gorm:"column:rule_id" json:"rule_id"`
	RuleName         string           `gorm:"column:rule_name" json:"rule_name"`
	Disabled         bool             `gorm:"column:disabled" json:"disabled"`
	MonitorPeriods   uint32           `gorm:"column:monitor_periods" json:"monitor_periods"`
	Severity         string           `gorm:"column:severity" json:"severity"`
	MetricsType      string           `gorm:"column:metrics_type" json:"metrics_type"`
	ConditionType    string           `gorm:"column:condition_type" json:"condition_type"`
	Thresholds       string           `gorm:"column:thresholds" json:"thresholds"`
	Unit             string           `gorm:"column:unit" json:"unit"`
	ConsecutiveCount uint32           `gorm:"column:consecutive_count" json:"consecutive_count"`
	Inhibit          bool             `gorm:"column:inhibit" json:"inhibit"`
	MetricName       string           `gorm:"column:metric_name" json:"metric_name"`
	Resources        []ResourceStatus `gorm:"column:resources" json:"resources"`
	CreateTime       time.Time        `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time        `gorm:"column:update_time" json:"update_time"`
	AlertStatus      string           `gorm:"column:alert_status"`
}

func AlertStatusToPb(alertStatus AlertStatus) *pb.AlertStatus {
	pbAlertStatus := pb.AlertStatus{}
	pbAlertStatus.RuleId = alertStatus.RuleId
	pbAlertStatus.RuleName = alertStatus.RuleName
	pbAlertStatus.MetricName = alertStatus.MetricName
	pbAlertStatus.Disabled = alertStatus.Disabled
	pbAlertStatus.MonitorPeriods = alertStatus.MonitorPeriods
	pbAlertStatus.Severity = alertStatus.Severity
	pbAlertStatus.MetricsType = alertStatus.MetricsType
	pbAlertStatus.ConditionType = alertStatus.ConditionType
	pbAlertStatus.Thresholds = alertStatus.Thresholds
	pbAlertStatus.Unit = alertStatus.Unit
	pbAlertStatus.ConsecutiveCount = alertStatus.ConsecutiveCount
	pbAlertStatus.Inhibit = alertStatus.Inhibit
	for _, resource := range alertStatus.Resources {
		pbResource := pb.ResourceStatus{}
		pbResource.ResourceName = resource.ResourceName
		pbResource.CurrentLevel = resource.CurrentLevel
		pbResource.PositiveCount = resource.PositiveCount
		pbResource.CumulatedSendCount = resource.CumulatedSendCount
		pbResource.NextResendInterval = resource.NextResendInterval
		pbResource.NextSendableTime = resource.NextSendableTime
		pbResource.AggregatedAlerts = resource.AggregatedAlerts

		pbAlertStatus.Resources = append(pbAlertStatus.Resources, &pbResource)
	}
	pbAlertStatus.CreateTime = pbutil.ToProtoTimestamp(alertStatus.CreateTime)
	pbAlertStatus.UpdateTime = pbutil.ToProtoTimestamp(alertStatus.UpdateTime)
	return &pbAlertStatus
}

func ParseAlsSet2PbSet(inAlss []AlertStatus) []*pb.AlertStatus {
	var pbAlss []*pb.AlertStatus
	for _, inAls := range inAlss {
		pbAls := AlertStatusToPb(inAls)
		pbAlss = append(pbAlss, pbAls)
	}
	return pbAlss
}

type AlertWrapper struct {
	ResourceMap string         `json:"resource_map"`
	RsFilter    ResourceFilter `json:"resource_filter"`
	Policy      Policy         `json:"policy"`
	Rules       []Rule         `json:"rules"`
	Action      Action         `json:"action"`
	Alert       Alert          `json:"alert"`
}

func RuleWrapperToPb(ruleWrapper *Rule) *pb.RuleWrapper {
	pbRuleWrapper := pb.RuleWrapper{}

	pbRuleWrapper.RuleName = ruleWrapper.RuleName
	pbRuleWrapper.Disabled = ruleWrapper.Disabled
	pbRuleWrapper.MonitorPeriods = ruleWrapper.MonitorPeriods
	pbRuleWrapper.Severity = ruleWrapper.Severity
	pbRuleWrapper.MetricsType = ruleWrapper.MetricsType
	pbRuleWrapper.ConditionType = ruleWrapper.ConditionType
	pbRuleWrapper.Thresholds = ruleWrapper.Thresholds
	pbRuleWrapper.Unit = ruleWrapper.Unit
	pbRuleWrapper.ConsecutiveCount = ruleWrapper.ConsecutiveCount
	pbRuleWrapper.Inhibit = ruleWrapper.Inhibit
	pbRuleWrapper.MetricId = ruleWrapper.MetricId

	return &pbRuleWrapper
}

func AlertWrapperToPb(alertWrapper *AlertWrapper) *pb.CreateAlertWrapperRequest {
	pbCreateAlertWrapperRequest := new(pb.CreateAlertWrapperRequest)

	pbCreateAlertWrapperRequest.ResourceMap = alertWrapper.ResourceMap

	pbCreateAlertWrapperRequest.RsFilter = new(pb.ResourceFilterWrapper)
	pbCreateAlertWrapperRequest.RsFilter.RsFilterName = alertWrapper.RsFilter.RsFilterName
	pbCreateAlertWrapperRequest.RsFilter.RsFilterParam = alertWrapper.RsFilter.RsFilterParam
	pbCreateAlertWrapperRequest.RsFilter.Status = alertWrapper.RsFilter.Status
	pbCreateAlertWrapperRequest.RsFilter.RsTypeId = alertWrapper.RsFilter.RsTypeId

	pbCreateAlertWrapperRequest.Policy = new(pb.PolicyWrapper)
	pbCreateAlertWrapperRequest.Policy.PolicyName = alertWrapper.Policy.PolicyName
	pbCreateAlertWrapperRequest.Policy.PolicyDescription = alertWrapper.Policy.PolicyDescription
	pbCreateAlertWrapperRequest.Policy.PolicyConfig = alertWrapper.Policy.PolicyConfig
	pbCreateAlertWrapperRequest.Policy.Creator = alertWrapper.Policy.Creator
	pbCreateAlertWrapperRequest.Policy.AvailableStartTime = alertWrapper.Policy.AvailableStartTime
	pbCreateAlertWrapperRequest.Policy.AvailableEndTime = alertWrapper.Policy.AvailableEndTime
	pbCreateAlertWrapperRequest.Policy.Language = alertWrapper.Policy.Language

	var rules []*pb.RuleWrapper
	for _, rule := range alertWrapper.Rules {
		pbRuleWrapper := RuleWrapperToPb(&rule)
		rules = append(rules, pbRuleWrapper)
	}
	pbCreateAlertWrapperRequest.Rules = rules

	pbCreateAlertWrapperRequest.Action = new(pb.ActionWrapper)
	pbCreateAlertWrapperRequest.Action.ActionName = alertWrapper.Action.ActionName
	pbCreateAlertWrapperRequest.Action.TriggerStatus = alertWrapper.Action.TriggerStatus
	pbCreateAlertWrapperRequest.Action.TriggerAction = alertWrapper.Action.TriggerAction
	pbCreateAlertWrapperRequest.Action.NfAddressListId = alertWrapper.Action.NfAddressListId

	pbCreateAlertWrapperRequest.Alert = new(pb.AlertWrapper)
	pbCreateAlertWrapperRequest.Alert.AlertName = alertWrapper.Alert.AlertName
	pbCreateAlertWrapperRequest.Alert.Disabled = alertWrapper.Alert.Disabled

	return pbCreateAlertWrapperRequest
}
