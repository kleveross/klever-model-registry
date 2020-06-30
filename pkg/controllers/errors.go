package controllers

const (
	// Success is the exit code of modeljob success
	Success = 0
	// ErrORMBLogin is the exit code of ormb login error
	ErrORMBLogin = 10000
	// ErrORMBPullModel is the exit code of ormb pull error
	ErrORMBPullModel = 10001
	// ErrORMBExportModel is the exit code of ormb export error
	ErrORMBExportModel = 10002
	// ErrRunTask is the exit code of run task error
	ErrRunTask = 10003
	// ErrORMBSaveModel is the exit code of ormb save error
	ErrORMBSaveModel = 10004
	// ErrORMBPushModel is the exit code of ormb push error
	ErrORMBPushModel = 10005
)
