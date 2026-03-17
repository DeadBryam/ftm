package shortener

type Provider interface {
	Name() string
	Shorten(longURL, custom string) (string, error)
}

type ShortenError struct {
	Reason  string
	Message string
}

func (e ShortenError) Error() string {
	return "shortener: " + e.Reason + " - " + e.Message
}

func IsDomainBlocked(err error) bool {
	if se, ok := err.(ShortenError); ok {
		return se.Reason == "DOMAIN_BLOCKED"
	}
	return false
}
