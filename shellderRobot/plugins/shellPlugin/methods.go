package shellPlugin

// GetUniqueId returns the unique-id of this command-container
// variable.
func (c *commandContainer) GetUniqueId() string {
	return c.result.UniqueId
}

// SetUniqueId method sets the specified string value as unique-id of
// this commandContainer variable.
func (c *commandContainer) SetUniqueId(value string) {
	c.result.UniqueId = value
}
