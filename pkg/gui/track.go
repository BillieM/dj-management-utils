package gui

// func (e *guiEnv) updateTracksList(trackWidget *trackWidget, track database.SoundCloudTrack, playlistName string) {

// 	trackWidget.nameLabel.SetText(track.Name)

// 	trackWidget.url.Text = track.Permalink
// 	trackWidget.url.SetURLFromString(track.Permalink)

// 	trackWidget.genrePropertyLabel.SetText(track.Genre)
// 	trackWidget.tagListPropertyLabel.SetText(track.TagList)
// 	trackWidget.publisherPropertyLabel.SetText(track.PublisherArtist)
// 	trackWidget.soundcloudUserPropertyLabel.SetText(track.SoundCloudUser)

// 	if track.HasDownloadsLeft {
// 		trackWidget.downloadFileBtn.Text = "Download File"
// 		trackWidget.downloadFileBtn.Enable()
// 		trackWidget.downloadFileBtn.OnTapped = func() {
// 			trackWidget.downloadFileBtn.Disable()

// 			opEnv := e.opEnv()

// 			go func() {
// 				trackWidget.downloadingFile.Show()
// 				opEnv.DownloadSoundCloudFile(track, playlistName)
// 				trackWidget.downloadingFile.Hide()
// 				trackWidget.downloadFileBtn.Enable()
// 			}()
// 		}
// 	} else {
// 		trackWidget.downloadFileBtn.Text = "No Download"
// 		trackWidget.downloadFileBtn.Disable()
// 	}

// 	if track.PurchaseTitle != "" {
// 		trackWidget.purchaseBtn.Text = track.PurchaseTitle
// 		trackWidget.purchaseBtn.Enable()
// 		trackWidget.purchaseBtn.OnTapped = func() {
// 			opEnv := e.opEnv()

// 			go func() {
// 				trackWidget.purchaseBtn.Disable()
// 				opEnv.OpenSoundCloudPurchase(track)
// 				trackWidget.purchaseBtn.Enable()
// 			}()
// 		}
// 	} else {
// 		trackWidget.purchaseBtn.Text = "No Purchase Link"
// 		trackWidget.purchaseBtn.Disable()
// 	}

// 	trackWidget.Refresh()
// }
