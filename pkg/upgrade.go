package via

// Upgrader provides a type for upgrading installed plans
type Upgrader struct {
	config   *Config
	upgrades Plans
}

// NewUpgrader creates and initializes a new Upgrader
func NewUpgrader(config *Config) *Upgrader {
	if config == nil {
		panic("config is nil")
	}
	return &Upgrader{
		config:   config,
		upgrades: Plans{},
	}
}

func isUpgradable(installed *Plan, git *Plan) bool {
	if installed.Cid != git.Cid {
		return true
	}
	return false
}

// Check will compared all installed plans against the git repository
// of plans. And returns a strings slice of plan names that can be
// upgraded
func (u *Upgrader) Check() ([]string, error) {
	files, err := u.config.DB.InstalledFiles(u.config)
	if err != nil {
		return []string{}, err
	}
	for _, f := range files {
		ip, err := ReadPath(u.config, f)
		if err != nil {
			return []string{}, err
		}
		np, err := NewPlan(u.config, ip.Name)
		if err != nil {
			return []string{}, err
		}
		if isUpgradable(ip, np) {
			u.upgrades = append(u.upgrades, np)
		}
	}
	return u.upgrades.Slice(), nil
}

// Upgrades returns a slice of plan names that are candidates for upgrading
func (u Upgrader) Upgrades() []string {
	return u.upgrades.Slice()
}

// Upgrade finally downloads and installs plan upgrades it returns a
// slice of errors if any occur.
//
// FIXME: we should make this transnational.
func (u Upgrader) Upgrade() []error {
	batch := NewBatch(u.config)
	return batch.ForEach(batch.DownloadInstall, u.upgrades)
}
