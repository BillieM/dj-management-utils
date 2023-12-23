# seren-dj-utils
A collection of utilities for DJ library management.


## Dependencies

ffmpeg (https://ffmpeg.org/)
- Required for audio processing/ demucs dependency
- Available on macOS via homebrew (`brew install ffmpeg`)
- Available on Ubuntu via apt (`sudo apt-get install ffmpeg`)
- Available on Windows via chocolatey (`choco install ffmpeg`)
- Binaries also available on the ffmpeg website (https://ffmpeg.org/download.html)
fyne (https://developer.fyne.io/ / https://github.com/fyne-io/fyne)
- UI framework
- requires golang/ a compatible c compiler/ graphics drivers on some platforms (macOS should work out of the box)
- detailed installation instructions available at `https://developer.fyne.io/started/`
demucs (github.com/facebookresearch/demucs)
- Used for audio source separation
- detailed installation instructions available at github.com/facebookresearch/demucs
- available via pip (`python3 -m pip install -U demucs`)
goose (github.com/pressly/goose)
- For running database migrations
- available on macOS via homebrew (`brew install goose`)
- avilable via go install (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
godotenv (github.com/joho/godotenv)
- to simplify loading environment variables from .env files to run development commands
- available via go install (`go install github.com/joho/godotenv/cmd/godotenv@latest`)

## Creating new DB migrations & queries
`github.com/sqlc-dev/sqlc`
- Main database queries within application logic are handled by sqlc
- sqlc generates go from queries inside ./db/queries
- add queries and run `sqlc generate`

`github.com/pressly/goose`
- Database migrations are handled by goose.
- migration status
    - `godotenv goose status`
- new migration
    - `godotenv goose create query_name sql`
- apply migrations
    - `godotenv goose up`

## XML Schema generation

DJ library integration is reliant on using generated go code from XSD schemas. We use xgen for this.

The results of this are not perfect for our usecase (likely a consequence of dodgy schema autogeneration due to the high complexity of the collection XMLs, particularly in the case of Traktor). But with some manual tweaking they seem to do the job. These generated go schema files are stored inside of `./pkg/collection`

xgen examples:
- `xgen -i ${xsd_schema_infile} -o ${go_schema_outfile} -l Go -p collection`

