package utils

import (
	"abir/models"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendSms(login, code string) error {
	msgId := strconv.Itoa(rand.Intn(10000000))
	data := map[string]interface{}{"messages": []map[string]interface{}{{
		"recipient":  login,
		"message-id": viper.GetString("sms_broker.msg_alias") + msgId,
		"sms": map[string]interface{}{
			"originator": viper.GetString("sms_broker.originator"),
			"content": map[string]string{
				"text": code,
			},
		},
	}}}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", viper.GetString("sms_broker.endpoint"), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+os.Getenv("SMS_BROKER_PASSWORD"))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("couldn't connect to broker api")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("wrong request body")
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logrus.Debugf("%s", jsonData)
	return nil
}

func ImageToReader(image image.Image, mimetype string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if mimetype != "png" {
		err := png.Encode(buf, image)
		if err != nil {
			return nil, err
		}
	}
	if mimetype != "jpeg" {
		err := jpeg.Encode(buf, image, nil)
		if err != nil {
			return nil, err
		}
	}
	imageBytes := buf.Bytes()
	imageReader := bytes.NewReader(imageBytes)
	//imageBase64 := base64.NewDecoder(base64.StdEncoding, imageReader)
	return imageReader, nil
}
func GenerateFileName() (string, error) {
	return uuid.GenerateUUID()
}

func GenerateFileURL(fileName, folder, size string) *string {
	endpoint := viper.GetString("storage.endpoint")
	path := fmt.Sprintf("https://%s/%s/%s/%s/%s", endpoint, viper.GetString("storage.bucket"), folder, size, fileName)
	return &path
}

func GetFileUrl(fileName []string) *string {
	if len(fileName) <= 1 {
		return nil
	} else {
		return GenerateFileURL(fileName[1], fileName[0], "original")
	}
}

func GetSmallFileUrl(fileName []string) *string {
	if len(fileName) <= 1 {
		return nil
	} else {
		return GenerateFileURL(fileName[1], fileName[0], "small")
	}
}

func EncryptMessage(message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher([]byte(viper.GetString("aes.cipher")))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	//if _, err = io.ReadFull(crand.Reader, iv); err != nil {
	//	return "", fmt.Errorf("could not encrypt: %v", err)
	//}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher([]byte(viper.GetString("aes.cipher")))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func HidePhone(phone string) string {
	phone = phone[:5] + "*****" + phone[10:]
	return phone
}

func StripString(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

//type LList struct {
//	Lat float64 `json:"lat"`
//	Lng float64 `json:"lng"`
//}
//type RList struct {
//	RegionId int `json:"region_id"`
//	Points []LList `json:"points"`
//}
//func InPolygon(jsn string, lat, lng float64) (int, error) {
//	start := time.Now()
//	var result []RList
//	err := json.Unmarshal([]byte(jsn), &result)
//	if err != nil {
//		return 0, err
//	}
//	var poly []*geo.Point
//	for _, list := range result {
//		poly = []*geo.Point{}
//		for _, point := range list.Points {
//			poly = append(poly, geo.NewPoint(point.Lat, point.Lng))
//		}
//		newPoly := geo.NewPolygon(poly)
//		contains := newPoly.Contains(geo.NewPoint(lat, lng))
//		if contains {
//			duration := time.Since(start)
//			logrus.Print(duration.Milliseconds())
//			return list.RegionId, nil
//		}
//	}
//
//	return 0, nil
//}
type DistrictIdResponse struct {
	RegionId   int `json:"region_id"`
	DistrictId int `json:"district_id"`
}
type GeoMe struct {
	Data []DistrictIdResponse `json:"data"`
}

type CalculateRouteResult struct {
	Inside  int `json:"inside"`
	Outside int `json:"outside"`
}
type CalculateRouteResponse struct {
	Response CalculateRouteResult `json:"response"`
}

func GetStats(driverId int, typ string, start int) (any, error) {
	url := fmt.Sprintf(viper.GetString("services.dash")+"/api/v1/mobile/driver/%s/statistics/?driver=%v&start=%v", typ, driverId, start)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 55d0058d0aadf54a2592dab8ba4543b23724ffbf")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("couldn't connect to dash api")
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("wrong request body")
	}
	defer resp.Body.Close()
	var res any
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return 0, errors.New("error while parsing json")
	}
	return res, nil
}
func GetStatOrders(driverId, langId int) (any, error) {
	url := fmt.Sprintf(viper.GetString("services.dash")+"/api/v1/mobile/driver/statistics/?driver=%v&lang=%v", driverId, langId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 55d0058d0aadf54a2592dab8ba4543b23724ffbf")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("couldn't connect to dash api")
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("wrong request body")
	}
	defer resp.Body.Close()
	var res any
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return 0, errors.New("error while parsing json")
	}
	return res, nil
}
func GetMyDistrictId(lat, lng float64) (int, error) {
	url := fmt.Sprintf(viper.GetString("services.geo")+"/me?lat=%v&lng=%v", lat, lng)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("couldn't connect to geo api")
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("wrong request body")
	}
	defer resp.Body.Close()
	var responseData GeoMe
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return 0, errors.New("error while parsing json")
	}
	if len(responseData.Data) == 0 {
		return 0, errors.New("district not found")
	}
	return responseData.Data[0].DistrictId, nil
}

