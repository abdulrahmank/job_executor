package time_based

import "github.com/abdulrahmank/job_executor/time_based/dao"

type Persistor struct {
	FileDao    dao.FileDao
	SettingDao dao.JobSettingDao
}

func (p *Persistor) SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int, jobFileContent []byte) {
	p.FileDao.SaveFile(fileName, jobFileContent)
	p.SettingDao.SaveJob(jobName, timeSlots, daysInWeek, fileName, numberOfWeeks)
}
