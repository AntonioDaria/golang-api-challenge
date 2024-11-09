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

func Test_User_Exists(t *testing.T) {
	// Arrange
	actionRepo, err := NewActionRepo("../data/actions.json")
	if err != nil {
		t.Fatalf("failed to create action repository: %v", err)
	}

	// Act
	userExists := actionRepo.UserExists(1)

	// Assert
	if !userExists {
		t.Fatal("expected user to exist")
	}
}

func Test_User_Does_Not_Exist(t *testing.T) {
	// Arrange
	actionRepo, err := NewActionRepo("../data/actions.json")
	if err != nil {
		t.Fatalf("failed to create action repository: %v", err)
	}

	// Act
	userExists := actionRepo.UserExists(1000)

	// Assert
	if userExists {
		t.Fatal("expected user to not exist")
	}
}
