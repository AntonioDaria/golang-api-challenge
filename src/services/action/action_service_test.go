package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/AntonioDaria/surfe/src/models"
	act_type "github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/repository/action"
	"github.com/AntonioDaria/surfe/src/repository/action/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CountActionsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	actionRepo := mock.NewMockRepository(ctrl)
	actionService := NewActionService(actionRepo)

	// Define expected behavior for UserExists
	actionRepo.EXPECT().UserExists(1).Return(true)

	// Define expected behavior for CountActionsByUserID
	actionRepo.EXPECT().CountActionsByUserID(1).Return(2)

	// Act
	count, err := actionService.GetActionCountByUserID(1)
	assert.NoError(t, err)

	// Assert
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func Test_CountActionsByUserID_User_Not_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	actionRepo := mock.NewMockRepository(ctrl)
	actionService := NewActionService(actionRepo)

	// Define expected behavior for UserExists
	actionRepo.EXPECT().UserExists(1).Return(false)

	// Act
	count, err := actionService.GetActionCountByUserID(1)
	assert.Error(t, err)
	assert.Equal(t, action.ErrUserNotFound, err)
	assert.Equal(t, 0, count)
}

func TestServiceImpl_GetNextActionProbabilities(t *testing.T) {
	// Define timestamps for test actions
	time1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2021, 1, 1, 0, 0, 1, 0, time.UTC)
	time3 := time.Date(2021, 1, 1, 0, 0, 2, 0, time.UTC)

	type fields struct {
		actionRepo *action.RepositoryImpl // Using RepositoryImpl directly
	}
	type args struct {
		actionType act_type.ActionType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[act_type.ActionType]float64
	}{
		{
			name: "Basic Case - Single User Sequence",
			fields: fields{
				actionRepo: &action.RepositoryImpl{
					Actions: []models.Action{
						{ID: 1, UserID: 1, Type: act_type.ActionTypeAddContact, CreatedAt: time1},
						{ID: 2, UserID: 1, Type: act_type.ActionTypeViewContacts, CreatedAt: time2},
					},
				},
			},
			args: args{actionType: act_type.ActionTypeAddContact},
			want: map[act_type.ActionType]float64{
				act_type.ActionTypeViewContacts: 1.0,
			},
		},
		{
			name: "Multiple Next Actions",
			fields: fields{
				actionRepo: &action.RepositoryImpl{
					Actions: []models.Action{
						{ID: 1, UserID: 1, Type: act_type.ActionTypeAddContact, CreatedAt: time1},
						{ID: 2, UserID: 1, Type: act_type.ActionTypeViewContacts, CreatedAt: time2},
						{ID: 3, UserID: 1, Type: act_type.ActionTypeEditContact, CreatedAt: time3},
					},
				},
			},
			args: args{actionType: act_type.ActionTypeAddContact},
			want: map[act_type.ActionType]float64{
				act_type.ActionTypeViewContacts: 0.5,
				act_type.ActionTypeEditContact:  0.5,
			},
		},
		{
			name: "No Next Actions",
			fields: fields{
				actionRepo: &action.RepositoryImpl{
					Actions: []models.Action{
						{ID: 1, UserID: 1, Type: act_type.ActionTypeAddContact, CreatedAt: time1},
					},
				},
			},
			args: args{actionType: act_type.ActionTypeAddContact},
			want: map[act_type.ActionType]float64{}, // No next actions
		},
		{
			name: "Action Type Not Found",
			fields: fields{
				actionRepo: &action.RepositoryImpl{
					Actions: []models.Action{
						{ID: 1, UserID: 1, Type: act_type.ActionTypeViewContacts, CreatedAt: time1},
					},
				},
			},
			args: args{actionType: act_type.ActionTypeAddContact},
			want: map[act_type.ActionType]float64{}, // Specified action type doesn't exist
		},
		{
			name: "Multiple Users - Independent Sequences",
			fields: fields{
				actionRepo: &action.RepositoryImpl{
					Actions: []models.Action{
						{ID: 1, UserID: 1, Type: act_type.ActionTypeAddContact, CreatedAt: time1},
						{ID: 2, UserID: 1, Type: act_type.ActionTypeViewContacts, CreatedAt: time2},
						{ID: 3, UserID: 2, Type: act_type.ActionTypeAddContact, CreatedAt: time1},
						{ID: 4, UserID: 2, Type: act_type.ActionTypeReferUser, CreatedAt: time2},
					},
				},
			},
			args: args{actionType: act_type.ActionTypeAddContact},
			want: map[act_type.ActionType]float64{
				act_type.ActionTypeViewContacts: 0.5,
				act_type.ActionTypeReferUser:    0.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceImpl{
				actionRepo: tt.fields.actionRepo,
			}
			if got := s.GetNextActionProbabilities(tt.args.actionType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceImpl.GetNextActionProbabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImpl_GetReferralIndex(t *testing.T) {
	// Mock data to create a referral chain
	actionRepo := &action.RepositoryImpl{
		Actions: []models.Action{
			{ID: 1, UserID: 1, Type: act_type.ActionTypeReferUser, TargetUser: 2},
			{ID: 2, UserID: 1, Type: act_type.ActionTypeReferUser, TargetUser: 3},
			{ID: 3, UserID: 2, Type: act_type.ActionTypeReferUser, TargetUser: 4},
			{ID: 4, UserID: 3, Type: act_type.ActionTypeReferUser, TargetUser: 5},
			// Expected: User 1 should have a referral index of 4 (2, 3, 4, 5)
		},
	}

	actionService := &ServiceImpl{
		actionRepo: actionRepo,
	}

	// Call the service function directly
	referralIndex := actionService.GetReferralIndex()

	// Assert expected referral indices
	assert.Equal(t, 4, referralIndex[1]) // User 1 referred 2, 3, 4, 5
	assert.Equal(t, 1, referralIndex[2]) // User 2 referred 4
	assert.Equal(t, 1, referralIndex[3]) // User 3 referred 5
	assert.Equal(t, 0, referralIndex[4]) // User 4 has no referrals
	assert.Equal(t, 0, referralIndex[5]) // User 5 has no referrals
}

func TestServiceImpl_GetReferralIndex_No_Referrals(t *testing.T) {
	// Mock data with no referral actions
	actionRepo := &action.RepositoryImpl{
		Actions: []models.Action{
			{ID: 1, UserID: 1, Type: act_type.ActionTypeAddContact},
			{ID: 2, UserID: 2, Type: act_type.ActionTypeViewContacts},
		},
	}

	actionService := &ServiceImpl{
		actionRepo: actionRepo,
	}

	// Call the service function directly
	referralIndex := actionService.GetReferralIndex()

	// Assert no referrals
	assert.Equal(t, 0, referralIndex[1])
	assert.Equal(t, 0, referralIndex[2])
}

func TestServiceImpl_GetReferralIndex_Circular_Referral(t *testing.T) {
	// Mock data to create a circular referral chain
	actionRepo := &action.RepositoryImpl{
		Actions: []models.Action{
			{ID: 1, UserID: 1, Type: act_type.ActionTypeReferUser, TargetUser: 2},
			{ID: 2, UserID: 2, Type: act_type.ActionTypeReferUser, TargetUser: 3},
			{ID: 3, UserID: 3, Type: act_type.ActionTypeReferUser, TargetUser: 1},
			// Expected: Circular referral chain, all users have a referral index of 1
		},
	}

	actionService := &ServiceImpl{
		actionRepo: actionRepo,
	}

	// Call the service function directly
	referralIndex := actionService.GetReferralIndex()

	// Assert expected referral indices for a circular referral chain
	assert.Equal(t, 2, referralIndex[1]) // User 1 has 2 indirect referrals (2 and 3)
	assert.Equal(t, 2, referralIndex[2]) // User 2 has 2 indirect referrals (3 and 1)
	assert.Equal(t, 2, referralIndex[3]) // User 3 has 2 indirect referrals (1 and 2)

}
