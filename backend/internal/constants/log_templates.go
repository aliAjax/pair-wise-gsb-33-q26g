package constants

// Central log message templates. Every handler/service/middleware must reference
// these templates instead of inline strings so that a business field change
// touches the log layer at the same time.
const (
	LogUserRegisterSuccess      = "user register success: username=%s"
	LogUserRegisterFailed       = "user register failed: username=%s"
	LogUserLoginSuccess         = "user login success: username=%s"
	LogUserLoginFailed          = "user login failed: username=%s"
	LogUserProfileUpdated       = "user profile updated: user_id=%d"
	LogUserAvatarUploaded       = "user avatar uploaded: user_id=%d"
	LogPlantCreateSuccess       = "plant species created: name=%s type=%s"
	LogPlantCreateFailed        = "plant species create failed: name=%s"
	LogPlantUpdateSuccess       = "plant species updated: id=%d"
	LogPlantUpdateFailed        = "plant species update failed: id=%d"
	LogPlantDeleteSuccess       = "plant species deleted: id=%d"
	LogPlantListSuccess         = "plant species list success: page=%d page_size=%d"
	LogPlantFavoriteSuccess     = "plant favorite success: plant_id=%d user_id=%d"
	LogPlantFavoriteFailed      = "plant favorite failed: plant_id=%d user_id=%d"
	LogArticleCreateSuccess     = "care article created: title=%s"
	LogArticleCreateFailed      = "care article create failed: title=%s"
	LogArticleUpdateSuccess     = "care article updated: id=%d"
	LogArticleDeleteSuccess     = "care article deleted: id=%d"
	LogArticleViewIncremented   = "care article view_count incremented: id=%d"
	LogPestCreateSuccess        = "disease pest created: name=%s"
	LogPestCreateFailed         = "disease pest create failed: name=%s"
	LogPestUpdateFailed         = "disease pest update failed: id=%d"
	LogPestSearchSuccess        = "disease pest search success: keyword=%s"
	LogReminderCreateSuccess    = "care reminder created: task_title=%s"
	LogReminderCreateFailed     = "care reminder create failed: task_title=%s"
	LogReminderStatusChanged    = "care reminder status changed: id=%d status=%s"
	LogGardenAddSuccess         = "garden item added: plant_id=%d user_id=%d"
	LogGardenAddFailed          = "garden item add failed: plant_id=%d user_id=%d"
	LogQuestionCreateSuccess    = "question created: title=%s"
	LogQuestionCreateFailed     = "question create failed: title=%s"
	LogAnswerCreateSuccess      = "answer created: question_id=%d"
	LogAnswerAdoptSuccess       = "answer adopted as best: id=%d"
	LogAnswerLikeSuccess        = "answer liked: id=%d"
	LogUploadSuccess            = "file upload success: path=%s"
	LogUploadFailed             = "file upload failed: filename=%s"
	LogRateLimited              = "request rate limited: path=%s ip=%s"
	LogRequestHandled           = "request handled: request_id=%s method=%s path=%s status=%d latency_ms=%d"
)

// LogTemplateCount returns the number of defined log templates (used by tests).
func LogTemplateCount() int {
	return 33
}
