package auth

import "github.com/ngocphuongnb/tetua/app/server"

var providers []server.AuthProvider

func New(authProviders ...server.AuthProvider) {
	for _, authProvider := range authProviders {
		addProvider(authProvider)
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
