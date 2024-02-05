package operations

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/collection"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
SeperateSingleStem separates stems from a single file
*/
func (e *OpEnv) SeparateSingleStem(ctx context.Context, opts SeparateSingleStemOpts) {

	_, err := opts.Check()

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error checking opts",
				"There was an error whilst checking the options for the operation",
			),
		))
		return
	}

	stemEnv := e.StemEnvBuilder()

	e.Logger.Info("Checking file to separate")
	stemTrackArray, alreadyExistsCnt, errs := stemEnv.GetStemTracks([]string{opts.InFilePath}, opts.OutDirPath, opts.Type)

	if len(errs) > 0 {
		e.FinishError(fault.Wrap(
			errs[0],
			fmsg.WithDesc(
				"error building stem track array",
				"There was an error processing the given file",
			),
		))
		return
	}

	if alreadyExistsCnt > 0 {
		e.Logger.Info("Output file(s) already exist(s)")
		e.FinishSuccess(nil)
		return
	}

	e.Logger.Info("Converting file to stems")
	stemEnv.ConvertStemTracks(ctx, stemTrackArray)
	e.Logger.Info("Finished")

	e.FinishSuccess(nil)
}

/*
SeparateFolderStem separates stems from all files in a folder
*/
func (e *OpEnv) SeparateFolderStem(ctx context.Context, opts SeparateFolderStemOpts) {

	_, err := opts.Check()

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error checking opts",
				"There was an error whilst checking the options for the operation",
			),
		))
		return
	}

	stemEnv := e.StemEnvBuilder()

	e.Logger.Info("Finding files to convert")
	stemFilePaths, err := stemEnv.GetStemPaths(opts.InDirPath, opts.Recursion)

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error getting stem paths",
				"There was an error getting the paths of the files to convert",
			),
		))
		return
	}

	e.Logger.Infof("Found %v potential files to convert", len(stemFilePaths))

	e.Logger.Info("Checking found files")
	stemTrackArray, alreadyExistsCnt, errs := stemEnv.GetStemTracks(stemFilePaths, opts.OutDirPath, opts.Type)
	e.Logger.Infof("%v files already exist, %v left to convert", alreadyExistsCnt, len(stemTrackArray))

	for _, err := range errs {
		e.Logger.NonFatalError(fault.Wrap(
			err,
			fmsg.With("error building stem track array"),
		))
	}

	if len(stemTrackArray) == 0 {
		e.Logger.Info("No files to convert")
		e.FinishSuccess(nil)
		return
	}

	e.Logger.Info("Converting files to stems")
	stemEnv.ConvertStemTracks(ctx, stemTrackArray)
	e.Logger.Info("Finished")

	e.FinishSuccess(nil)
}

/*
ConvertSingleMp3 converts a single file to mp3
*/
func (e *OpEnv) ConvertSingleMp3(ctx context.Context, opts ConvertSingleMp3Opts) {
	_, err := opts.Check()

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error checking opts",
				"There was an error whilst checking the options for the operation",
			),
		))
		return
	}

	mp3Env := e.Mp3EnvBuilder()

	e.Logger.Info("Checking file to convert")
	convertTrackArray, alreadyExistsCnt, errs := mp3Env.GetMp3Tracks([]string{opts.InFilePath}, opts.OutDirPath)

	if len(errs) > 0 {
		e.FinishError(fault.Wrap(
			errs[0],
			fmsg.WithDesc(
				"error building convert track array",
				"There was an error processing the given file",
			),
		))
		return
	}

	if alreadyExistsCnt > 0 {
		e.Logger.Info("Output file(s) already exist(s)")
		e.FinishSuccess(nil)
		return
	}

	e.Logger.Info("Converting file to mp3")
	mp3Env.ConvertMp3Tracks(ctx, convertTrackArray)
	e.Logger.Info("Finished")

	e.FinishSuccess(nil)
}

/*
ConvertFolderMp3 converts all files in a folder to mp3
*/
func (e *OpEnv) ConvertFolderMp3(ctx context.Context, opts ConvertFolderMp3Opts) {

	_, err := opts.Check()

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error checking opts",
				"There was an error whilst checking the options for the operation",
			),
		))
		return
	}

	mp3Env := e.Mp3EnvBuilder()

	e.Logger.Info("Finding files to convert")
	convertFilePaths, err := mp3Env.GetMp3Paths(opts.InDirPath, opts.Recursion)
	e.Logger.Infof("Found %v potential files to convert", len(convertFilePaths))

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error getting convert paths",
				"There was an error getting the paths of the files to convert",
			),
		))
		return
	}

	e.Logger.Info("Checking found files")
	convertTrackArray, alreadyExistsCnt, errs := mp3Env.GetMp3Tracks(convertFilePaths, opts.OutDirPath)
	e.Logger.Infof("%v files already exist, %v left to convert", alreadyExistsCnt, len(convertTrackArray))

	for _, err := range errs {
		e.Logger.NonFatalError(fault.Wrap(
			err,
			fmsg.With("error building convert track array"),
		))
	}

	if len(convertTrackArray) == 0 {
		e.Logger.Info("No files to convert")
		e.FinishSuccess(nil)
		return
	}

	e.Logger.Info("Converting files to mp3")
	mp3Env.ConvertMp3Tracks(ctx, convertTrackArray)
	e.Logger.Info("Finished")

	e.FinishSuccess(nil)
}

