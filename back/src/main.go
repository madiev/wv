package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"wv/ws"
)

type videoInfo struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	ChannelId    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`
	Duration     string `json:"duration"`
	PublishedAt  string `json:"published_at"`
}

type result struct {
	Video    []videoInfo       `json:"video"`
	Channel  map[string]string `json:"channel"`
	Playlist map[string]string `json:"playlist"`
}

var developerKey = os.Getenv("YT_KEY")

func findFileWithRegex(dir, pattern string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Ошибка компиляции регулярного выражения: %v\n", err)
		return ""
	}

	var foundPath string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ошибка обхода пути: %v\n", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if re.MatchString(info.Name()) {
			foundPath = path
			log.Println("Найден файл: " + path)
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return foundPath
}

func handlerDownloadTask(c *gin.Context, ws *ws.Client) {
	var req struct {
		ID string `form:"v" binding:"required"`
	}
	err := c.Bind(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	dir := os.Getenv("DIR_DOWNLOAD")
	find := fmt.Sprintf("https://www.youtube.com/watch?v=%s", req.ID)
	cmd := exec.Command("/bin/sh", "../download/download.sh", dir, find)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		log.Println(err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		err = ws.Conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("write message:", err)
			break
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v\n", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	// err = cmd.Run()
	// if err != nil {
	// 	log.Printf("Ошибка выполнения скрипта: %v\n", err)
	// 	log.Printf("Ошибки вывода: %s\n", stderr.String())
	// 	_ = c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	response := map[string]string{
		"status": "ok",
	}

	c.JSON(http.StatusOK, response)
}

func handlerWatchVideo(c *gin.Context) {
	var req struct {
		ID string `form:"v" binding:"required"`
	}
	err := c.Bind(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	directory := "../download"
	regexPattern := `^.+\[` + req.ID + `\]\..+$`

	filePath := findFileWithRegex(directory, regexPattern)
	if filePath == "" {
		log.Println("Файл не найден")
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Файл не найден: %v\n", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	size := fileInfo.Size()

	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Ошибка открытия файла: %v\n", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	err = os.Remove(filePath)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//c.Request.Header.Add("Ассept-Ranges", "bytes")
	//c.Request.Header.Add("Content-Type", "video/mp4")
	//c.Request.Header.Add("Content-Length", strconv.Itoa(len(data)))
	//c.Data(http.StatusOK, "video/mp4", data)

	rangeHeader := c.Request.Header.Get("Range")
	if len(rangeHeader) > 0 {
		var start, end string

		rangeBytes := rangeHeader[len("bytes="):]
		i := strings.Index(rangeBytes, "-")
		start = rangeBytes[:i]
		end = rangeBytes[i+1:]

		extraHeaders := map[string]string{
			"Accept-Ranges": "bytes",
			"Content-Range": fmt.Sprintf("bytes %s-%s/%d", start, end, size),
		}
		c.DataFromReader(
			http.StatusPartialContent,
			size,
			"video/mp4",
			f,
			extraHeaders,
		)
	} else {
		bs := make([]byte, size)
		_, err = bufio.NewReader(f).Read(bs)
		if err != nil && err != io.EOF {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "video/mp4", bs)
	}
}

func handlerSearch(c *gin.Context) {
	var req struct {
		//Query  string `form:"q" binding:"required"`
		Query string `form:"q"`
		Limit int64  `form:"limit" binding:"min=0,max=50"`
	}

	err := c.Bind(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	arr := []string{"id", "snippet"}

	// Make the API call to YouTube.
	call := service.Search.List(arr).
		Q(req.Query).
		MaxResults(req.Limit)
	response, err := call.Do()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	var resp result
	var video []videoInfo
	channel := make(map[string]string)
	playlist := make(map[string]string)

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			//fmt.Printf("%#v\n", item.Snippet)
			var info videoInfo

			call := service.Videos.List([]string{"contentDetails"}).Id(item.Id.VideoId)
			response, err := call.Do()
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			duration := response.Items[0].ContentDetails.Duration

			info.ID = item.Id.VideoId
			info.Title = item.Snippet.Title
			info.Description = item.Snippet.Description
			info.Thumbnail = item.Snippet.Thumbnails.High.Url
			info.ChannelId = item.Snippet.ChannelId
			info.ChannelTitle = item.Snippet.ChannelTitle
			info.Duration = duration
			info.PublishedAt = item.Snippet.PublishedAt

			video = append(video, info)
		case "youtube#channel":
			channel[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			//fmt.Printf("%#v\n", item.Snippet)
			playlist[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	resp.Video = video
	resp.Channel = channel
	resp.Playlist = playlist

	c.JSON(http.StatusOK, resp)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	client := ws.Client{}

	api := router.Group("/api")
	api.GET("/search", handlerSearch)
	api.GET("/task", func(c *gin.Context) {
		handlerDownloadTask(c, &client)
	})
	api.GET("/video", handlerWatchVideo)
	api.GET("/ws", client.HandlerWebsocket)

	router.Run(":9999")
}
