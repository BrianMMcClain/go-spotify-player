Go Spotify Player
===

A RESTful interface to simplify interacting with the Spotify Web API

Current features include:

- Play/Pause track
- Choice of playback device
- Track search

Example Config
---
```
{
  "port": 8000,
  "key": "SPOTIFY_OAUTH_KEY",
  "secret": "SPOTIFY_OAUTH_SECRET"
}
```

You can create your Spotify OAuth key and secret [here](https://developer.spotify.com/my-applications/#!/applications)

REST API
---
| Method | Endpoint | Action          | Example Data          | Example Response |
|--------|----------|-----------------|-----------------------|------------------|
| GET | /status     | Playback Status | N/A                   | {"playing":true,"deviceID":"x","deviceName":"x","url":"x","progress":0,"track":"x","artist":"x"} |
| GET    | /play    | Resume Playback | N/A                   | None             |
| POST   | /play    | Play Track      | {"uri": "TRACK_URI"}  | None             |
| GET    | /devices | Get List of Devices | N/A               | [{"name": "Device Name", "id": "xxxxxxxxx"}] |
| GET    | /pause   | Pause Playback  | N/A                   | None             |
| GET | /search/{keyword} | Search for Track | N/A            | {"name": "Track Name", "artist": "Track Artist", "URI": "xxx"} |
