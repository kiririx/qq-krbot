package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/json_util"
	"io"
	"os"
	"path"
	"regexp"
	"time"
)

var (
	client_id      = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	client_secret  = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	hash_secret    = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
	base_hosts     = "https://app-api.pixiv.net"
	_REFRESH_TOKEN = "vE94AY1QvM8BGcJuA6o1lPFBpfwv8YrDeuaH1AQWZRQ"
	AccessToken    = ""
	proxyUrl       = "http://127.0.0.1:7890"
	ExpireTime     int64
)

type PixivClient struct {
	Headers map[string]string
}

func GetHeaders() map[string]string {
	localTime := time.Now().Format(time.RFC3339)
	px := make(map[string]string)
	px["Accept-Language"] = "en-us"
	px["X-Client-Time"] = localTime
	px["X-Client-Hash"] = genClientHash(localTime)
	px["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	if AccessToken == "" {
		Auth()
	}
	px["Authorization"] = "Bearer " + AccessToken
	return px
}

func DownloadImg(url string) (string, error) {
	ext, _ := GetFileExt(url)
	fileName := algo_util.MD5(url) + "." + ext
	_, err := os.Stat("./photo/" + fileName)
	if err == nil {
		return fileName, nil
	}
	referer := "https://app-api.pixiv.net/"
	resp, err := http_util.Client(time.Second * 20).Proxy(proxyUrl).Headers(map[string]string{
		"Referer": referer,
	}).Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	f, err := os.Create("./photo/" + fileName)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	io.Copy(f, resp.Body)
	f.Close()
	return fileName, nil
}

// GetFileExt 获取文件的扩展名
func GetFileExt(fileAddr string) (string, error) {
	fileName := path.Base(fileAddr)
	reg, err := regexp.Compile("\\.(" + "jpg|png|jpeg|JPG|JPEG|PNG" + ")")
	if err != nil {
		return "", err
	}
	matchedExtArr := reg.FindAllString(fileName, -1)
	if matchedExtArr != nil && len(matchedExtArr) > 0 {
		ext := matchedExtArr[len(matchedExtArr)-1]
		return ext[1:], nil
	}
	return "", errors.New("获取文件扩展名失败")
}
func Auth() {
	if AccessToken != "" && time.Now().Unix() < ExpireTime {
		return
	}
	localTime := time.Now().Format(time.RFC3339)
	px := PixivClient{}
	px.Headers = make(map[string]string)
	px.Headers["Accept-Language"] = "en-us"
	px.Headers["X-Client-Time"] = localTime
	px.Headers["X-Client-Hash"] = genClientHash(localTime)
	px.Headers["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	json, err := http_util.Client(time.Second*10).
		Proxy("http://127.0.0.1:7890").
		PostFormGetJSON("https://oauth.secure.pixiv.net"+"/auth/token", map[string]string{
			"client_id":      client_id,
			"client_secret":  client_secret,
			"grant_type":     "refresh_token",
			"get_secure_url": "1",
			"refresh_token":  _REFRESH_TOKEN,
		})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	AccessToken = json["access_token"].(string)
	ExpireTime = time.Now().Unix() + int64(json["expires_in"].(float64))
}

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, hash_secret)
	return hex.EncodeToString(h.Sum(nil))
}

func Rank() []string {
	host := "https://app-api.pixiv.net"
	url := host + "/v1/illust/ranking"
	mode := "day_male_r18"
	filter := "for_ios"
	json, err := http_util.Client(time.Second*3).Proxy(proxyUrl).Headers(GetHeaders()).GetJSON(url, map[string]string{
		"mode":   mode,
		"filter": filter,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	photos := make([]string, 0)
	m, err := json_util.JSON2Map(json)
	illusts := m["illusts"].([]interface{})
	for _, illust := range illusts {
		image := illust.(map[string]interface{})["image_urls"].(map[string]interface{})
		if image["large"] != "" {
			photos = append(photos, image["large"].(string))
		} else if image["medium"] != "" {
			photos = append(photos, image["medium"].(string))
		} else if image["square_medium"] != "" {
			photos = append(photos, image["square_medium"].(string))
		}
	}
	return photos
}

func Recommend() (map[string]interface{}, error) {
	req := base_hosts + "/v1/illust/recommended"
	headers := GetHeaders()
	headers[`include_ranking_label`] = "true"
	json, err := http_util.Client(time.Second*3).Proxy(proxyUrl).Headers(headers).GetJSON(req, map[string]string{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	j, err := json_util.JSON2Map(json)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return j, nil
}

func Search(tag string) (map[string]interface{}, error) {
	url := base_hosts + "/v1/search/illust"
	headers := GetHeaders()
	json, err := http_util.Client(time.Second*10).Proxy(proxyUrl).Headers(headers).GetJSON(url, map[string]string{
		"word":          tag,
		"search_target": "partial_match_for_tags",
		"sort":          "popular_desc",
		"filter":        "for_ios",
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	j, err := json_util.JSON2Map(json)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return j, nil
}
