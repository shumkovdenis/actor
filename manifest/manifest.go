package manifest

var m *manifest

func init() {
	m = newManifest()
}

type manifest struct {
	Version string `mapstructure:"version" validate:"required"`
}

func newManifest() *manifest {
	return &manifest{}
}

func Version() string {
	return m.Version
}

func Get() *manifest {
	return m
}
