package internal

type Site struct {
	domain string
	path string
	docRoot string
	isWordPress bool
	dbName string
	dbPass string
}

func (s *Site) GetDocRoot() string {
	if s.docRoot == "" {
		return s.domain
	}

	return s.domain + "/" + s.docRoot
}
