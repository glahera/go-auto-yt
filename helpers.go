package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetChannels() []Channel {
	log.Info("Getting all channels from channels.json")
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Error("There was an error reading channels.json: ", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		log.Error("There was an error unmarshalling json: ", err)
	}
	log.Info("Successfully read all channels")
	return db
}

func CheckAll() {
	log.Info("Checking for all channels")
	allChannelsInDb := GetChannels()

	for _, item := range allChannelsInDb {
		channelName, err := GetChannelName(item.ChannelURL)
		if err != nil {
			log.Error("There was an error getting channel name: ", err)
		}
		channelType, err := GetChannelType(item.ChannelURL)
		if err != nil {
			log.Error("There was an error getting channel type: ", err)
		}

		if strings.Contains(item.ChannelURL, channelName) {
			videoId, _ := GetLatestVideo(channelName, channelType)

			if item.LatestDownloaded == videoId {
				log.Info("No new videos found for: ", item.ChannelURL)
			} else {
				log.Info("New video detected for: ", item.ChannelURL)
				go Download(channelName, channelType, "Audio Only")
				UpdateLatestDownloaded(item.ChannelURL, videoId)
			}
		}
	}
}

func CheckNow(channelName string, channelType string) Response {
	log.Info("Checking for new videos")
	allChannelsInDb := GetChannels()

	videoId, _ := GetLatestVideo(channelName, channelType)

	for _, item := range allChannelsInDb {
		if strings.Contains(item.ChannelURL, channelName) {
			if item.LatestDownloaded == videoId {
				log.Info("No new videos found for: ", channelName)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected"}
			} else {
				log.Info("New video detected for: ", channelName)
				Download(channelName, channelType, "Audio Only")
				UpdateLatestDownloaded(channelName, videoId)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected"}
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}
}

func GetChannelName(channelURL string) (string, error) {
	if channelURL != "" {
		return strings.Split(channelURL, "/")[4], nil
	}

	return "", fmt.Errorf("channelURL string is either empty or cant be parsed properly")
}

func GetChannelType(channelURL string) (string, error) {
	if channelURL != "" {
		return strings.Split(channelURL, "/")[3], nil
	}

	return "", fmt.Errorf("channelURL string is either empty or cant be parsed properly")
}

func CreateDirIfNotExist(dirName string) {
	log.Info("Creating channel directory")
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			log.Error("Couldn't create channel directory: ", err)
		} else {
			log.Info("Channel directory created successfully")
		}
	}
}

func RemoveAtIndex(s []Channel, index int) []Channel {
	return append(s[:index], s[index+1:]...)
}
