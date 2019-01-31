package via

type Upgrader struct {
	config   *Config
	upgrades Plans
}

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

func (u Upgrader) Upgrades() []string {
	return u.upgrades.Slice()
}

func (u Upgrader) Upgrade() []error {
	batch := NewBatch(u.config)
	return batch.ForEach(batch.DownloadInstall, u.upgrades)
}
