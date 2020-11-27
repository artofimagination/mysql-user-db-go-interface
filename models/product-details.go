package models

const (
	SupportClients = "support_clients"
	ClientUI       = "client_ui"
	ProjectUI      = "project_ui"
	Requires3D     = "requires_3d"
	HasTrial       = "has_trial"
	IsFree         = "is_free"
)

// Errors called in multiple places (for example in unittests).

var ErrProductDetailsNotInitialised = "Details map not initialised"
