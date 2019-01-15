package persistor

import (
	"github.com/abdulrahmank/job_executor/time_based/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestSaveJob(t *testing.T) {
	t.Run("Should save given job", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFileDao := mocks.NewMockFileDao(mockCtrl)
		mockSettingsDao := mocks.NewMockJobSettingDao(mockCtrl)

		persistor := &Persistor{FileDao: mockFileDao, SettingDao: mockSettingsDao}
		jobName := "hw"
		fileName := "hw.sh"
		content := []byte("echo 'hello world'")

		mockFileDao.EXPECT().SaveFile(fileName, content)
		mockSettingsDao.EXPECT().SaveJob(jobName, fileName)

		persistor.SaveJob(jobName, fileName, content)
	})
}

func TestPersistor_DeleteJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileDao := mocks.NewMockFileDao(mockCtrl)
	mockSettingsDao := mocks.NewMockJobSettingDao(mockCtrl)

	jobName := "job"
	fileName := "file"
	mockSettingsDao.EXPECT().GetFileName(jobName).Return(fileName)
	mockSettingsDao.EXPECT().DeleteJob(jobName)
	mockFileDao.EXPECT().DeleteFile(fileName)

	persistor := &Persistor{FileDao: mockFileDao, SettingDao: mockSettingsDao}

	persistor.DeleteJob(jobName)
}

func TestPersistor_ConfigureTimeBasedJob(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileDao := mocks.NewMockFileDao(mockCtrl)
	mockSettingsDao := mocks.NewMockJobSettingDao(mockCtrl)

	persistor := &Persistor{FileDao: mockFileDao, SettingDao: mockSettingsDao}
	jobName := "hw"
	timeSlots := "10:00,11:00"
	daysInWeek := "wed,fri"
	numberOfWeeks := 4

	mockSettingsDao.EXPECT().SaveTimedJob(jobName, timeSlots, daysInWeek, numberOfWeeks)

	persistor.ConfigureTimeBasedJob(jobName, timeSlots, daysInWeek, numberOfWeeks)

}

func TestPersistor_SaveEvenBasedJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileDao := mocks.NewMockFileDao(mockCtrl)
	mockSettingsDao := mocks.NewMockJobSettingDao(mockCtrl)

	persistor := &Persistor{FileDao: mockFileDao, SettingDao: mockSettingsDao}
	jobName := "hw"
	eventName := "event"

	mockSettingsDao.EXPECT().SaveEventBasedJob(jobName, eventName)

	persistor.ConfigureEventBasedJob(jobName, eventName)
}
