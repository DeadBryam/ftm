package shortener

type MultiProvider struct {
	providers []Provider
}

func NewMulti(providers ...Provider) *MultiProvider {
	return &MultiProvider{
		providers: providers,
	}
}

func DefaultMulti() *MultiProvider {
	return NewMulti(
		NewISGD(),
		NewCleanURI(),
		NewDagd(),
		NewTinyURL(),
	)
}

func (m *MultiProvider) Name() string {
	return "multi"
}

func (m *MultiProvider) Shorten(longURL, custom string) (string, error) {
	var lastErr error
	
	for _, p := range m.providers {
		shortURL, err := p.Shorten(longURL, custom)
		if err == nil {
			return shortURL, nil
		}
		
		lastErr = err
		
		if !IsDomainBlocked(err) {
			return "", err
		}
	}
	
	return "", lastErr
}
