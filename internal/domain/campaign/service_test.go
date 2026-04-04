package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internalErrors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Body Hi!",
		Emails:    []string{"test1@gmail.com"},
		CreatedBy: "test@test.com",
	}
	repositoryMock *internalmock.CampaignRepositoryMock
	service        = campaign.ServiceImp{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
}

func Test_Create_Campaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setUp()
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_RepositorySave(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("Error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)

	campaignReturned, _ := service.GetBy(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status, campaignReturned.Status)
	assert.Equal(campaign.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_When_Campaign_Dont_Exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignIdInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_status_not_equals_pendind(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)

	err := service.Delete(campaign.ID)

	assert.Equal("Campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problem(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body !", []string{"teste@email.com"}, newCampaign.CreatedBy)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignFound.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_has_success(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body !", []string{"teste@email.com"}, newCampaign.CreatedBy)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(nil)

	err := service.Delete(campaignFound.ID)

	assert.Nil(err)
}

func Test_Start_ReturnRecordNotFound_When_Campaign_Dont_Exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignIdInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_status_not_equals_pendind(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)

	err := service.Start(campaign.ID)

	assert.Equal("Campaign status invalid", err.Error())
}

func Test_Start_should_send_email(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	service.Repository = repositoryMock
	sentMail := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignSaved.ID {
			sentMail = true
		}
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)

	assert.True(sentMail)
}

func Test_Start_ReturnError_when_func_SendMail_fail(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	service.Repository = repositoryMock
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignSaved.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Pending}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignSaved.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)
	service.Repository = repositoryMock
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)

	assert.Equal(campaign.Done, campaignSaved.Status)
}
