package auth

import (
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
)

var providers []server.AuthProvider

type NewProviderFn func(cfg map[string]string) server.AuthProvider

func New(newProviderFns map[string]NewProviderFn) {
	if config.Auth == nil || len(config.Auth.Providers) == 0 {
		return
	}

	for providerName, newProviderFn := range newProviderFns {
		if providerName != "local" && !utils.SliceContains(config.Auth.EnabledProviders, providerName) {
			continue
		}

		provider := newProviderFn(config.Auth.Providers[providerName])
		addProvider(provider)
	}
}

func Providers() []server.AuthProvider {
	return providers
}

func addProvider(authProvider server.AuthProvider) {
	for _, provider := range providers {
		if provider.Name() == authProvider.Name() {
			return
		}
	}

	providers = append(providers, authProvider)
}

func GetProvider(name string) server.AuthProvider {
	for _, provider := range providers {
		if provider.Name() == name {
			return provider
		}
	}

	return nil
}