/*
ReadCollection reads a collection for a given platform and stores it in the database
*/
func (e *OpEnv) ReadCollection(ctx context.Context, opts collection.ReadCollectionOpts) {

	collection := opts.Build(e.Config)

	err := collection.ReadCollection()

	if err != nil {
		e.FinishError(fault.Wrap(
			err,
			fmsg.WithDesc(
				"error reading collection",
				"There was an error reading the collection",
			),
		))
		return
	}
}

/*
GetPlaylist gets a playlist for a given platform and stores it in the database
*/

func (e *OpEnv) GetSoundCloudPlaylist(ctx context.Context, opts GetSoundCloudPlaylistOpts, p func(streaming.SoundCloudPlaylist, error)) {

	if !opts.Refresh {
		// check if playlist with same url already exists in database
		numPlaylists, err := e.SerenDB.GetNumSoundCloudPlaylistByURL(
			context.Background(),
			sql.NullString{Valid: true, String: opts.PlaylistURL},
		)

		if err != nil {
			p(
				streaming.SoundCloudPlaylist{},
				fault.Wrap(
					err,
					fmsg.WithDesc(
						"err checking if playlist exists in db by url",
						"Error checking if playlist with the same URL already exists in the applications database",
					),
				),
			)
			return
		}

		if numPlaylists > 0 {
			p(
				streaming.SoundCloudPlaylist{},
				fault.Wrap(
					helpers.ErrPlaylistAlreadyExists,
					fmsg.WithDesc(
						"playlist with same url already exists in db",
						"A playlist with that URL already exists",
					),
				),
			)
			return
		}
	}

	s := streaming.SoundCloud{
		ClientID: e.Config.SoundCloudClientID,
	}

	// get playlist from SoundCloud
	downloadedPlaylist, err := s.GetSoundCloudPlaylist(ctx, opts.PlaylistURL)

	if err != nil {
		p(
			streaming.SoundCloudPlaylist{},
			fault.Wrap(
				err,
				fmsg.WithDesc(
					"error getting playlist info from SoundCloud",
					"Error getting playlist information from SoundCloud",
				)),
		)
		return
	}

	if !opts.Refresh {
		// check if playlist with same external id already exists in database
		numPlaylists, err := e.SerenDB.GetNumSoundCloudPlaylistByExternalID(
			context.Background(),
			sql.NullInt64{Valid: true, Int64: downloadedPlaylist.ExternalID},
		)

		if err != nil {
			p(
				streaming.SoundCloudPlaylist{},
				fault.Wrap(
					err,
					fmsg.With("error checking if playlist already exists in database by external id"),
				),
			)
			return
		}

		if numPlaylists > 0 {
			p(
				streaming.SoundCloudPlaylist{},
				fault.Wrap(
					helpers.ErrPlaylistAlreadyExists,
					fmsg.With("playlist with same external id already exists in db"),
				),
			)
			return
		}
	}

	downloadedPlaylist.SearchUrl = opts.PlaylistURL

	dataP, dataT := downloadedPlaylist.ToDB()

	// save playlist to database
	err = e.SerenDB.TxUpsertSoundCloudPlaylistAndTracks(dataP, dataT)

	if err != nil {
		p(
			streaming.SoundCloudPlaylist{},
			fault.Wrap(err, fmsg.With("error saving playlist to database")),
		)
		return
	}

	p(downloadedPlaylist, nil)
}

/*
DownloadSoundCloudFile downloads a file straight from SoundCloud

# This only works for files with download enabled

playlistName is optional and is used to create a folder for the playlist within the download directory
*/
func (e *OpEnv) DownloadSoundCloudFile(track streaming.SoundCloudTrack, playlistName string) {

	if e.Config.SoundCloudClientID == "" {
		e.FinishError(
			fault.New("SoundCloud Client ID not set"),
		)
		return
	}

	downloadDir := e.Config.DownloadDir

	if playlistName != "" {
		downloadDir = helpers.JoinFilepathToSlash(downloadDir, playlistName)
	}

	s := streaming.SoundCloud{
		ClientID: e.Config.SoundCloudClientID,
	}

	filePath, err := s.DownloadFile(
		downloadDir,
		track.ExternalID,
	)

	if err != nil {
		e.FinishError(
			fault.Wrap(
				err,
				fmsg.With("error downloading file from SoundCloud"),
			),
		)
		return
	}

	track.LocalPath = filePath

	err = e.SerenDB.TxUpsertSoundCloudTracks([]data.SoundcloudTrack{track.ToDB()})
	if err != nil {
		e.FinishError(
			fault.Wrap(
				err,
				fmsg.With("error saving track to database"),
			),
		)
		return
	}

	e.FinishSuccess(
		map[string]any{
			"filepath": filePath,
		},
	)
}

/*
Flatten directory iteraves through a directory recursively and moves all files to the root of the directory
*/
func (e *OpEnv) FlattenDirectory(dirPath string) {
	// Get the list of files in the specified directory
	filePaths, err := helpers.GetFilesInDir(dirPath, true)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through each file and move it to the root of the directory
	// TODO: move move file to path to a helper func
	for _, filePath := range filePaths {

		fileName, err := helpers.GetFileNameFromFilePath(filePath)

		if err != nil {
			fmt.Println(err)
			return
		}

		fileExt, err := helpers.GetFileExtensionFromFilePath(filePath)

		if err != nil {
			fmt.Println(err)
			return
		}

		newPath := helpers.JoinFilepathToSlash(dirPath, fileName+fileExt)

		if filePath != newPath {
			err = os.Rename(filePath, newPath)

			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// TODO: remove directories in the specified directory...

}
