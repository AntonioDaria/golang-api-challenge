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
	actions := actionRepo.CountActionsByUserID(1)

	// Assert
	if actions != 49 {
		t.Fatalf("expected actions count to be 49, got %d", actions)
	}
}
