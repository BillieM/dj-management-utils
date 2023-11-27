package operations

import (
	"context"
	"fmt"
	"os"

	"github.com/billiem/seren-management/pkg/collection"
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
This file serves as an entrypoint for all operations
*/

/*
StepHandler is used to provide callbacks to the operations package
*/
type StepHandler interface {
	StepCallback(StepInfo)
	ExitCallback()
}

/*
StepInfo is returned to the StepCallback after each step

It provides information about the step that just finished
*/
type StepInfo struct {
	SkipLog    bool
	Progress   float64
	Message    string
	Error      error
	Importance helpers.Importance
}

/*
SeperateSingleStem separates stems from a single file
*/
func (e *OpEnv) SeparateSingleStem(ctx context.Context, opts SeparateSingleStemOpts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking file to separate"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray([]string{opts.InFilePath}, opts.OutDirPath, opts.Type)

	if len(errs) > 0 {
		e.step(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		e.step(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	e.step(stageStepInfo("Converting file to stems"))
	e.parallelProcessStemTrackArray(ctx, stemTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
SeparateFolderStem separates stems from all files in a folder
*/
func (e *OpEnv) SeparateFolderStem(ctx context.Context, opts SeparateFolderStemOpts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Finding files to convert"))
	stemFilePaths, err := e.getStemPaths(opts.InDirPath, opts.Recursion)
	e.step(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(stemFilePaths))))

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking found files"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray(stemFilePaths, opts.OutDirPath, opts.Type)
	e.step(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(stemTrackArray))))

	for _, err := range errs {
		e.step(warningStepInfo(err))
	}

	e.step(stageStepInfo("Converting files to stems"))
	e.parallelProcessStemTrackArray(ctx, stemTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ConvertSingleMp3 converts a single file to mp3
*/
func (e *OpEnv) ConvertSingleMp3(ctx context.Context, opts ConvertSingleMp3Opts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking file to convert"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray([]string{opts.InFilePath}, opts.OutDirPath)

	if len(errs) > 0 {
		e.step(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		e.step(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	e.step(stageStepInfo("Converting file to mp3"))
	e.parallelProcessConvertTrackArray(ctx, convertTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ConvertFolderMp3 converts all files in a folder to mp3
*/
func (e *OpEnv) ConvertFolderMp3(ctx context.Context, opts ConvertFolderMp3Opts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Finding files to convert"))
	convertFilePaths, err := e.getConvertPaths(opts.InDirPath, opts.Recursion)
	e.step(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(convertFilePaths))))

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking found files"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray(convertFilePaths, opts.OutDirPath)
	e.step(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(convertTrackArray))))

	for _, err := range errs {
		e.step(warningStepInfo(err))
	}

	e.step(stageStepInfo("Converting files to mp3"))
	e.parallelProcessConvertTrackArray(ctx, convertTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ReadCollection reads a collection for a given platform and stores it in the database
*/
func (e *OpEnv) ReadCollection(ctx context.Context, opts collection.ReadCollectionOpts) {

	collection := opts.Build(e.Config)

	err := collection.ReadCollection()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}
}

/*
GetPlaylist gets a playlist for a given platform and stores it in the database
*/

func (e *OpEnv) GetSoundCloudPlaylist(ctx context.Context, opts GetSoundCloudPlaylistOpts, p func(database.SoundCloudPlaylist, error)) {

	// check if playlist with same url already exists in database
	playlistByUrlCheck, err := e.SerenDB.GetSoundCloudPlaylistByURL(opts.PlaylistURL)

	if err != nil {
		e.step(dangerStepInfo(err))
		p(database.SoundCloudPlaylist{}, err)
		return
	}

	if playlistByUrlCheck.ID != 0 {
		p(playlistByUrlCheck, helpers.ErrPlaylistAlreadyExists)
		return
	}

	s := streaming.SoundCloud{
		ClientID: e.Config.SoundCloudClientID,
	}

	// get playlist from SoundCloud
	downloadedPlaylist, err := s.GetSoundCloudPlaylist(ctx, opts.PlaylistURL)

	if err != nil {
		e.step(dangerStepInfo(err))
		p(database.SoundCloudPlaylist{}, err)
		return
	}

	// check if playlist with same external id already exists in database
	playlistByExternalIDCheck, err := e.SerenDB.GetSoundCloudPlaylistByExternalID(downloadedPlaylist.ExternalID)

	if err != nil {
		e.step(dangerStepInfo(err))
		p(database.SoundCloudPlaylist{}, err)
		return
	}

	if playlistByExternalIDCheck.ID != 0 {
		p(playlistByExternalIDCheck, helpers.ErrPlaylistAlreadyExists)
		return
	}

	downloadedPlaylist.SearchUrl = opts.PlaylistURL

	// save playlist to database
	err = e.SerenDB.CreateSoundCloudPlaylist(downloadedPlaylist)

	if err != nil {
		e.step(dangerStepInfo(err))
		p(database.SoundCloudPlaylist{}, err)
		return
	}

	p(downloadedPlaylist, nil)
}

/*
DownloadSoundCloudFile downloads a file straight from SoundCloud

# This only works for files with download enabled

playlistName is optional and is used to create a folder for the playlist within the download directory
*/
func (e *OpEnv) DownloadSoundCloudFile(track database.SoundCloudTrack, playlistName string) {

	downloadDir := e.Config.DownloadDir

	if playlistName != "" {
		downloadDir = helpers.JoinFilepathToSlash(downloadDir, playlistName)
	}

	fmt.Println("downloading id: ", track.ExternalID)

	s := streaming.SoundCloud{
		ClientID: e.Config.SoundCloudClientID,
	}

	err := s.DownloadFile(
		downloadDir,
		track.ExternalID,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	/*
		https://api-v2.soundcloud.com/tracks/1367921728/download
			?client_id=1ZRkRXa5klyxfeCePlMbkWl1xNzz1Bu3&app_version=1700828706&app_locale=en
	*/
}

func (e *OpEnv) OpenSoundCloudPurchase(track database.SoundCloudTrack) {
	fmt.Printf("opening purchase for id: %v, URL: %s\n", track.ExternalID, track.PurchaseURL)

	// this is based on OS
	// gonna have to explore this at some point

	helpers.CmdExec("cmd", "/C", "start", "chrome.exe", track.PurchaseURL)

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
