package config

// Application metadata
const (
	// AppName is the application identifier used for directories
	AppName = "silentcast"

	// ConfigName is the name used for configuration files
	ConfigName = "spellbook"

	// AppDisplayName is the user-friendly name shown in UI
	AppDisplayName = "SilentCast"

	// AppDescription is the application description
	AppDescription = "Silent Hotkey Task Runner"

	// AppRepo is the GitHub repository (without organization)
	AppRepo = "silentcast"

	// AppOrg is the GitHub organization
	AppOrg = "SphereStacking"
)

// YAML top-level key names - customize these to change the configuration structure
var (
	KeyDaemon    = "daemon"   // デーモン設定
	KeyHotkeys   = "hotkeys"  // ホットキー設定
	KeyShortcuts = "spells"   // ショートカット定義（Config構造体ではSpells）
	KeyActions   = "grimoire" // アクション定義（Config構造体ではGrimoire）
	KeyLogger    = "logger"   // ログ設定
	KeyUpdater   = "updater"  // 更新設定
)
