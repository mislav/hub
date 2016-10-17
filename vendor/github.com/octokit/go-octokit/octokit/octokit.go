package octokit

const (
	gitHubAPIURL     = "https://api.github.com"
	userAgent        = "Octokit Go " + version
	version          = "0.3.0"
	defaultMediaType = "application/vnd.github.v3+json;charset=utf-8"
	diffMediaType    = "application/vnd.github.v3.diff;charset=utf-8"
	patchMediaType   = "application/vnd.github.v3.patch;charset=utf-8"
	textMediaType    = "text/plain;charset=utf-8"
)
