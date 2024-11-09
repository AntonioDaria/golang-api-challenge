package action

import (
	"testing"
)

func Test_Count_Actions_By_User_ID(t *testing.T) {
	// Arrange
	actionRepo, err := NewActionRepo("../data/actions.json")
	if err != nil {
		t.Fatalf("failed to create action repository: %v", err)
	}

	// Act
	actions, err := actionRepo.CountActionsByUserID(1)

	// Assert
	if err != nil {
		t.Fatalf("failed to count actions by user ID: %v", err)
	}

	if actions != 49 {
		t.Fatalf("expected actions count to be 49, got %d", actions)
	}
}
