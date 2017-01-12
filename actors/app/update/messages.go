package update

// Check -> command.update.check
type Check struct {
}

// No -> event.update.no
type No struct {
}

// Available -> event.update.available
type Available struct {
}

// Download -> command.update.download
type Download struct {
}

// DownloadProgress -> event.update.download.progress
type DownloadProgress struct {
	Progress float64 `json:"progress"`
}

// DownloadComplete -> event.update.download.complete
type DownloadComplete struct {
}

// Install -> command.update.install
type Install struct {
}

// InstallComplete -> event.update.install.complete
type InstallComplete struct {
}

// InstallRestart -> event.update.install.restart
type InstallRestart struct {
}

// Fail -> event.update.fail
type Fail struct {
	Message string `json:"message"`
}
