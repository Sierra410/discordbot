package main

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"strconv"

// 	"github.com/bwmarrin/discordgo"
// )

// func init() {
// 	mustAddExplicitCommand(&explicitCommand{
// 		adminOnly:   true,
// 		chatType:    chatTypeAny,
// 		command:     "shitshow",
// 		helpMessage: "Nothing is gonna help",
// 		function:    commandShitshow,
// 	})
// }

// var shitshowJobs = NewSmap()

// func commandShitshow(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
// 	var flag *Flag

// 	v, ok := shitshowJobs.Get(cmd.message.ChannelID)

// 	if ok {
// 		v.(*Flag).Set(false)

// 		return "", nil
// 	}

// 	flag = NewFlag(true)
// 	shitshowJobs.Set(cmd.message.ChannelID, flag)

// 	defer func() {
// 		shitshowJobs.Delete(cmd.message.ChannelID)
// 	}()

// outer:
// 	for flag.Get() {
// 		images, err := getRandomDbrImages("")
// 		if err != nil {
// 			return "", err
// 		}

// 		for _, x := range images {
// 			sendFileFromUrl(
// 				session,
// 				cmd.message.ChannelID,
// 				x.Reprs.Small,
// 				strconv.Itoa(x.Id)+"."+x.Format,
// 				"",
// 			)

// 			if !flag.Get() {
// 				break outer
// 			}
// 		}
// 	}

// 	return "Shitshow was aborted", nil
// }

// var (
// 	dbrDomain = "https://derpibooru.org"
// 	dbrSearch = dbrDomain + "/api/v1/json/search/images"
// 	dbrKey    = "QJP4wT83_ToWxxy9q2Wj"
// )

// type DbrImageRepresentations struct {
// 	Full       string `json:"full"`
// 	Large      string `json:"large"`
// 	Medium     string `json:"medium"`
// 	Small      string `json:"small"`
// 	Tall       string `json:"tall"`
// 	Thumb      string `json:"thumb"`
// 	ThumbSmall string `json:"thumb_small"`
// 	ThumbTiny  string `json:"thumb_tiny"`
// }

// type DbrImage struct {
// 	Id       int                     `json:"id"`
// 	Format   string                  `json:"format"`
// 	MimeType string                  `json:"mime_type"`
// 	Reprs    DbrImageRepresentations `json:"representations"`
// }

// func sendFileFromUrl(session *discordgo.Session, channelId string, url, name, contentTypeOverride string) (*discordgo.Message, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if contentTypeOverride == "" {
// 		contentTypeOverride = resp.Header.Get("Content-Type")
// 	}

// 	return session.ChannelMessageSendComplex(
// 		channelId,
// 		&discordgo.MessageSend{
// 			Files: []*discordgo.File{
// 				&discordgo.File{
// 					Name:        name,
// 					ContentType: contentTypeOverride,
// 					Reader:      resp.Body,
// 				},
// 			},
// 		},
// 	)
// }

// func getRandomDbrImages(q string) ([]DbrImage, error) {
// 	q = url.QueryEscape("explicit,upvotes.gte:500,-watersports" + ",-foalcon")
// 	u := dbrSearch + "?q=" + q + "&sf=random&sd=desc&filder_id=42053&per_page=50&key=" + dbrKey

// 	resp, err := http.Get(u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	bytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s := &struct {
// 		I []DbrImage `json:"images"`
// 	}{}

// 	err = json.Unmarshal(bytes, s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.I, nil
// }
