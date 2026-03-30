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
	assert := assert.New(t) // arrange

	campaign, _ := NewCampaign(name, content, contacts, createdBy) //act

	assert.Equal(campaign.Name, name) // assert
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(campaign.CreatedBy, createdBy)
}

func Test_NewCampaign_IDIsNotNil(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.NotNil(campaign.ID)
}

func Test_NewCampaign_MustStatusStartPending(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(Pending, campaign.Status)
}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Greater(campaign.CreatedOn, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts, createdBy)

	assert.Equal("Name is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)
	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, createdBy)

	assert.Equal("Name is required with max 24", err.Error())
}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "", contacts, createdBy)

	assert.Equal("Content is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts, createdBy)

	assert.Equal("Content is required with max 1024", err.Error())
}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, nil, createdBy)

	assert.Equal("Contacts is required with min 1", err.Error())
}

func Test_NewCampaign_MustValidateContactsMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"email invalid"}, createdBy)

	assert.Equal("Email is invalid", err.Error())
}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, contacts, "")

	assert.Equal("CreatedBy is invalid", err.Error())
}
