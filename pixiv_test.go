package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"qq-krbot/handler"
	"regexp"
	"testing"
	"time"
)

var (
	client_id     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	client_secret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	hash_secret   = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
	base_hosts    = "https://app-api.pixiv.net"
	// get your refresh_token, and replace _REFRESH_TOKEN
	//  https://github.com/upbit/pixivpy/issues/158#issuecomment-778919084
	_REFRESH_TOKEN = "vE94AY1QvM8BGcJuA6o1lPFBpfwv8YrDeuaH1AQWZRQ"
	proxyUrl       = "http://127.0.0.1:7890"
)

type PixivClient struct {
	Headers map[string]string
}

func TestPv(t *testing.T) {
	// Instantiate default collector
	c := colly.NewCollector(colly.AllowURLRevisit())
	// Rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:7890")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		t.Log(e.Text)
	})
	c.OnRequest(func(r *colly.Request) {
		localTime := time.Now().Format(time.RFC3339)
		r.Headers.Set("Authorization", "Bearer DW01MWKqPnWUW76bcsI-p0n-FtV-ZYwap-1fi9VFGGQ")
		r.Headers.Set("Accept-Language", "en-us")
		r.Headers.Set("X-Client-Hash", genClientHash(localTime))
		r.Headers.Set("X-Client-Time", localTime)
		r.Headers.Set("app-os", "ios")
		r.Headers.Set("app-os-version", "14.6")
		r.Headers.Set("user-agent", "PixivIOSApp/7.13.3 (iOS 14.6; iPhone13,2)")
	})
	c.Visit("https://www.pixiv.net")
}

func TestLogin(t *testing.T) {
	// auth("user_rsgn7527", "Javajava321", _REFRESH_TOKEN, nil)
	// illust_ranking()
	DownLoad()
}

type authParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

func GetHeaders() map[string]string {
	localTime := time.Now().Format(time.RFC3339)
	px := make(map[string]string)
	px["Accept-Language"] = "en-us"
	px["X-Client-Time"] = localTime
	px["X-Client-Hash"] = genClientHash(localTime)
	px["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	px["Authorization"] = "Bearer DW01MWKqPnWUW76bcsI-p0n-FtV-ZYwap-1fi9VFGGQ"
	return px
}

var a = "https://i.pximg.net/c/360x360_70/img-master/img/2022/06/24/03/26/39/99258375_p0_square1200.jpg"

func DownLoad() {
	exec.Command("rm", "-rf", "./photo/").Run()
	exec.Command("mkdir", "photo").Run()
	referer := "https://app-api.pixiv.net/"

	photos := illust_ranking()
	if photos != nil {
		for _, photo := range photos {
			ext, _ := GetFileExt(photo)

			resp, err := http_util.Client().Proxy(proxyUrl).Headers(map[string]string{
				"Referer": referer,
			}).Get(photo, nil)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			f, err := os.Create("./photo/" + algo_util.UUID() + "." + ext)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			io.Copy(f, resp.Body)
			f.Close()
		}
	}

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
func auth(username, password, refreshToken string, headers map[string]any) {
	localTime := time.Now().Format(time.RFC3339)
	px := PixivClient{}
	px.Headers = make(map[string]string)
	px.Headers["Accept-Language"] = "en-us"
	px.Headers["X-Client-Time"] = localTime
	px.Headers["X-Client-Hash"] = genClientHash(localTime)
	px.Headers["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	// todo may have error
	json, err := http_util.Client().
		Proxy("http://127.0.0.1:7890").
		PostFormGetJSON("https://oauth.secure.pixiv.net"+"/auth/token", map[string]string{
			"client_id":     client_id,
			"client_secret": client_secret,
			"grant_type":    "refresh_token",
			// "username":       username,
			// "password":       password,
			"get_secure_url": "1",
			"refresh_token":  refreshToken,
		})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(json)
}

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, hash_secret)
	return hex.EncodeToString(h.Sum(nil))
}

func illust_ranking() []string {
	host := "https://app-api.pixiv.net"
	url := host + "/v1/illust/ranking"
	mode := "day_male_r18"
	filter := "for_ios"
	json, err := http_util.Client().Proxy(proxyUrl).Headers(GetHeaders()).GetJSON(url, map[string]string{
		"mode":   mode,
		"filter": filter,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil
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
	return photos
}

func TestC(t *testing.T) {
	// handler.DownloadImg("https://i.pximg.net/c/600x1200_90/img-master/img/2017/11/05/14/41/36/65760041_p0_master1200.jpg")
	// handler.Auth()
	m, err := handler.Search("裸足", handler.PixivPageSize)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m)
	photos := make([]string, 0)
	illusts := m["illusts"].([]interface{})
	for _, illust := range illusts {
		image := illust.(map[string]interface{})["meta_single_page"].(map[string]interface{})
		var imgUrl = image["original_image_url"]
		if imgUrl != nil {
			photos = append(photos, imgUrl.(string))
		} else {
			images := illust.(map[string]interface{})["meta_pages"].([]interface{})
			for _, v := range images {
				img := v.(map[string]interface{})
				url := img["image_urls"].(map[string]interface{})["original"]
				if url != nil {
					photos = append(photos, url.(string))
				}
			}
		}
	}
	p := photos[algo_util.RandomInt(0, len(photos))]
	handler.DownloadImg(p)
}

func TestRan(t *testing.T) {
	// t.Log(algo_util.RandomInt(0, 20))
	// t.Log(rand.Int())
	handler.Auth()
}

func TestB(t *testing.T) {
	result, err := handler.Recommend()
	if err != nil {
		t.Error(err)
		return
	}
	photos := make([]string, 0)
	illusts := result["illusts"].([]interface{})
	for _, illust := range illusts {
		image := illust.(map[string]interface{})["meta_single_page"].(map[string]interface{})
		if image["original_image_url"] != nil {
			photos = append(photos, image["original_image_url"].(string))
			handler.DownloadImg(image["original_image_url"].(string))
		}

	}
}