func CalculateRoute(points [][2]float64) (int, int, error) {
	var newLatLng []models.LatLng
	for _, point := range points {
		newLatLng = append(newLatLng, models.LatLng{Lat: point[1], Lng: point[0]})
	}
	newPoints := models.PointsRequest{Points: newLatLng}
	jsonData, err := json.Marshal(newPoints)
	if err != nil {
		return 0, 0, err
	}
	req, err := http.NewRequest("POST", viper.GetString("services.geo")+"/calculate", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, errors.New("couldn't connect to geo api")
	}

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errors.New("wrong request body")
	}
	defer resp.Body.Close()
	var responseData CalculateRouteResponse
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return 0, 0, errors.New("error while parsing json")
	}
	return responseData.Response.Inside, responseData.Response.Outside, nil
}

type SocketResponse struct {
	Success bool `json:"success"`
}

func SearchTaxi(orderId int) error {
	url := fmt.Sprintf(viper.GetString("services.socket")+"/search-drivers/%v", orderId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("couldn't connect to socket api")
	}
	defer resp.Body.Close()
	var responseData SocketResponse
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return errors.New("error while parsing json")
	}
	if !responseData.Success {
		return errors.New("error while searching")
	}
	return nil
}
func SkipTaxi(orderId int) error {
	url := fmt.Sprintf(viper.GetString("services.socket")+"/search-drivers/%v/skip", orderId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("couldn't connect to socket api")
	}
	defer resp.Body.Close()
	var responseData SocketResponse
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return errors.New("error while parsing json")
	}
	if !responseData.Success {
		return errors.New("error while skipping")
	}
	return nil
}

func CancelTaxi(orderId int) error {
	url := fmt.Sprintf(viper.GetString("services.socket")+"/search-drivers/%v/cancel", orderId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("couldn't connect to socket api")
	}
	defer resp.Body.Close()
	var responseData SocketResponse
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return errors.New("error while parsing json")
	}
	if !responseData.Success {
		return errors.New("error while cancelling")
	}
	return nil
}

func AcceptTaxi(orderId int) error {
	url := fmt.Sprintf(viper.GetString("services.socket")+"/search-drivers/%v/accept", orderId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("couldn't connect to socket api")
	}
	defer resp.Body.Close()
	var responseData SocketResponse
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return errors.New("error while parsing json")
	}
	if !responseData.Success {
		return errors.New("error while cancelling")
	}
	return nil
}
