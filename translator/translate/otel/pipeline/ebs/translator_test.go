// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package ebs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/confmap"

	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

func TestTranslate(t *testing.T) {
	translator := NewTranslator()

	// Test with nil config
	_, err := translator.Translate(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing key in JSON")

	// Test with empty config
	conf := confmap.New()
	_, err = translator.Translate(conf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing key in JSON")

	// Test with metrics section
	conf = confmap.New()
	conf.Set(common.ConfigKey(common.MetricsKey, common.MetricsCollectedKey), map[string]interface{}{})
	result, err := translator.Translate(conf)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify that we have the expected components
	assert.Equal(t, 0, result.Receivers.Len(), "Should have no receivers")
	assert.Equal(t, 0, result.Processors.Len(), "Should have no processors")
	assert.Equal(t, 1, result.Exporters.Len(), "Should have one exporter")
	assert.Equal(t, 2, result.Extensions.Len(), "Should have two extensions")
}