package client

// TaskSummary is a simplified representation of a task for notification rendering.
type TaskSummary struct {
	Name       string
	Status     string
	QD         string
	SignID     string
	QDTimeText string
}

// ToSummary converts a TaskList to a slice of TaskSummary.
func (tl *TaskList) ToSummary() []TaskSummary {
	if tl == nil {
		return nil
	}
	result := make([]TaskSummary, 0, len(tl.List.Items))
	for _, item := range tl.List.Items {
		statusText := "未签到"
		if item.QD == StatusSignSuccessfullyOk {
			statusText = "已签到"
		}
		result = append(result, TaskSummary{
			Name:       item.Name,
			Status:     statusText,
			QD:         item.QD,
			SignID:     item.SignID,
			QDTimeText: item.QDTimeText,
		})
	}
	return result
}
