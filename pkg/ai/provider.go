package ai

type AiPlatform int

const (
	Cerebras AiPlatform = iota
	Groq
)

var platformToUrl = map[AiPlatform]string{
	Cerebras: "https://api.cerebras.ai",
	Groq:     "https://api.groq.com/openai",
}

type RatePolicy struct {
	RequestsMinute int
	RequestsDay    int

	TokensMinute int
	TokensDay    int
}

type Provider interface {
	Token() string
	Model() string
	Limits() RatePolicy
	BaseUrl() string
}

type CerebrasLllamaProvider struct {
	token string
}

func NewCerebrasLllamaProvider(token string) *CerebrasLllamaProvider {
	return &CerebrasLllamaProvider{
		token: token,
	}
}

func (clm *CerebrasLllamaProvider) Token() string {
	return clm.token
}

func (clm *CerebrasLllamaProvider) Model() string {
	return "llama3.1-8b"
}

func (clm *CerebrasLllamaProvider) Limits() RatePolicy {
	return RatePolicy{
		RequestsMinute: 30,
		RequestsDay:    14_400,
		TokensMinute:   30_000,
		TokensDay:      1_000_000,
	}
}

func (clm *CerebrasLllamaProvider) BaseUrl() string {
	return platformToUrl[Cerebras]
}

type GroqOss120Provider struct {
	token string
}

func NewGroqOss120Provider(token string) *GroqOss120Provider {
	return &GroqOss120Provider{
		token: token,
	}
}

func (gop *GroqOss120Provider) Token() string {
	return gop.token
}

func (gop *GroqOss120Provider) Model() string {
	return "openai/gpt-oss-120b"
}

func (gop *GroqOss120Provider) Limits() RatePolicy {
	return RatePolicy{
		RequestsMinute: 30,
		RequestsDay:    1_000,

		TokensMinute: 8_000,
		TokensDay:    200_000,
	}
}

func (gop *GroqOss120Provider) BaseUrl() string {
	return platformToUrl[Groq]
}
