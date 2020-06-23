package models

type ApplicationList struct {
	ApplicationID int32 `json:"application_id"`
	UserID        int32 `json:"user_id"`
	Position      int32 `json:"position"`
}

type ApplicationListInput struct {
	ApplicationID   int32 `json:"applicationId"`
	UserID          int32 `json:"userId"`
	Position        int32 `json:"position"`
	DesiredPosition int32 `json:"desiredPosition"`
}
