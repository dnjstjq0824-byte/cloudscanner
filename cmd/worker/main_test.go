package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHarvestInterval guards the DEPLOY-documented 15-minute catalog/
// deal-ingestion harvest cadence.
func TestHarvestInterval(t *testing.T) {
	assert.Equal(t, 15*time.Minute, harvestInterval)
}
