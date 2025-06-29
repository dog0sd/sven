package config


type elevenlabsVoiceSettings struct {
	SimilarityBoost float32 `required:"false" default:"1.0"`
	Stability       float32 `required:"false" default:"1.0"`
	Style           float32 `required:"false" default:"0.0"`
	SpeakerBoost    bool    `required:"false" default:"true"`
}

type ElevenLabsConfig struct {
	Enabled  bool                    `default:"false"`
	VoiceId  string                  `required:"true"`
	Model    string                  `default:"eleven_multilingual_v2"`
	Token    string                  `required:"true"`
	Settings elevenlabsVoiceSettings `required:"false"`
}

type Config struct {
	Elevenlabs ElevenLabsConfig
	Port       string `default:":8080"`
}