package main

import (
	"encoding/json"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getUserMetadata(channelName string) Video {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", "https://www.youtube.com/user/"+channelName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	metaData := &ChannelInformation{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Fatal(err)
	}

	return Video{VideoID: metaData.ID}
}

func getChannelMetadata(channelName string) Video {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", "https://www.youtube.com/channel/"+channelName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	metaData := &ChannelInformation{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Fatal(err)
	}

	return Video{VideoID: metaData.ID}
}

func (c ChannelBasicInfo) GetLatestVideo() Video {
	if c.Type == "user" {
		return getUserMetadata(c.Name)
	}
	return getChannelMetadata(c.Name)
}

func (v Video) DownloadYTDL(fileExtension, downloadQuality string) error {
	cmd := exec.Command("youtube-dl", "-f", downloadQuality, "-o", "downloads/ %(uploader)s/video/ %(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
	log.Info(cmd.String())

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
		return err
	}

	return nil
}

func (v Video) DownloadAudioYTDL(fileExtension, downloadQuality string) error {
	if downloadQuality == "best" {
		downloadQuality = "0"
	} else if downloadQuality == "medium" {
		downloadQuality = "5"
	} else if downloadQuality == "worst" {
		downloadQuality = "9"
	}
	fileExtension = strings.Replace(fileExtension, ".", "", 1)
	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", fileExtension, "--audio-quality", downloadQuality, "-o", "downloads/ %(uploader)s/audio/ %(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)

	log.Info(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
		return err
	}

	return nil
}
