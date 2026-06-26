package history

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tagline   string `json:"tagline"`
	Path      string `json:"path"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Message struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Phase     string `json:"phase"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
