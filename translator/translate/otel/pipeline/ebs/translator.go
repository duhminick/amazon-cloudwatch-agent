// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package ebs

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/pipeline"

	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/exporter/awscloudwatch"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/extension/agenthealth"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/processor/ec2taggerprocessor"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/processor/filterprocessor"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/receiver/awsebsnvme"
)

const (
	ebsPipelineName = common.PipelineNameEbs
)

type translator struct {
	pipelineName string
	common.DestinationProvider
}

var (
	baseKey = common.ConfigKey(common.MetricsKey, common.MetricsCollectedKey)
	diskioKey = common.ConfigKey(baseKey, common.DiskIOKey)
	_ common.PipelineTranslator = (*translator)(nil)
)

func NewTranslator(
	opts ...common.TranslatorOption,
) common.PipelineTranslator {
	t := &translator{pipelineName: ebsPipelineName}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *translator) ID() pipeline.ID {
	return pipeline.NewIDWithName(pipeline.SignalMetrics, t.pipelineName)
}

// Translate creates a pipeline for EBS metrics
func (t *translator) Translate(conf *confmap.Conf) (*common.ComponentTranslators, error) {
	if conf == nil || !conf.IsSet(diskioKey) {
		return nil, &common.MissingKeyError{ID: t.ID(), JsonKey: common.DiskIOKey}
	}

	if measurement := common.GetArray[any](conf, common.ConfigKey(diskioKey, common.MeasurementKey)); measurement != nil && len(measurement) == 0 {
		return nil, &common.MissingKeyError{ID: t.ID(), JsonKey: common.ConfigKey(diskioKey, common.MeasurementKey)}
	}


	translators := common.ComponentTranslators{
		Receivers: common.NewTranslatorMap(
			awsebsnvme.NewTranslator(common.WithName(t.pipelineName)),
		),
		Processors: common.NewTranslatorMap(
			filterprocessor.NewTranslator(common.WithName(t.pipelineName)),
		),
		Exporters: common.NewTranslatorMap[component.Config, component.ID](),
		Extensions: common.NewTranslatorMap(
			agenthealth.NewTranslator(agenthealth.MetricsName, []string{agenthealth.OperationPutMetricData}),
			agenthealth.NewTranslatorWithStatusCode(agenthealth.StatusCodeName, nil, true),
		),
	}

	// 3. Metric transform processor
	// translators.Processors.Set(metricstransformprocessor.NewTranslator(
	// 	metricstransformprocessor.WithPipelineName(PipelineName),
	// ))

	// 4. EC2 tagger processor
	if conf.IsSet(common.ConfigKey(common.MetricsKey, common.AppendDimensionsKey)) {
		translators.Processors.Set(ec2taggerprocessor.NewTranslator())
	}

	// 5. Cumulative to delta processor
	// translators.Processors.Set(cumulativetodeltaprocessor.NewTranslator(
	// 	common.WithName(PipelineName),
	// 	cumulativetodeltaprocessor.WithDefaultKeys(),
	// ))

	switch t.Destination() {
	case common.DefaultDestination, common.CloudWatchKey:
		translators.Exporters.Set(awscloudwatch.NewTranslator())
	default:
		return nil, fmt.Errorf("pipeline (%s) does not support destination (%s) in configuration", t.pipelineName, t.Destination())
	}

	return &translators, nil
}
