# Simple Voice Enhanced Narrator

A Go-based text-to-speech (TTS) HTTP server and CLI tool supporting cloud-based engines (currently ElevenLabs). 
SVEN enables you to convert text to speech via a simple API or command line, with customizable voice settings and autoplay it on your audio device.

---

## Table of Contents
- [Features](#features)
- [Getting Started](#getting-started)
  - [Download](#download)
  - [Configuration](#configuration)
  - [Running SVEN](#running-sven)
- [API Usage](#api-usage)
  - [Request Format](#request-format)
  - [Response Format](#response-format)
  - [Examples](#examples)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [TODO](#todo)

Take skill for Claude Code: [SKILL.md](./SKILL.md).

---

## Features
- Text-to-speech conversion via HTTP API and CLI
- MCP (stdio) and SKILL.md
- Eleven v3 tags
- Pluggable audio backends: PulseAudio/PipeWire (`pulse`) and ALSA (`oto`)
- Voice settings (stability, speed, style, etc.)
- Context-aware synthesis with `ptext` (previous text)

---

## Getting Started

### Download
Download the latest binary from the [releases](/releases) page.


### Elevenlabs API KEY

Permissions for API KEY:
```
text_to_speech
voices_read
models_read (optional)
```

### Configuration

Configuration is loaded and merged from multiple locations (in priority order):
1. `sven.yml` / `sven.yaml` in the current directory
2. `~/.config/sven.yml` / `~/.config/sven.yaml`
3. `/etc/sven.yml` / `/etc/sven.yaml`

You can also set `ELEVENLABS_API_KEY` environment variable to override the token from config.

Example `sven.yml`:
```yaml
# Listen address for HTTP server (default: ":8080")
# listen: ":8080"
# listen: "127.0.0.1:8080"

# Audio backend: "pulse" (PulseAudio/PipeWire) or "oto" (ALSA)
audiobackend: pulse

# Logging: level (debug, info, warn, error), format (text, json)
# loglevel: info
# logformat: text

elevenlabs:
  voiceid: iP95p4xoKVk53GoZ742B
  # Or use voice name instead of voiceid:
  # voicename: Rachel
  model: eleven_turbo_v2_5
  token: <PASTE YOUR ELEVENLABS API KEY>
  # API timeout in seconds (default: 30)
  # timeout: 30
```

### Running SVEN
**Command Line:**
```bash
./sven "Hello World!"
```

**With voice settings:**
```bash
./sven -stability 0.5 -similarity 0.8 -speed 1.1 "Hello World!"
```

**With a specific audio backend:**
```bash
./sven -backend oto "Hello World!"
```

**List available voices:**
```bash
./sven voices
```

**List available models:**
```bash
./sven models
```

**All CLI flags:**
| Flag | Range | Description |
|------|-------|-------------|
| `-backend` | `pulse`, `oto` | Audio backend |
| `-stability` | 0.0 - 1.0 | Voice stability |
| `-similarity` | 0.0 - 1.0 | Voice similarity boost |
| `-style` | 0.0 - 1.0 | Voice style |
| `-speed` | 0.7 - 1.2 | Voice speed |

**Server Mode** (starts on `:8080` by default, configure with `listen` in config):
```bash
./sven
```

**MCP Server Mode:**
```bash
./mcp-server
```
The MCP server exposes a `voice_it` tool over stdio, allowing AI assistants (e.g. Claude) to speak text aloud. Configure it via `ELEVENLABS_API_KEY` and `ELEVENLABS_VOICE_ID` environment variables.

---

## API Usage

### Request Format
Send a POST request to `/tts` with JSON body:
```json
{
  "text": "Hello World!",
  "voice_settings": {
    "speed": 1.2
  },
  "ptext": "Previous context text (optional)"
}
```
- `text` (string, required): The text to synthesize.
- `voice_settings` (object, optional): Voice parameters (e.g., speed, style). Look into [sven.yml](./sven.yml) for more voice settings.
- `ptext` (string, optional): Previous text for context.

### Response Format
The response is a text (OK or Internal Server Error). Didn't you think I will return you mp3 file?

### Examples
**Simple request:**
```bash
curl -H "content-type: application/json" -d '{"text":"Hello World!"}' localhost:8080/tts
```

**With voice settings:**
```bash
curl -H "content-type: application/json" -d '{"text":"Hello World!","voice_settings":{"speed":1.2}}' localhost:8080/tts
```

**With context:**
```bash
curl -H "content-type: application/json" -d '{"text":"How are you?","ptext":"Hi, my love!"}' localhost:8080/tts
```

---

## Audio Tags (eleven_v3)

The `eleven_v3` model supports inline audio tags using square brackets directly in text. These allow adding emotions, sound effects, pauses, and non-verbal sounds to speech.

```bash
./sven "[sighs] I can't believe it... [laughs] just kidding!"
```

Tags include emotions (`[happy]`, `[sad]`, `[angry]`, `[sarcastic]`, ...), laughter (`[laughs]`, `[chuckles]`, ...), breathing (`[sighs]`, `[exhales]`, ...), pauses (`[short pause]`, `[long pause]`), and sound effects (`[applause]`, `[gunshot]`, ...).

See [SKILL.md](./SKILL.md) for the full list of supported tags.

---

## Dependencies
- Go 1.22+
- ElevenLabs account
- PulseAudio/PipeWire (`pulse` backend) or ALSA (`oto` backend)

---

## Contributing
Contributions, issues, and feature requests are welcome! Please open an issue or submit a pull request.

---

## TODO
- ~~`--help` flag~~
- Specify exact output format for API request
- Async audio playback
- ~~Validate value boundaries in request and config at startup~~
- ~~Specify voice by name~~
- ~~Structured log format for the server~~
- Pronunciation dictionary support in config
- ~~Optional previous_text for context~~
- ~~MCP server mode~~
- ~~Pluggable audio backends~~

