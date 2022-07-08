package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/str_util"
	"io"
	"os"
	"path"
	"qq-krbot/env"
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
	proxyUrl       = ""
	ExpireTime     int64
	PixivPageSize  = 30
)

func init() {
	proxyUrl = env.Conf["proxy.url"]
}

type PixivClient struct {
	Headers map[string]string
}

func GetHeaders() (map[string]string, error) {
	localTime := time.Now().Format(time.RFC3339)
	px := make(map[string]string)
	px["Accept-Language"] = "en-us"
	px["X-Client-Time"] = localTime
	px["X-Client-Hash"] = genClientHash(localTime)
	px["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	if AccessToken == "" {
		err := Auth()
		if err != nil {
			return nil, err
		}
	}
	px["Authorization"] = "Bearer " + AccessToken
	return px, nil
}

func DownloadImg(url string) (string, error) {
	ext, _ := GetFileExt(url)
	fileName := algo_util.MD5(url) + "." + ext
	_, err := os.Stat("./photo/" + fileName)
	if err == nil {
		return fileName, nil
	}
	referer := "https://www.pixiv.net/"
	resp, err := http_util.Client().Timeout(time.Second*99).Proxy(proxyUrl).Headers(map[string]string{
		"Referer": referer,
	}).Get(url, nil)
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

func DownloadImgAndTag(url string, tag string) (string, error) {
	dirPath := "./photo/" + tag
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	ext, _ := GetFileExt(url)
	fileName := algo_util.MD5(url) + "." + ext
	_, err = os.Stat(dirPath + "/" + fileName)
	if err == nil {
		return fileName, nil
	}
	referer := "https://www.pixiv.net/"
	resp, err := http_util.Client().Timeout(time.Second*99).Proxy(proxyUrl).Headers(map[string]string{
		"Referer": referer,
	}).Get(url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	f, err := os.Create(dirPath + "/" + fileName)
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
func Auth() error {
	if AccessToken != "" && time.Now().Unix() < ExpireTime {
		return nil
	}
	localTime := time.Now().Format(time.RFC3339)
	px := PixivClient{}
	px.Headers = make(map[string]string)
	px.Headers["Accept-Language"] = "en-us"
	px.Headers["X-Client-Time"] = localTime
	px.Headers["X-Client-Hash"] = genClientHash(localTime)
	px.Headers["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	json, err := http_util.Client().
		Timeout(time.Second*10).
		Proxy("http://127.0.0.1:7890").
		PostFormGetJSON("https://oauth.secure.pixiv.net"+"/auth/token", map[string]string{
			"client_id":      client_id,
			"client_secret":  client_secret,
			"grant_type":     "refresh_token",
			"get_secure_url": "1",
			"refresh_token":  _REFRESH_TOKEN,
		})
	if err != nil {
		return errors.New("Get pixiv token failed :" + err.Error())
	}
	AccessToken = json["access_token"].(string)
	ExpireTime = time.Now().Unix() + int64(json["expires_in"].(float64))
	return nil
}

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, hash_secret)
	return hex.EncodeToString(h.Sum(nil))
}

func Rank() ([]string, error) {
	host := "https://app-api.pixiv.net"
	url := host + "/v1/illust/ranking"
	mode := "day_male_r18"
	filter := "for_ios"
	headers, err := GetHeaders()
	if err != nil {
		return nil, err
	}
	json, err := http_util.Client().Proxy(proxyUrl).Headers(headers).GetJSON(url, map[string]string{
		"mode":   mode,
		"filter": filter,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	photos := make([]string, 0)
	illusts := json["illusts"].([]interface{})
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
	return photos, nil
}

func Recommend() (map[string]interface{}, error) {
	req := base_hosts + "/v1/illust/recommended"
	headers, err := GetHeaders()
	if err != nil {
		return nil, err
	}
	headers[`include_ranking_label`] = "true"
	json, err := http_util.Client().Proxy(proxyUrl).Headers(headers).GetJSON(req, map[string]string{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return json, nil
}

func Search(tag string, offset int) (map[string]interface{}, error) {
	url := base_hosts + "/v1/search/illust"
	headers, err := GetHeaders()
	if err != nil {
		return nil, err
	}
	json, err := http_util.Client().Timeout(time.Second*10).Proxy(proxyUrl).Headers(headers).GetJSON(url, map[string]string{
		"word":          tag,
		"search_target": "partial_match_for_tags",
		"sort":          "popular_desc",
		"filter":        "for_ios",
		"offset":        str_util.ToStr(offset),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return json, nil
}

func GetImgUrlForSearch(tag string, offset int) ([]string, error) {
	if offset < 1 {
		offset = PixivPageSize
	}
	m, err := Search(tag, offset)
	if err != nil {
		return nil, err
	}
	photos := make([]string, 0)
	illusts := m["illusts"].([]interface{})
re:
	for _, illust := range illusts {
		itype := illust.(map[string]interface{})["type"].(string)
		// 只采集插画资源
		if itype != "illust" {
			continue
		}
		tags := illust.(map[string]interface{})["tags"].([]interface{})
		for _, t := range tags {
			_t := t.(map[string]interface{})
			// 过滤掉这些tag
			if str_util.Contains(_t["name"].(string), "描き方", "参考", "講座", "資料", "漫画") {
				continue re
			}
		}
		image := illust.(map[string]interface{})["meta_single_page"].(map[string]interface{})
		var imgUrl = image["original_image_url"]
		if imgUrl != nil {
			photos = append(photos, imgUrl.(string))
		} else {
			images := illust.(map[string]interface{})["meta_pages"].([]interface{})
			// 过滤掉3张以上的作品，因为3张以上的很可能是漫画
			if len(images) > 3 {
				continue re
			}
			for _, v := range images {
				img := v.(map[string]interface{})
				url := img["image_urls"].(map[string]interface{})["original"]
				if url != nil {
					photos = append(photos, url.(string))
				}
			}
		}
	}
	return photos, nil
}
