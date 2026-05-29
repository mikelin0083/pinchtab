package stealth

import (
	"strings"

	"github.com/pinchtab/pinchtab/internal/config"
)

type LaunchContract struct {
	Args  []string
	Flags map[string]bool
}

func BuildLaunchContract(cfg *config.RuntimeConfig, level Level) LaunchContract {
	persona := BrowserPersona{}
	customUA := ""
	if cfg != nil {
		persona = BuildPersona(cfg.UserAgent, cfg.ChromeVersion)
		customUA = strings.TrimSpace(cfg.UserAgent)
	}

	args := []string{
		"--disable-automation",
		"--enable-automation=false",
		"--disable-blink-features=AutomationControlled",
		"--enable-network-information-downlink-max",
	}
	// Only pin --user-agent when the operator configured an EXPLICIT custom UA.
	// Passing --user-agent makes Chrome return EMPTY high-entropy UA Client Hints
	// (architecture, platformVersion, uaFullVersion, fullVersionList) from
	// navigator.userAgentData.getHighEntropyValues — a fingerprint inconsistency a
	// real Chrome never exhibits. With no custom UA, Chrome's native UA + native
	// UA-CH are already self-consistent, so leave them intact.
	pinUA := customUA != "" && persona.UserAgent != ""
	if pinUA {
		args = append(args, "--user-agent="+persona.UserAgent)
	}
	if persona.Language != "" {
		args = append(args, "--lang="+persona.Language)
	}

	return LaunchContract{
		Args: args,
		Flags: map[string]bool{
			"automationControlledDisabled": true,
			"enableAutomationFalse":        true,
			"downlinkMaxFlag":              true,
			"globalUserAgent":              pinUA,
			"globalLanguage":               persona.Language != "",
		},
	}
}

func HasLaunchArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

func HasLaunchArgPrefix(args []string, prefix string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}
