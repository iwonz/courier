package sshx

// AgentEndpoint returns an explicit SSH_AUTH_SOCK value or the native
// platform default when the platform defines one.
func AgentEndpoint(explicit string) string {
	if explicit != "" {
		return explicit
	}
	return defaultAgentEndpoint()
}
