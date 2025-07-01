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

---

## Features
- Text-to-speech conversion via HTTP API and CLI
- Support for ElevenLabs cloud TTS
- Customizable voice settings (speed, style, etc.)
- Context-aware synthesis with `ptext` (previous text)

---

## Getting Started

### Download
Download the latest binary from the [releases](/releases) page.

### Configuration
Create a `sven.yml` file in the same directory as the binary with the following content:
```yaml
elevenlabs:
  voiceid: iP95p4xoKVk53GoZ742B
  model: eleven_turbo_v2_5
  token: <PASTE YOUR ELEVENLABS API KEY>
```

### Running SVEN
**Command Line:**
```bash
./sven "Hello World!"
```

**Server Mode:**
```bash
./sven
```

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

## Dependencies
- Go 1.20+
- Elevenlabs account

---

## Contributing
Contributions, issues, and feature requests are welcome! Please open an issue or submit a pull request.

---

## TODO
- ~~`--help` flag~~
- Specify exact output format for API request
- Async audio playback
- ~~Validate value boundaries in request and config at startup~~
- Specify voice by name
- Structured log format for the server
- Pronunciation dictionary support in config
- ~~Optional previous_text for context~~

