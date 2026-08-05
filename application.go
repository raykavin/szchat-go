package szchat

import "context"

// ApplicationAutoMessage is an automatic message toggle (wrong option,
// service started/completed).
type ApplicationAutoMessage struct {
	Message string `json:"message"`
	Enable  bool   `json:"enable"`
}

// ApplicationConfig is the tenant-wide configuration object. FinalGrade and
// InitialGrade are typed any because the API returns them as strings from
// some endpoints and numbers from others.
type ApplicationConfig struct {
	ID               string                  `json:"_id,omitempty"`
	Pagination       int                     `json:"pagination,omitempty"`
	Timezone         string                  `json:"timezone,omitempty"`
	Language         string                  `json:"language,omitempty"`
	DDI              int                     `json:"ddi,omitempty"`
	ClosingCommand   string                  `json:"closingCommand,omitempty"`
	APIToken         string                  `json:"apiToken,omitempty"`
	CompletedService *ApplicationAutoMessage `json:"completed_service,omitempty"`
	StartedService   *ApplicationAutoMessage `json:"started_service,omitempty"`
	WrongOption      *ApplicationAutoMessage `json:"wrong_option,omitempty"`
	Distribuition    any                     `json:"distribuition,omitempty"`
	FinalGrade       any                     `json:"final_grade,omitempty"`
	InitialGrade     any                     `json:"initial_grade,omitempty"`
	StorageID        string                  `json:"storage_id,omitempty"`
	CreatedAt        string                  `json:"created_at,omitempty"`
	UpdatedAt        string                  `json:"updated_at,omitempty"`
}

// ApplicationConfigList is the payload returned by ApplicationAPI.Get.
type ApplicationConfigList struct {
	Configurations []ApplicationConfig `json:"configurations"`
}

// ApplicationUpdateRequest is the payload for ApplicationAPI.Update.
type ApplicationUpdateRequest struct {
	Pagination int    `json:"pagination,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	Language   string `json:"language,omitempty"`
	DDI        int    `json:"ddi,omitempty"`
	StorageID  string `json:"storage_id,omitempty"`
}

// ApplicationAttendanceGrades is the request/response payload for the
// /application/attendance grade range endpoints.
type ApplicationAttendanceGrades struct {
	InitialGrade any `json:"initial_grade,omitempty"`
	FinalGrade   any `json:"final_grade,omitempty"`
}

// ApplicationMessages is the request/response payload for the
// /application/messages automatic message endpoints.
type ApplicationMessages struct {
	WrongOption      *ApplicationAutoMessage `json:"wrong_option,omitempty"`
	StartedService   *ApplicationAutoMessage `json:"started_service,omitempty"`
	CompletedService *ApplicationAutoMessage `json:"completed_service,omitempty"`
}

// ApplicationAPI groups the /application endpoints.
type ApplicationAPI struct {
	client *Client
}

// Get returns the tenant configuration via GET /application.
func (a *ApplicationAPI) Get(ctx context.Context) (*ApplicationConfigList, error) {
	var resp ApplicationConfigList
	if err := a.client.get(ctx, "/application", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates the tenant configuration via PUT /application.
func (a *ApplicationAPI) Update(ctx context.Context, req ApplicationUpdateRequest) (*ApplicationConfig, error) {
	var resp ApplicationConfig
	if err := a.client.put(ctx, "/application", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAttendanceGrades returns the attendance grade range via
// GET /application/attendance.
func (a *ApplicationAPI) GetAttendanceGrades(ctx context.Context) (*ApplicationAttendanceGrades, error) {
	var resp ApplicationAttendanceGrades
	if err := a.client.get(ctx, "/application/attendance", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAttendanceGrades updates the attendance grade range via
// PUT /application/attendance.
func (a *ApplicationAPI) UpdateAttendanceGrades(ctx context.Context, req ApplicationAttendanceGrades) (*ApplicationConfig, error) {
	var resp ApplicationConfig
	if err := a.client.put(ctx, "/application/attendance", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAutoMessages returns the automatic messages configuration via
// GET /application/messages.
func (a *ApplicationAPI) GetAutoMessages(ctx context.Context) (*ApplicationMessages, error) {
	var resp ApplicationMessages
	if err := a.client.get(ctx, "/application/messages", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAutoMessages updates the automatic messages configuration via
// PUT /application/messages.
func (a *ApplicationAPI) UpdateAutoMessages(ctx context.Context, req ApplicationMessages) (*ApplicationConfig, error) {
	var resp ApplicationConfig
	if err := a.client.put(ctx, "/application/messages", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
