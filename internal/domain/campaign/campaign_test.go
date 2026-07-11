package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	createdBy = "test@test.com"
	content   = "Body Hi!"
	contacts  = []string{"email1@g.com", "email2@g.com"}
	fake      = faker.New()
)

func Test_NewCampaign_CreateCampaign(t *testing.T) {
	t.Parallel()

	campaign, err := NewCampaign(name, content, contacts, createdBy)

	assert.NoError(t, err)
	assert.Equal(t, name, campaign.Name)
	assert.Equal(t, content, campaign.Content)
	assert.Len(t, campaign.Contacts, len(contacts))
	assert.Equal(t, createdBy, campaign.CreatedBy)
}

func Test_NewCampaign_GeneratesIdentifierAndPendingStatus(t *testing.T) {
	t.Parallel()

	campaign, err := NewCampaign(name, content, contacts, createdBy)

	assert.NoError(t, err)
	assert.NotEmpty(t, campaign.ID)
	assert.Equal(t, Pending, campaign.Status)
}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {
	t.Parallel()

	now := time.Now().Add(-time.Second)
	campaign, err := NewCampaign(name, content, contacts, createdBy)

	assert.NoError(t, err)
	assert.True(t, campaign.CreatedOn.After(now))
}

func Test_NewCampaign_ValidationErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		campaign    string
		contentText string
		emails      []string
		createdBy   string
		expectedErr string
	}{
		{
			name:        "name min",
			campaign:    "",
			contentText: content,
			emails:      contacts,
			createdBy:   createdBy,
			expectedErr: "Name is required with min 5",
		},
		{
			name:        "name max",
			campaign:    fake.Lorem().Text(30),
			contentText: content,
			emails:      contacts,
			createdBy:   createdBy,
			expectedErr: "Name is required with max 24",
		},
		{
			name:        "content min",
			campaign:    name,
			contentText: "",
			emails:      contacts,
			createdBy:   createdBy,
			expectedErr: "Content is required with min 5",
		},
		{
			name:        "content max",
			campaign:    name,
			contentText: fake.Lorem().Text(1040),
			emails:      contacts,
			createdBy:   createdBy,
			expectedErr: "Content is required with max 1024",
		},
		{
			name:        "contacts min",
			campaign:    name,
			contentText: content,
			emails:      nil,
			createdBy:   createdBy,
			expectedErr: "Contacts is required with min 1",
		},
		{
			name:        "invalid email",
			campaign:    name,
			contentText: content,
			emails:      []string{"email invalid"},
			createdBy:   createdBy,
			expectedErr: "Email is invalid",
		},
		{
			name:        "created by invalid",
			campaign:    name,
			contentText: content,
			emails:      contacts,
			createdBy:   "",
			expectedErr: "CreatedBy is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCampaign(tc.campaign, tc.contentText, tc.emails, tc.createdBy)
			assert.EqualError(t, err, tc.expectedErr)
		})
	}
}

func Test_Campaign_StatusTransitionsUpdateTimestamp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		action         func(*Campaign)
		expectedStatus string
	}{
		{name: "done", action: func(c *Campaign) { c.Done() }, expectedStatus: Done},
		{name: "cancel", action: func(c *Campaign) { c.Cancel() }, expectedStatus: Canceled},
		{name: "delete", action: func(c *Campaign) { c.Delete() }, expectedStatus: Deleted},
		{name: "fail", action: func(c *Campaign) { c.Fail() }, expectedStatus: Fail},
		{name: "start", action: func(c *Campaign) { c.Started() }, expectedStatus: Started},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			campaign := &Campaign{Status: Pending, UpdatedOn: time.Time{}}
			before := time.Now().Add(-time.Second)

			tc.action(campaign)

			assert.Equal(t, tc.expectedStatus, campaign.Status)
			assert.True(t, campaign.UpdatedOn.After(before))
		})
	}
}
