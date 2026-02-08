# SVEN — Simple Voice Enhanced Narrator

Say it out loud.

## Quick Start

```bash
export ELEVENLABS_API_KEY=your_key
sven -voice Rachel "Hello, world!"
```

No config file needed.

## Install

Download the latest binary from the [releases](/releases) page.

### ElevenLabs API Key

Required permissions:
```
text_to_speech
voices_read
models_read (optional)
```

## Usage

### CLI

```bash
# Speak text
sven -voice Rachel "Hello, world!"

# Choose a model
sven -voice Rachel -model eleven_v3 "Hello!"

# Tune voice settings
sven -voice Rachel -stability 0.5 -speed 1.1 "Hello!"

# List available voices
sven voices

# List available models
sven models
```

### Server Mode

```bash
sven -voice Rachel
```

Starts HTTP server on `:8080` (configurable via config file). Send POST to `/tts`:

```bash
curl -d '{"text":"Hello!"}' localhost:8080/tts
```

With voice settings:

```bash
curl -d '{"text":"Hello!","voice_settings":{"speed":1.2}}' localhost:8080/tts
```

With context (previous text):

```bash
curl -d '{"text":"How are you?","ptext":"Hi, my love!"}' localhost:8080/tts
```

### MCP Server

```bash
ELEVENLABS_API_KEY=... ELEVENLABS_VOICE_ID=... ./mcp-server
```

Exposes `voice_it` tool over stdio for AI assistants. See [SKILL.md](./SKILL.md).

### All Flags

| Flag | Range | Description |
|------|-------|-------------|
| `-voice` | | Voice name (e.g. Rachel) |
| `-model` | | Model ID (default: `eleven_turbo_v2_5`) |
| `-backend` | `pulse`, `oto` | Audio backend |
| `-stability` | 0.0 - 1.0 | Voice stability |
| `-similarity` | 0.0 - 1.0 | Voice similarity boost |
| `-style` | 0.0 - 1.0 | Voice style |
| `-speed` | 0.7 - 1.2 | Voice speed |

---

## Configuration (Optional)

For persistent settings, create `sven.yml` in one of:
1. Current directory
2. `~/.config/sven.yml`
3. `/etc/sven.yml`

```yaml
elevenlabs:
  voicename: Rachel
  # Or use voice ID directly:
  # voiceid: iP95p4xoKVk53GoZ742B
  model: eleven_v3
  token: your_key
  settings:
    stability: 0.5
    speed: 1.0
    similarityboost: 0.5
    speakerboost: true

# Server settings
# listen: ":8080"
# audiobackend: pulse
# loglevel: info
# logformat: text
```

CLI flags override config file values. `ELEVENLABS_API_KEY` env overrides `token` in config.

---

## Audio Tags (eleven_v3)

The `eleven_v3` model supports inline tags:

```bash
sven -voice Rachel -model eleven_v3 "[sighs] I can't believe it... [laughs] just kidding!"
```

Tags: emotions (`[happy]`, `[sad]`, `[angry]`, `[sarcastic]`, ...), laughter (`[laughs]`, `[chuckles]`, ...), breathing (`[sighs]`, `[exhales]`, ...), pauses (`[short pause]`, `[long pause]`), sound effects (`[applause]`, `[gunshot]`, ...).

See [SKILL.md](./SKILL.md) for the full list.

---

## Dependencies

- ElevenLabs account with API key
- PulseAudio/PipeWire (`pulse` backend) or ALSA (`oto` backend)
