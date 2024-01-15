# seren-dj-utils
A collection of utilities for DJ library management.


## Dependencies

ffmpeg ( [ffmpeg.org](https://ffmpeg.org/) )
- Required for audio processing/ demucs dependency
- Available on macOS via homebrew (`brew install ffmpeg`)
- Available on Ubuntu via apt (`sudo apt-get install ffmpeg`)
- Available on Windows via chocolatey (`choco install ffmpeg`)
- Binaries also available on the ffmpeg website (https://ffmpeg.org/download.html)

fyne ( [developer.fyne.io](https://developer.fyne.io/) / [github.com/fyne-io/fyne](https://github.com/fyne-io/fyne) )
- UI framework
- requires golang/ a compatible c compiler/ graphics drivers on some platforms (macOS should work out of the box)
- detailed installation instructions available at `https://developer.fyne.io/started/`

demucs ( [github.com/facebookresearch/demucs](https://github.com/facebookresearch/demucs) )
- Used for audio source separation
- detailed installation instructions available at github.com/facebookresearch/demucs
- available via pip (`python3 -m pip install -U demucs`)

goose ( [github.com/pressly/goose](https://github.com/pressly/goose) )
- For running database migrations
- available on macOS via homebrew (`brew install goose`)
- avilable via go install (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

godotenv ( [github.com/joho/godotenv](https://github.com/joho/godotenv) )
- to simplify loading environment variables from .env files to run development commands
- available via go install (`go install github.com/joho/godotenv/cmd/godotenv@latest`)

## Creating new DB migrations & queries

goose is used in order to run the migrations inside `./db/migrations` on the database, we combine this with godotenv to simplify these commands. For example...

- apply migrations
    - `godotenv goose up`
- new migration
    - `godotenv goose create query_name sql`
- migration status
    - `godotenv goose status`

We also use sqlc to generate go code from queries inside `./db/queries`. This is done by running `sqlc generate`

## XML Schema generation

DJ library integration is reliant on using generated go code from XSD schemas. We use xgen for this.

The results of this are not perfect for our usecase (likely a consequence of dodgy schema autogeneration due to the high complexity of the collection XMLs, particularly in the case of Traktor). But with some manual tweaking they seem to do the job. These generated go schema files are stored inside of `./pkg/collection`. As a result of the requirement for manual tweaking, we generate these once as part of the development process, not as part of CI/CD, should you need to regenerate the schema for any reason, an example is as follows:

xgen examples:
- `xgen -i ${xsd_schema_infile} -o ${go_schema_outfile} -l Go -p collection`

