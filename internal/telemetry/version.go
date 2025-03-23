package telemetry

func Version() string {
	return "1.0.0"
}

func SemVersion() string {
	return "semver:" + Version()
}
