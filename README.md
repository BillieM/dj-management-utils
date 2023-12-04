# seren-dj-utils
A collection of utilities for DJ library management.


## Dependencies

`ffmpeg`
- Required for audio processing/ demucs dependecy
`sudo apt-get update && sudo apt-get install gcc libgl1-mesa-dev libegl1-mesa-dev libgles2-mesa-dev libx11-dev xorg-dev`
- fyne development dependencies
`github.com/facebookresearch/demucs`
- Stem separation model

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

