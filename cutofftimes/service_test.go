package cutofftimes

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/zamedic/go2hal/alert/mock_alert"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceservices/mock_gceSelenium"
	"github.com/zamedic/go2hal/remoteTelegramCommands/mock_remoteTelegramCommands"
	"github.com/zamedic/go2hal/remoteTelegramCommands"
	"context"
	"github.com/zamedic/go2hal/halSelenium/mock_selenium"
	"github.com/tebeka/selenium"
	"github.com/pkg/errors"
	gomock2 "github.com/zamedic/go2hal/gomock"
)

func TestService_DoCheck(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockAlert := mock_alert.NewMockService(ctrl)
	mockSelenium := mock_gceSelenium.NewMockService(ctrl)
	mockStore := NewMockStore(ctrl)

	mockClient := mock_remoteTelegramCommands.NewMockRemoteCommand_RegisterCommandClient(ctrl)
	mockRemoteCommand := mock_remoteTelegramCommands.NewMockRemoteCommandClient(ctrl)

	mockRemoteCommand.EXPECT().RegisterCommand(context.Background(), &remoteTelegramCommands.RemoteCommandRequest{Name: "DisableGCECutoffTimes", Description: "Disable GCE Cutoff time check"}).Return(mockClient, nil)
	mockRemoteCommand.EXPECT().RegisterCommand(context.Background(), &remoteTelegramCommands.RemoteCommandRequest{Name: "GCECutofftimes", Description: "Trigger GCE Cutofftime check"}).Return(mockClient, nil)

	mockDriver := mock_selenium.NewMockWebDriver(ctrl)
	mockWebElement := mock_selenium.NewMockWebElement(ctrl)

	mockSelenium.EXPECT().Driver().Times(17).Return(mockDriver)
	mockSelenium.EXPECT().WaitForWaitFor().Times(4)

	mockDriver.EXPECT().Wait(gomock.Any()).Times(4).Return(nil)
	mockDriver.EXPECT().FindElement(selenium.ByPartialLinkText,"Service Options").Return(mockWebElement,nil)
	mockDriver.EXPECT().FindElement(selenium.ByPartialLinkText,"INWARD SERVICE OPTIONS").Return(mockWebElement,nil)
	mockClient.EXPECT().Recv().Return(nil,errors.New("Out of scope for this test"))

	mockWebElement.EXPECT().Click().Times(3).Return(nil)

	mockDriver.EXPECT().FindElement(selenium.ByXPATH,"//table[@id='TABLEINWARDSERVICES']/tbody/tr[1]/td[2]").Times(2).Return(mockWebElement,nil)
	mockWebElement.EXPECT().Text().Times(2).Return("Service A",nil)

	mockWebElement2 := mock_selenium.NewMockWebElement(ctrl)
	mockDriver.EXPECT().FindElement(selenium.ByXPATH,"//table[@id='TABLEINWARDSERVICES']/tbody/tr[1]/td[3]").Times(2).Return(mockWebElement2,nil)
	mockWebElement2.EXPECT().Text().Times(2).Return("SUB SERVICE A",nil)

	mockWebElement3 := mock_selenium.NewMockWebElement(ctrl)
	mockDriver.EXPECT().FindElement(selenium.ByXPATH,"//table[@id='TABLEINWARDSERVICES']/tbody/tr[1]/td[4]").Times(2).Return(mockWebElement3,nil)
	mockWebElement3.EXPECT().Text().Times(2).Return("DEST A",nil)

	mockWebElement13 := mock_selenium.NewMockWebElement(ctrl)
	mockDriver.EXPECT().FindElement(selenium.ByXPATH,"//table[@id='TABLEINWARDSERVICES']/tbody/tr[1]/td[13]").Times(2).Return(mockWebElement13,nil)
	mockWebElement13.EXPECT().Text().Times(2).Return("STATUS",nil)

	mockDriver.EXPECT().FindElement(selenium.ByXPATH,"//table[@id='TABLEINWARDSERVICES']/tbody/tr[2]/td[2]").Times(2).Return(nil,errors.New("Loop Breakout"))

	mockDriver.EXPECT().FindElement(selenium.ByPartialLinkText,"2").Return(mockWebElement,nil)

	mockStore.EXPECT().cutoffExists("Service ASUB SERVICE A","DEST A").Times(2).Return(true)
	mockStore.EXPECT().isInStartOfDay("Service ASUB SERVICE A","DEST A").Times(2).Return(true)

	mockSelenium.EXPECT().HandleSeleniumError(false,gomock2.ErrorMsgMatches(errors.New("invalid status for service Service ASUB SERVICE A, sub service DEST A, status STATUS\ninvalid status for service Service ASUB SERVICE A, sub service DEST A, status STATUS\n")))

	s := NewService(mockStore,mockSelenium,mockRemoteCommand,mockAlert)

	s.DoCheck(true)

}