package persistor

import "github.com/abdulrahmank/job_executor/time_based/dao"

type Persistor struct {
	FileDao    dao.FileDao
	SettingDao dao.JobSettingDao
}

func (p *Persistor) SaveJob(jobName, fileName string, jobFileContent []byte) {
	p.FileDao.SaveFile(fileName, jobFileContent)
	p.SettingDao.SaveJob(jobName, fileName)
}

func (p *Persistor) DeleteJob(jobName string) {
	fileName := p.SettingDao.GetFileName(jobName)
	p.FileDao.DeleteFile(fileName)
	p.SettingDao.DeleteJob(jobName)
}

func (p *Persistor) SaveEvenBasedJob(jobName, fileName, eventName string, jobFileContent []byte) {
	p.FileDao.SaveFile(fileName, jobFileContent)
	p.SettingDao.SaveEventBasedJob(jobName, fileName, eventName)
}

func (p *Persistor) ConfigureTimeBasedJob(jobName, timeSlots, daysInWeek string, numberOfWeeks int) {
	p.SettingDao.SaveTimedJob(jobName, timeSlots, daysInWeek, numberOfWeeks)
}
