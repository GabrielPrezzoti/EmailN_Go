package campaign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCampaign(t *testing.T) {
	assert := assert.New(t)                                           // arrange
	name := "campaing X"
	content := "Body"                                                   
	contacts := []string{"email1@g.com", "email2@g.com"}

	campaign := NewCampaign(name, content, contacts)               //act

	assert.Equal(campaign.ID, "1")                                  // assert
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
}
