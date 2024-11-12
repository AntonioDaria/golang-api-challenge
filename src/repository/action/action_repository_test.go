package action

import (
	"reflect"
	"testing"
	"time"

	"github.com/AntonioDaria/surfe/src/models"
)

// crate a func to load the actions from the json file
func loadActionRepo(t *testing.T) *RepositoryImpl {
	t.Helper()

	actionRepo, err := NewActionRepo("../data/actions.json")
	if err != nil {
		t.Fatalf("failed to create action repository: %v", err)
	}

	return actionRepo
}

func Test_Count_Actions_By_User_ID(t *testing.T) {
	// Arrange
	actionRepo := loadActionRepo(t)

	// Act
	actions := actionRepo.CountActionsByUserID(1)

	// Assert
	if actions != 49 {
		t.Fatalf("expected actions count to be 49, got %d", actions)
	}
}

func Test_User_Exists(t *testing.T) {
	// Arrange
	actionRepo := loadActionRepo(t)

	// Act
	userExists := actionRepo.UserExists(1)

	// Assert
	if !userExists {
		t.Fatal("expected user to exist")
	}
}

func Test_User_Does_Not_Exist(t *testing.T) {
	// Arrange
	actionRepo := loadActionRepo(t)

	// Act
	userExists := actionRepo.UserExists(1000)

	// Assert
	if userExists {
		t.Fatal("expected user to not exist")
	}
}

func Test_Get_Sorted_Actions(t *testing.T) {
	// Arrange
	actionRepo := loadActionRepo(t)

	// Act
	actions := actionRepo.GetSortedActions()

	// Assert
	if len(actions) != 22938 {
		t.Fatalf("expected actions count to be 22938, got %d", len(actions))
	}

	// assert that actions are sorted by timestamp within each user
	for i := 0; i < len(actions)-1; i++ {
		if actions[i].UserID == actions[i+1].UserID && actions[i].CreatedAt.After(actions[i+1].CreatedAt) {
			t.Fatalf("actions are not sorted by timestamp within each user")
		}
	}
}

func Test_Get_All_Actions(t *testing.T) {
	// Arrange
	actionRepo := loadActionRepo(t)

	// Act
	actions := actionRepo.GetAllActions()

	// Assert
	if len(actions) != 22938 {
		t.Fatalf("expected actions count to be 22938, got %d", len(actions))
	}
}

func TestRepositoryImpl_GetSortedActions(t *testing.T) {
	// Define the time format and parse timestamps
	timeFormat := "2006-01-02T15:04:05Z"
	time1, _ := time.Parse(timeFormat, "2021-01-01T00:00:00Z")
	time2, _ := time.Parse(timeFormat, "2021-01-01T00:00:01Z")
	time3, _ := time.Parse(timeFormat, "2021-01-01T00:00:02Z")
	time4, _ := time.Parse(timeFormat, "2021-01-01T00:00:03Z")

	type fields struct {
		actions []models.Action
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.Action
	}{
		{
			name: "Test GetSortedActions",
			fields: fields{
				actions: []models.Action{
					{
						ID:        1,
						UserID:    1,
						Type:      models.ActionTypeAddContact,
						CreatedAt: time1,
					},
					{
						ID:        4,
						UserID:    2,
						Type:      models.ActionTypeAddContact,
						CreatedAt: time4,
					},
					{
						ID:        3,
						UserID:    2,
						Type:      models.ActionTypeAddContact,
						CreatedAt: time3,
					},

					{
						ID:        2,
						UserID:    1,
						Type:      models.ActionTypeAddContact,
						CreatedAt: time2,
					},
				},
			},
			want: []models.Action{
				{
					ID:        1,
					UserID:    1,
					Type:      models.ActionTypeAddContact,
					CreatedAt: time1,
				},
				{
					ID:        2,
					UserID:    1,
					Type:      models.ActionTypeAddContact,
					CreatedAt: time2,
				},
				{
					ID:        3,
					UserID:    2,
					Type:      models.ActionTypeAddContact,
					CreatedAt: time3,
				},
				{
					ID:        4,
					UserID:    2,
					Type:      models.ActionTypeAddContact,
					CreatedAt: time4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepositoryImpl{
				Actions: tt.fields.actions,
			}
			if got := r.GetSortedActions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RepositoryImpl.GetSortedActions() = %v, want %v", got, tt.want)
			}
		})
	}
}
