package forms

// Add an error message to the map
type errors map[string][]string


// Add an error message to the map
func (e errors) Add(field, message string) {
	es := e[field]
	e[field] = append(es, message)
}

// Get the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

 
