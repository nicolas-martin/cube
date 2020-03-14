package repository

// Repository is an in memory database and
// holds the list of employee reports
type Repository struct {
	data map[string]interface{}
}

// NewRepository creates a new repository
func NewRepository() *Repository {
	d := make(map[string]interface{}, 0)
	return &Repository{data: d}
}

// Create creates a PayRollReport entry in the store
func (r *Repository) Create(id string, d interface{}) error {
	r.data[id] = d
	return nil
}

// Retrieve retrieves the PayrollReports given an ID
func (r *Repository) Retrieve(id string) (interface{}, error) {
	return r.data[id], nil

}
