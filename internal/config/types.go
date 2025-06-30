package config

type elevenlabsVoiceSettings struct {
	SimilarityBoost float32 `required:"false" default:"1.0"`
	Stability       float32 `required:"false" default:"1.0"`
	Style           float32 `required:"false" default:"0.0"`
	SpeakerBoost    bool    `required:"false" default:"true"`
	Speed           float32 `required:"false" default:"1.0"`
}

type ElevenLabsConfig struct {
	Enabled  bool                    `default:"false"`
	VoiceId  string                  `required:"true"`
	Model    string                  `required:"true"`
	Token    string                  `required:"true"`
	Settings elevenlabsVoiceSettings `required:"false"`
}

type Config struct {
	Elevenlabs ElevenLabsConfig
	Port       string `default:":8080"`
}
