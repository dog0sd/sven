---
name: sven-tts
description: Simple Voice Enhanced Narrator - CLI and HTTP server for Text-to-Speech via ElevenLabs
license: GPLv3
---

# Simple Voice Enhanced Narrator

## Usage

No config file required. Set `ELEVENLABS_API_KEY` env and use CLI flags.

### CLI
```shell
sven -voice Rachel "Hello, world!"
sven -voice Adam -model eleven_v3 "[sighs] Finally... [laughs] just kidding!"
```

| Flag | Range | Description |
|------|-------|-------------|
| `-voice` | | Voice name (e.g. Rachel) |
| `-model` | | Model ID (default: `eleven_turbo_v2_5`) |
| `-backend` | `pulse`, `oto` | Audio backend |
| `-stability` | 0.0 - 1.0 | Voice stability |
| `-similarity` | 0.0 - 1.0 | Voice similarity boost |
| `-style` | 0.0 - 1.0 | Voice style |
| `-speed` | 0.7 - 1.2 | Voice speed |

### HTTP

```shell
sven -voice Rachel
```

**POST /tts**
```json
{
  "text": "Text to speak",
  "voice_settings": {
    "similarity_boost": 0.8,
    "stability": 0.7,
    "style": 0.0,
    "speed": 1.0
  },
  "ptext": "Previous context (optional)"
}
```

Voice settings use pointers (`*float32`) to distinguish "not set" from "0.0".

## Audio Tags (eleven_v3 only)

The `eleven_v3` model supports inline audio tags using square brackets `[]` directly in text. Tags can be mixed with regular text.

Example: `[sighs] I can't believe it... [laughs] just kidding!`

### Emotions / Delivery
| Tag | Description |
|-----|-------------|
| `[happy]` | Happy tone |
| `[sad]` | Sad tone |
| `[angry]` | Angry tone |
| `[annoyed]` | Annoyed tone |
| `[appalled]` | Appalled tone |
| `[excited]` | Excited tone |
| `[curious]` | Curious tone |
| `[surprised]` | Surprised tone |
| `[thoughtful]` | Thoughtful tone |
| `[sarcastic]` | Sarcastic tone |
| `[mischievously]` | Mischievous delivery |
| `[crying]` | Crying voice |
| `[whispers]` | Whispering voice |
| `[muttering]` | Muttering under breath |

### Laughter
| Tag | Description |
|-----|-------------|
| `[laughs]` | Standard laugh |
| `[laughs harder]` | Intense laugh |
| `[starts laughing]` | Begins laughing |
| `[chuckles]` | Soft chuckle |
| `[wheezing]` | Wheezing laugh |
| `[snorts]` | Snort laugh |

### Breathing / Non-Verbal
| Tag | Description |
|-----|-------------|
| `[sighs]` | Sigh |
| `[exhales]` | Exhale |
| `[exhales sharply]` | Sharp exhale |
| `[inhales deeply]` | Deep inhale |
| `[clears throat]` | Throat clearing |
| `[swallows]` | Swallowing sound |
| `[gulps]` | Gulping sound |

### Pauses
| Tag | Description |
|-----|-------------|
| `[short pause]` | Brief pause |
| `[long pause]` | Extended pause |

### Sound Effects
| Tag | Description |
|-----|-------------|
| `[gunshot]` | Gunshot sound |
| `[explosion]` | Explosion sound |
| `[applause]` | Applause |
| `[clapping]` | Clapping |

### Other
| Tag | Description |
|-----|-------------|
| `[sings]` | Singing voice |
| `[woo]` | Exclamation |
| `[fart]` | Fart sound |

> Note: tag effectiveness depends on the chosen voice and its training data. Not all tags work equally well with every voice.

