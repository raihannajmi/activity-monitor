package handlers

import "testing"

func TestFloatingPositionArithmetic(t *testing.T) {
	// Simulate Trello position arithmetic logic
	calcNewPos := func(hasPrev, hasNext bool, posPrev, posNext float64) float64 {
		if hasPrev && hasNext {
			return (posPrev + posNext) / 2.0
		} else if hasPrev {
			return posPrev + 1000.0
		} else if hasNext {
			return posNext / 2.0
		}
		return 1000.0
	}

	// Test 1: Insert between two tasks
	if pos := calcNewPos(true, true, 1000.0, 2000.0); pos != 1500.0 {
		t.Errorf("Expected 1500.0, got %f", pos)
	}

	// Test 2: Insert after a task (append to end)
	if pos := calcNewPos(true, false, 2000.0, 0.0); pos != 3000.0 {
		t.Errorf("Expected 3000.0, got %f", pos)
	}

	// Test 3: Insert before a task (prepend to start)
	if pos := calcNewPos(false, true, 0.0, 1000.0); pos != 500.0 {
		t.Errorf("Expected 500.0, got %f", pos)
	}

	// Test 4: Insert into empty list
	if pos := calcNewPos(false, false, 0.0, 0.0); pos != 1000.0 {
		t.Errorf("Expected 1000.0, got %f", pos)
	}
}
