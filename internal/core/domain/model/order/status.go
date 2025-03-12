package order

const (
	StatusCreated   Status = "Created"
	StatusAssigned  Status = "Assigned"
	StatusCompleted Status = "Completed"
)

type Status string
