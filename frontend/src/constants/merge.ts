// MergeTaskStatus mirrors backend model.MergeTask* constants.
export type MergeTaskStatus = 'pending' | 'processing' | 'succeeded' | 'failed'

export const MergeTaskStatusMap: Record<MergeTaskStatus, string> = {
  pending: '待执行',
  processing: '执行中',
  succeeded: '已完成',
  failed: '失败',
}

export const MergeTaskStatusTagType: Record<MergeTaskStatus, 'info' | 'warning' | 'success' | 'danger'> = {
  pending: 'info',
  processing: 'warning',
  succeeded: 'success',
  failed: 'danger',
}

// el-alert uses 'error' where el-tag uses 'danger'.
export const MergeTaskStatusAlertType: Record<MergeTaskStatus, 'info' | 'warning' | 'success' | 'error'> = {
  pending: 'info',
  processing: 'warning',
  succeeded: 'success',
  failed: 'error',
}
