---
name: sven-tts
description: Simple Voice Enhanced Narrator - command-line utility use Elevenlabs for Text-to-Speech and seemless playback on Pulse
license: GPLv3
---

# Simple Voice Enhanced Narrator

## Configuration

Config is loaded from (in order): `./sven.yml`, `~/.config/sven.yml`, `/etc/sven.yml`.


```yaml
audiobackend: pulse  # pulse (Linux) or oto (Windows/macOS)
elevenlabs:
  voiceid: <voice-id>
  model: eleven_turbo_v2_5
  token: <api-key>
  settings:
    similarityboost: 0.8
    stability: 0.7
    style: 0.0
    speed: 1.0
```

## Usage

### blocking cli-mode
```shell
sven "Hello, world!"
```
**All CLI flags:**
| Flag | Range | Description |
|------|-------|-------------|
| `-backend` | `pulse`, `oto` | Audio backend |
| `-stability` | 0.0 - 1.0 | Voice stability |
| `-similarity` | 0.0 - 1.0 | Voice similarity boost |
| `-style` | 0.0 - 1.0 | Voice style |
| `-speed` | 0.7 - 1.2 | Voice speed |

### HTTP 

```shell
sven
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

