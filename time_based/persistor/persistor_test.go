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
		timeSlots := "10:00,11:00"
		daysInWeek := "wed,fri"
		fileName := "hw.sh"
		content := []byte("echo 'hello world'")
		numberOfWeeks := 4

		mockFileDao.EXPECT().SaveFile(fileName, content)
		mockSettingsDao.EXPECT().SaveJob(jobName, timeSlots, daysInWeek, fileName, numberOfWeeks)

		persistor.SaveJob(jobName, timeSlots, daysInWeek, fileName, numberOfWeeks, content)
	})
}
