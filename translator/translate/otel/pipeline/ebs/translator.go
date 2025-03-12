// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package ebs

import (
	"fmt"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/pipeline"

	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/exporter/awscloudwatch"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/extension/agenthealth"
)

const (
	// PipelineName is the name of the EBS metrics pipeline
	PipelineName = "metric/ebs"
)

var (
	metricsKey = common.ConfigKey(common.MetricsKey, common.MetricsCollectedKey)
)

type translator struct {
}

var _ common.PipelineTranslator = (*translator)(nil)

func NewTranslator() common.PipelineTranslator {
	return &translator{}
}

func (t *translator) ID() pipeline.ID {
	return pipeline.NewIDWithName(pipeline.SignalMetrics, PipelineName)
}

// Translate creates a pipeline for EBS metrics
func (t *translator) Translate(conf *confmap.Conf) (*common.ComponentTranslators, error) {
	if conf == nil || !conf.IsSet(metricsKey) {
		return nil, &common.MissingKeyError{ID: t.ID(), JsonKey: metricsKey}
	}

	log.Printf("D! Setting up EBS metrics pipeline")

	translators := common.ComponentTranslators{
		Receivers:  common.NewTranslatorMap[component.Config, component.ID](),
		Processors: common.NewTranslatorMap[component.Config, component.ID](),
		Exporters:  common.NewTranslatorMap(awscloudwatch.NewTranslator()),
		Extensions: common.NewTranslatorMap(
			agenthealth.NewTranslator(agenthealth.MetricsName, []string{agenthealth.OperationPutMetricData}),
			agenthealth.NewTranslatorWithStatusCode(agenthealth.StatusCodeName, nil, true),
		),
	}

	// For now, receivers are left empty as per requirements
	// In the future, specific EBS receivers would be added here

	// Check if we have any receivers to create a pipeline
	if translators.Receivers.Len() == 0 {
		log.Printf("D! No receivers configured for EBS metrics pipeline")
		// We're still creating the pipeline structure even without receivers
		// as requested in the requirements
	}

	return &translators, nil
}