package routes

const (
	//User Endpoints

	GenericUserSignup = "/signup"
	GenericUserLogin  = "/login"
	ResumeUpload      = "/uploadResume"

	//Admin EndPoints
	CreatrJobOpeningRoute   = "/admin/job"
	GetJobInfoByJobIdRoute  = "/admin/job/{job_id}"
	FetchAllUsersRoute      = "/admin/applicants"
	FetchApplicantDataRoute = "/admin/applicant/{applicant_id}"

	//Job Endpoints
	FetchingJobOpeningsRoute = "/jobs"
	ApplyParticularJobRoute  = "/jobs/apply/{job_id}"
)
