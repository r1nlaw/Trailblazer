package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
)

type Landmark struct {
	repo repository.Landmark
	config.ParserConfig
	cookies   []*http.Cookie
	ApiConfig config.GeocoderConfig
}

func NewLandmarkService(landmark repository.Landmark, cfg config.ParserConfig) *Landmark {
	return &Landmark{
		repo:         landmark,
		ParserConfig: cfg,
	}
}

func (s *Landmark) Crawl() error {
	if err := s.SetCookies(); err != nil {
		slog.Error(err.Error())
		return err
	}
	pageStr := s.GetHtml(s.ParserConfig.BaseURL)
	var urls = make([]string, 0)
	urls = s.GetCardURL(pageStr)

	landmarks := make([]*models.Landmark, 0)
	for _, url := range urls {
		landmark := s.GetInfoByLink(url)
		landmarks = append(landmarks, landmark)
	}
	var err error
	err = s.repo.SaveLandmarks(context.Background(), landmarks)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (s *Landmark) SaveLandmarks(c context.Context, landmarks []*models.Landmark) error {
	//TODO implement me
	panic("implement me")
}

func (s *Landmark) SetCookies() error {
	resp, err := http.Get(s.ParserConfig.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}
	defer resp.Body.Close()

	s.cookies = resp.Cookies()
	return nil
}

func (s *Landmark) GetHtml(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	for _, cookie := range s.cookies {
		req.AddCookie(cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	return string(body)
}

func (s *Landmark) GetCardURL(str string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	url := make([]string, 0)
	//landmarks := []*models.{}
	sel := doc.Find(".full-card-content")
	for i := range sel.Nodes {
		a, ok := sel.Eq(i).Find("a").First().Attr("href")
		if !ok {
			slog.Error("fasf")
		}
		url = append(url, a)
		//landmark := s.GetInfoByLink(a)
		//landmark.Name = sel.Eq(i).Find("h3>a").Text()
		//landmark.Name = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(landmark.Name, "")
		//landmarks = append(landmarks, landmark)
	}
	doc.Find("ul.wp-block-list").Find("a").Each(func(i int, s *goquery.Selection) {
		a, ok := s.First().Attr("href")
		if !ok {
			slog.Error(err.Error())
			return
		}
		url = append(url, a)
	})
	return url
}
func (s *Landmark) GetInfoByLink(link string) *models.Landmark {
	page := s.GetHtml(link)
	slog.Info(link)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	sel := doc.Find(".promo-text")
	var landmark = &models.Landmark{}
	sel.Find("p").Each(
		func(i int, selection *goquery.Selection) {
			label := selection.Find("strong").Text()
			switch label {
			case "Адрес:":
				value := strings.TrimSpace(strings.Replace(selection.Text(), label, "", 1))
				landmark.Address = value
			case "Координаты:":
				value := strings.TrimSpace(strings.Replace(selection.Text(), label, "", 1))
				value = strings.Replace(value, " ", "", -1)
				coord := strings.Split(value, ",")

				landmark.Location.Lat, err = strconv.ParseFloat(coord[0], 64)
				if err != nil {
					slog.Error(err.Error())
					return
				}

				landmark.Location.Lng, err = strconv.ParseFloat(coord[1], 64)
				if err != nil {
					slog.Error(err.Error())
					return
				}
			}
		})
	history := doc.Find("section.wp-block-rest-gutenberg-blocks-guten-block.content-block")
	historyText := ""
	history.Find("p").Each(
		func(i int, selection *goquery.Selection) {
			historyText += selection.Text()
		})
	landmark.Name = doc.Find(".promo-layout-columns").Find("h1").Text()
	landmark.History = historyText
	return landmark
}
