package config

type ElevenlabsVoiceSettings struct {
	SimilarityBoost float32 `required:"false" default:"1.0"`
	Stability       float32 `required:"false" default:"1.0"`
	Style           float32 `required:"false" default:"0.0"`
	SpeakerBoost    bool    `required:"false" default:"true"`
	Speed           float32 `required:"false" default:"1.0"`
}

type ElevenLabsConfig struct {
	Enabled   bool                     `default:"false"`
	VoiceId   string                   `required:"false"`
	VoiceName string                   `required:"false"`
	Model     string                   `required:"true"`
	Token     string                   `required:"true"`
	Timeout   int                      `default:"30"`
	Settings  ElevenlabsVoiceSettings  `required:"false"`
}

type Config struct {
	Elevenlabs   ElevenLabsConfig
	Listen       string `default:":8080"`
	AudioBackend string `default:"pulse"`
	LogLevel     string `default:"info"`
	LogFormat    string `default:"text"`
}
