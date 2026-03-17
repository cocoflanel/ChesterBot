package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// throws a basic GET request at the provided endpoint
func restGet(url string) string {
	resp, err := http.Get(url)
	errorCheck(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	errorCheck(err)

	return string(data)
}

// uploads a single image to imgchest, returns a link to the image.
func restPutImg(imgUrl string) (string, error) {

	if len(imgUrl) == 0 {
		return "", errors.New("Blank URL Provided to REST PUT")
	}

	// loads API key from a file called imgchest.key
	apiKey := loadFile("imgchest.key")
	if len(apiKey) == 0 {
		return "", errors.New("imgchest.key file empty or does not exist, please create file and populate with an ImgChest API key")
	}

	// HTTP endpoint for uploading new posts
	createPostEndpoint := "https://api.imgchest.com/v1/post"

	// pull image from URL
	resp, err := http.Get(imgUrl)
	if err != nil {
		return "", errors.Join(errors.New("Error retrieving the file, "), err)
	}

	defer resp.Body.Close()

	//error out if we don't get a good response code
	if resp.StatusCode != 200 {
		return "", errors.New("Recieved response code " + strconv.Itoa(resp.StatusCode) + " " + resp.Status + " when trying to retrieve image from url " + imgUrl)
	}
	//read the image to memory
	img, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Join(errors.New("Error reading the response body, "), err)
	}

	//title string based on date and current time
	titleString := "Chester Upload " + time.Now().Format(time.RFC3339Nano)

	/*Now, to upload to the ImgChest API.
	  Documentation for this is extremely vague about how this works, but, it requires a multipart form data request
	  using mime/multipart, we set up a form, write each field (per API documentation) and then write the image.
	  Then we create a POST request, set this form as the body, add headers, and send it.

	  The images[] field is the trickiest part and took me a lot of looking around to understand.
	  in this case, create a Form File, name it images[], and then pass it the names of the image being uploaded.
	  Have yet to test it with more than one image at a time but hey this works great for now.

	  TODO: figure out how to set the Anonymous field to exactly a boolean value as the API expects.
	  Writing it as is throws an error that the
	*/

	bodyBuffer := &bytes.Buffer{}
	mpw := multipart.NewWriter(bodyBuffer)
	mpw.WriteField("title", titleString)
	mpw.WriteField("privacy", "secret")
	//mpw.WriteField("anonymous", "false")
	imgWriter, err := mpw.CreateFormFile("images[]", "webp.webp")
	if err != nil {
		return "", errors.Join(errors.New("Error constructing Form Data: "), err)
	}
	imgWriter.Write(img)
	mpw.Close()

	// Creating a POST request on createPostEndpoint
	request, err := http.NewRequest("POST", createPostEndpoint, bodyBuffer)
	if err != nil {
		return "", err
	}

	// Setting Headers
	//Content type to tell the API this is a multipart form, and we accept a JSON response
	request.Header.Add("Content-Type", mpw.FormDataContentType())
	request.Header.Add("Accept", "application/json")
	//Adding API Key as a header
	request.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{}

	// Performing POST request
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Reading the response to the request
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseStruct postResponse
	err = json.Unmarshal(responseBody, &responseStruct)
	if err != nil || len(responseStruct.Data.Images) == 0 {
		return "", errors.Join(errors.New("error in API response: "+string(responseBody)+" - "), err)
	}
	//debug
	responseBodyString := string(responseBody)
	fmt.Println(responseBodyString)

	//return a link to the created image on imgChest
	return responseStruct.Data.Images[0].Link, nil
}
func deserializeTest() (string, error) {
	input := []byte(`{"data":{"id":"ne7bnmvzj75","title":"Chester Upload 2026-02-17T19:25:04.340400342-08:00","username":"cocoflanne728","privacy":"secret","report_status":1,"views":0,"nsfw":0,"image_count":1,"created":"2026-02-18T03:25:04.000000Z","delete_url":"https:\/\/imgchest.com\/p\/7pyveh6m9jk48\/delete","images":[{"id":"162e2f602561","description":null,"link":"https:\/\/cdn.imgchest.com\/files\/162e2f602561.webp","position":1,"nsfw":false,"created":"2026-02-18T03:25:04.000000Z","original_name":"webp.webp"}]}}`)
	var responseStruct postResponse
	err := json.Unmarshal(input, &responseStruct)
	if err != nil || len(responseStruct.Data.Images) == 0 {
		return "", errors.Join(errors.New("error in API response: "+string(input)+" - "), err)
	}
	//debug
	responseBodyString := string(input)
	fmt.Println(responseBodyString)

	//return a link to the created image on imgChest
	return responseStruct.Data.Images[0].Link, nil
}

type imgFileInfo struct {
	Id            string `json:"id"`
	Description   string `json:"description"`
	Link          string `json:"link"`
	Position      int    `json:"position"`
	Nsfw          bool   `json:"nsfw"`
	Created       string `json:"created"`
	Original_name string `json:"original_name"`
}
type postResponse struct {
	Data struct {
		Id            string        `json:"id"`
		Title         string        `json:"title"`
		Username      string        `json:"username"`
		Privacy       string        `json:"privacy"`
		Report_status int           `json:"report_status"`
		Views         int           `json:"views"`
		Nsfw          int           `json:"nsfw"`
		Image_count   int           `json:"image_count"`
		Created       string        `json:"created"`
		Delete_url    string        `json:"delete_url"`
		Images        []imgFileInfo `json:"images"`
	} `json:"data"`
}
