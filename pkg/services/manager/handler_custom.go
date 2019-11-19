package manager

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	. "kubesphere.io/alert/pkg/pb"
	rs "kubesphere.io/alert/pkg/services/manager/resource_control"
)

//0.Alert
//********************************************************************************************************
func (s *Server) DescribeAlertsWithResource(ctx context.Context, req *DescribeAlertsWithResourceRequest) (*DescribeAlertsWithResourceResponse, error) {
	als, alCnt, err := rs.DescribeAlertsWithResource(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alerts With Resource, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	alPbSet := models.ParseAlSet2PbSet(als)
	res := &DescribeAlertsWithResourceResponse{
		Total:    uint32(alCnt),
		AlertSet: alPbSet,
	}

	logger.Debug(ctx, "Describe Alerts With Resource successfully, Alerts=[%+v].", res)
	return res, nil
}

func (s *Server) DescribeAlertDetails(ctx context.Context, req *DescribeAlertDetailsRequest) (*DescribeAlertDetailsResponse, error) {
	alds, aldCnt, err := rs.DescribeAlertDetails(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alert Details, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	aldPbSet := models.ParseAldSet2PbSet(alds)
	res := &DescribeAlertDetailsResponse{
		Total:          uint32(aldCnt),
		AlertdetailSet: aldPbSet,
	}

	logger.Debug(ctx, "Describe Alert Details successfully, Alert Details=[%+v].", res)
	return res, nil
}

func (s *Server) DescribeAlertStatus(ctx context.Context, req *DescribeAlertStatusRequest) (*DescribeAlertStatusResponse, error) {
	alss, alsCnt, err := rs.DescribeAlertStatus(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alert Status, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	alsPbSet := models.ParseAlsSet2PbSet(alss)
	res := &DescribeAlertStatusResponse{
		Total:          uint32(alsCnt),
		AlertstatusSet: alsPbSet,
	}

	logger.Debug(ctx, "Describe Alert Status successfully, Alert Status=[%+v].", res)
	return res, nil
}

func (s *Server) removeResourceFilter(ctx context.Context, rsFilterId string) {
	var reqDeleteResourceFilter = &DeleteResourceFiltersRequest{
		RsFilterId: strings.Split(rsFilterId, ","),
	}

	s.DeleteResourceFilters(ctx, reqDeleteResourceFilter)
}

func (s *Server) removePolicy(ctx context.Context, policyId string) {
	var reqDeletePolicies = &DeletePoliciesRequest{
		PolicyId: strings.Split(policyId, ","),
	}

	s.DeletePolicies(ctx, reqDeletePolicies)
}

func (s *Server) CreateAlertWrapper(ctx context.Context, req *CreateAlertWrapperRequest) (*CreateAlertResponse, error) {
	//1. Check Rules Length
	if len(req.Rules) == 0 {
		logger.Error(nil, "CreateAlertWrapper Rules Length error")
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, errors.New("CreateAlertWrapper Rules Length error"), gerr.ErrorCreateResourcesFailed)
	}

	//2. Check if resource_type match
	rsTypeIds := strings.Split(req.RsFilter.RsTypeId, ",")

	var reqRsType = &DescribeResourceTypesRequest{
		RsTypeId: rsTypeIds,
	}

	respRsType, err := s.DescribeResourceTypes(ctx, reqRsType)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper DescribeResourceTypes failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if respRsType.Total != 1 {
		logger.Error(nil, "CreateAlertWrapper resource type error")
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	resourceMap := make(map[string]string)
	err = json.Unmarshal([]byte(req.ResourceMap), &resourceMap)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Unmarshal ResourceMap error %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if respRsType.ResourceTypeSet[0].RsTypeName != resourceMap["rs_type_name"] {
		logger.Error(nil, "CreateAlertWrapper resource type mismatch")
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	//3. Check workspace namespace matches
	rsFilterURI := make(map[string]string)
	err = json.Unmarshal([]byte(req.RsFilter.RsFilterParam), &rsFilterURI)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Unmarshal rsFilterURI Error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	uriCorrect := true
	switch resourceMap["rs_type_name"] {
	case "cluster":
		break
	case "node":
		break
	case "workspace":
		if resourceMap["ws_name"] != rsFilterURI["ws_name"] {
			uriCorrect = false
		}
	case "namespace":
		if resourceMap["ns_name"] != rsFilterURI["ns_name"] {
			uriCorrect = false
		}
	case "workload":
		if resourceMap["ns_name"] != rsFilterURI["ns_name"] {
			uriCorrect = false
		}
	case "pod":
		if resourceMap["ns_name"] != rsFilterURI["ns_name"] {
			uriCorrect = false
		}
	case "container":
		if resourceMap["ns_name"] != rsFilterURI["ns_name"] {
			uriCorrect = false
		}
	}

	if !uriCorrect {
		logger.Error(nil, "CreateAlertWrapper uri mismatch")
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	//4. Check if same alert_name exists
	resourceSearch, _ := json.Marshal(resourceMap)
	alertNames := strings.Split(req.Alert.AlertName, ",")
	var reqCheck = &DescribeAlertsWithResourceRequest{
		ResourceSearch: string(resourceSearch),
		AlertName:      alertNames,
	}

	respCheck, err := s.DescribeAlertsWithResource(ctx, reqCheck)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper check alert name failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if respCheck.Total != 0 {
		logger.Error(nil, "CreateAlertWrapper alert name already exists")
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	//5. Create Resource Filter
	var reqRsFilter = &CreateResourceFilterRequest{
		RsFilterName:  req.RsFilter.RsFilterName,
		RsFilterParam: req.RsFilter.RsFilterParam,
		Status:        req.RsFilter.Status,
		RsTypeId:      req.RsFilter.RsTypeId,
	}

	respRsFilter, err := s.CreateResourceFilter(ctx, reqRsFilter)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Resource Filter failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debug(nil, "CreateAlertWrapper Resource Filter success: %+v", respRsFilter)
	rsFilterId := respRsFilter.RsFilterId

	//6. Create Policy
	var reqPolicy = &CreatePolicyRequest{
		PolicyName:         req.Policy.PolicyName,
		PolicyDescription:  req.Policy.PolicyDescription,
		PolicyConfig:       req.Policy.PolicyConfig,
		Creator:            req.Policy.Creator,
		AvailableStartTime: req.Policy.AvailableStartTime,
		AvailableEndTime:   req.Policy.AvailableEndTime,
		Language:           req.Policy.Language,
		RsTypeId:           req.RsFilter.RsTypeId,
	}

	respPolicy, err := s.CreatePolicy(ctx, reqPolicy)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Policy failed: %+v", err)

		s.removeResourceFilter(ctx, rsFilterId)

		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debug(nil, "CreateAlertWrapper Policy success: %+v", respPolicy)
	policyId := respPolicy.PolicyId

	//7. Create Action
	var reqAction = &CreateActionRequest{
		ActionName:      req.Action.ActionName,
		PolicyId:        policyId,
		NfAddressListId: req.Action.NfAddressListId,
	}

	respAction, err := s.CreateAction(ctx, reqAction)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Action failed: %+v", err)

		s.removeResourceFilter(ctx, rsFilterId)
		s.removePolicy(ctx, policyId)

		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debug(nil, "CreateAlertWrapper Action success: %+v", respAction)

	//8. Create Rules
	createRulesSuccess := true
	for _, rule := range req.Rules {
		var reqRule = &CreateRuleRequest{
			RuleName:         rule.RuleName,
			Disabled:         rule.Disabled,
			MonitorPeriods:   rule.MonitorPeriods,
			Severity:         rule.Severity,
			MetricsType:      rule.MetricsType,
			ConditionType:    rule.ConditionType,
			Thresholds:       rule.Thresholds,
			Unit:             rule.Unit,
			ConsecutiveCount: rule.ConsecutiveCount,
			Inhibit:          rule.Inhibit,
			PolicyId:         policyId,
			MetricId:         rule.MetricId,
		}

		_, err := s.CreateRule(ctx, reqRule)
		if err != nil {
			createRulesSuccess = false
			break
		}
	}

	if !createRulesSuccess {
		logger.Error(nil, "CreateAlertWrapper Rules failed")

		s.removeResourceFilter(ctx, rsFilterId)
		s.removePolicy(ctx, policyId)

		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debug(nil, "CreateAlertWrapper Rules success")

	//9. Create Alert
	var reqAlert = &CreateAlertRequest{
		AlertName:  req.Alert.AlertName,
		PolicyId:   policyId,
		RsFilterId: rsFilterId,
	}

	respAlert, err := s.CreateAlert(ctx, reqAlert)
	if err != nil {
		logger.Error(nil, "CreateAlertWrapper Alert failed: %+v", err)

		s.removeResourceFilter(ctx, rsFilterId)
		s.removePolicy(ctx, policyId)

		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debug(nil, "CreateAlertWrapper Alert success: %+v", respAlert)

	return respAlert, nil
}

//1.History
//********************************************************************************************************
func (s *Server) DescribeHistoryDetail(ctx context.Context, req *DescribeHistoryDetailRequest) (*DescribeHistoryDetailResponse, error) {
	hsds, hsdCnt, err := rs.DescribeHistoryDetail(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe History Detail, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	hsdPbSet := models.ParseHsdSet2PbSet(hsds)
	res := &DescribeHistoryDetailResponse{
		Total:            uint32(hsdCnt),
		HistorydetailSet: hsdPbSet,
	}

	logger.Debug(ctx, "Describe History Detail successfully, Histories=[%+v].", res)
	return res, nil
}
