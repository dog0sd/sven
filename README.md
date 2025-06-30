# Simple Voice Enhanced Narrator

A Go-based text-to-speech HTTP server(and cli tool too) that supports cloud-based (ElevenLabs) engines.

## Get started

Download binary from [releases](/releases)  
Make a `sven.yml` file with the following content:
```yaml
elevenlabs:
  voiceid: iP95p4xoKVk53GoZ742B
  model: eleven_turbo_v2_5
  token: <PASTE YOUR ELEVENLABS API KEY>
```
**running in command line**  
Run `./sven "Hello World!"`

**running as a server**  
`./sven`

**simple request**  
```bash
curl -H "content-type: application/json" -d '{"text":"Hello World!"}' localhost:8080/tts
```

**voice speed, style**  
For adjusting voice settings, look into config [sven.yml](/sven.yml).

Overriding voice settings available only in server mode.

The same settings you can adjust during any HTTP request:  
```bash
curl -H "content-type: application/json" -d '{"text":"Hello World!","voice_settings":{"speed":1.2}}' localhost:8080/tts
```

## Features

- Text-to-speech conversion via HTTP API
- Text-to-speech conversion via command line
- Support for ElevenLabs cloud TTS
- Tweaking voice settings

## TODO
- `--help` flag
- Add [cloud] google TTS support
- Add batch processing support
- write README