package operations_test

import (
	"testing"

	"github.com/billiem/seren-management/pkg/operations"
)

func TestProgress(t *testing.T) {

	// test table
	tests := []struct {
		name            string
		numProcesses    int
		stepsPerProcess int
		expStepVal      float64
		expCompleteVal  float64
	}{
		{
			name:            "test",
			numProcesses:    2,
			stepsPerProcess: 10,
			expStepVal:      0.05,
			expCompleteVal:  0.5,
		},
		{
			name:            "test",
			numProcesses:    10,
			stepsPerProcess: 1,
			expStepVal:      0.1,
			expCompleteVal:  0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			P := operations.BuildProgress(tt.numProcesses, tt.stepsPerProcess)

			if P.Step(1) != tt.expStepVal {
				t.Errorf("expected %f, got %f", tt.expStepVal, P.Step(1))
			}

			if P.Complete(1) != tt.expCompleteVal {
				t.Errorf("expected %f, got %f", tt.expCompleteVal, P.Complete(1))
			}
		})
	}
}
