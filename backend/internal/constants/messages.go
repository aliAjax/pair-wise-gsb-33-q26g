package constants

// Shared user-facing and backend messages. A wording change here may affect
// frontend hints, backend responses and log content simultaneously.
const (
	MsgRegisterOK           = "注册成功"
	MsgLoginOK              = "登录成功"
	MsgLogoutOK             = "退出登录成功"
	MsgPlantAddedToGarden   = "已加入我的花园"
	MsgPlantRemovedGarden   = "已从我的花园移除"
	MsgFavoriteAdded        = "收藏成功"
	MsgFavoriteRemoved      = "已取消收藏"
	MsgArticlePublished     = "文章发布成功"
	MsgArticleSaved         = "文章已保存"
	MsgReminderCreated      = "养护提醒已创建"
	MsgReminderDone         = "提醒已标记完成"
	MsgQuestionCreated      = "问题发布成功"
	MsgAnswerAdopted        = "已采纳该回答"
	MsgAnswerLiked          = "点赞成功"
	MsgProfileUpdated       = "资料已更新"
	MsgUploadOK             = "上传成功"
	MsgInvalidCredentials   = "用户名或密码错误"
	MsgUsernameTaken        = "用户名已存在"
	MsgEmailTaken           = "邮箱已被注册"
	MsgPlantMergeOK         = "品种合并完成"
	MsgPlantMergeProcessing = "相同合并任务正在执行中，请稍后查看结果"
	MsgPlantMergeExists     = "相同合并任务已完成"
	MsgPlantMergeFailed     = "品种合并失败，可重试从检查点恢复"
)
