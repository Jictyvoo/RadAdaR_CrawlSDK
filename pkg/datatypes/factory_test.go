package datatypes

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestNewConstructorFactory(t *testing.T) {
	tests := []struct {
		name          string
		constructor   func() int
		expectedValue int
		expectPtr     bool
		expectError   bool
	}{
		{
			name:          "should return value from constructor 42",
			constructor:   func() int { return 42 },
			expectedValue: 42,
			expectPtr:     false,
			expectError:   false,
		},
		{
			name:          "should return pointer to value 100 from constructor",
			constructor:   func() int { return 100 },
			expectedValue: 100,
			expectPtr:     true,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewConstructorFactory(tt.constructor)

			if tt.expectPtr {
				resultPtr, err := factory.NewPtr()
				if tt.expectError && err == nil {
					t.Errorf("expected error, but got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if resultPtr == nil {
					t.Errorf("expected non-nil pointer result, got nil")
				} else if *resultPtr != tt.expectedValue {
					t.Errorf("expected %d, got %d", tt.expectedValue, *resultPtr)
				}
			} else {
				result, err := factory.New()
				if tt.expectError && err == nil {
					t.Errorf("expected error, but got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if result != tt.expectedValue {
					t.Errorf("expected %d, got %d", tt.expectedValue, result)
				}
			}
		})
	}
}

func TestNewInjectionFactory(t *testing.T) {
	var calledTimes uint16
	const injectionKey = "injectionKey"
	inj := remy.NewCycleDetectorInjector()
	remy.RegisterConstructor(inj, remy.Factory[uint16], func() (result uint16) {
		result = calledTimes
		calledTimes += 1
		return calledTimes
	}, injectionKey)
	factory := NewInjectionFactory[uint16](inj, injectionKey)

	for index := range uint16(42) {
		resultPtr, err := factory.NewPtr()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
			t.FailNow()
		} else if resultPtr == nil {
			t.Errorf("expected non-nil result pointer, got nil")
			t.FailNow()
		}

		var result uint16
		if result, err = factory.New(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if expectVal := (index + 1) << 1; result+*resultPtr != (expectVal<<1)-1 {
			t.Errorf("expected %d, got %d + %d", expectVal, result, *resultPtr)
		}
	}
}

func TestFactoryWithoutConstructorOrInjector(t *testing.T) {
	var factory GenericFactory[int]
	result, err := factory.NewPtr()
	if result != nil {
		t.Errorf("expected nil, got %d", result)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}

	result = new(int)
	if *result, err = factory.New(); *result != 0 {
		t.Errorf("expected 0, got %d", *result)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}
