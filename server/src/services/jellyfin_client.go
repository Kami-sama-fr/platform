package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type JellyfinClient struct {
	baseURL    string
	apiKey     string
	userID     string
	httpClient *http.Client
}

type JellyfinConfig struct {
	URL    string
	APIKey string
	UserID string
}

func NewJellyfinClient(cfg JellyfinConfig) *JellyfinClient {
	return &JellyfinClient{
		baseURL: strings.TrimRight(cfg.URL, "/"),
		apiKey:  cfg.APIKey,
		userID:  cfg.UserID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *JellyfinClient) Name() string { return "jellyfin" }

func (c *JellyfinClient) do(ctx context.Context, method, path string, params url.Values) ([]byte, error) {
	u := c.baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Emby-Authorization", fmt.Sprintf(`MediaBrowser Client="KamiSama", Device="Server", DeviceId="kamisama-server", Version="1.0.0", Token="%s"`, c.apiKey))
	if c.apiKey != "" {
		q := req.URL.Query()
		q.Set("api_key", c.apiKey)
		req.URL.RawQuery = q.Encode()
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("jellyfin %s %s: %d: %s", method, path, resp.StatusCode, string(body))
	}
	return body, nil
}

type jfUser struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type jfUserViewsResponse struct {
	Items []jfView `json:"Items"`
}

type jfView struct {
	ID       string `json:"Id"`
	Name     string `json:"Name"`
	ItemType string `json:"Type"`
}

type jfItemsResponse struct {
	Items            []jfItem `json:"Items"`
	TotalRecordCount int      `json:"TotalRecordCount"`
}

type jfItem struct {
	ID                 string         `json:"Id"`
	Name               string         `json:"Name"`
	OriginalTitle      string         `json:"OriginalTitle"`
	Type               string         `json:"Type"`
	Year               int            `json:"ProductionYear"`
	OfficialRating     string         `json:"OfficialRating"`
	Overview           string         `json:"Overview"`
	CommunityRating    float64        `json:"CommunityRating"`
	SeriesName         string         `json:"SeriesName"`
	ParentIndexNumber  *int           `json:"ParentIndexNumber"`
	IndexNumber        *int           `json:"IndexNumber"`
	SeasonName         string         `json:"SeasonName"`
	RecursiveItemCount int            `json:"RecursiveItemCount"`
	MediaStreams       []jfMediaStream `json:"MediaStreams"`
	ImageTags          map[string]string `json:"ImageTags"`
	RunTimeTicks       int64          `json:"RunTimeTicks"`
	Container          string         `json:"Container"`
	Genres             []string       `json:"Genres"`
	SeriesId           string         `json:"SeriesId"`
	ParentId           string         `json:"ParentId"`
}

type jfMediaStream struct {
	Type         string  `json:"Type"`
	Codec        string  `json:"Codec"`
	Width        int     `json:"Width"`
	Height       int     `json:"Height"`
	BitRate      int64   `json:"BitRate"`
	Language     string  `json:"Language"`
	DisplayTitle string  `json:"DisplayTitle"`
}

type jfPlaybackInfoResponse struct {
	MediaStreams []jfMediaStream `json:"MediaSources"`
}

func ticksToSeconds(ticks int64) float64 {
	return float64(ticks) / 10000000.0
}

func mapJFItemToMediaSourceItem(item jfItem) map[string]interface{} {
	result := map[string]interface{}{
		"id":            item.ID,
		"sourceId":      item.ID,
		"name":          item.Name,
		"originalTitle": item.OriginalTitle,
		"type":          item.Type,
		"year":          item.Year,
		"rating":        item.CommunityRating,
		"overview":      item.Overview,
		"duration":      ticksToSeconds(item.RunTimeTicks),
		"container":     item.Container,
	}
	if item.SeriesName != "" {
		result["name"] = item.SeriesName
		result["originalTitle"] = item.OriginalTitle
	}
	if item.ParentIndexNumber != nil {
		result["seasonNumber"] = *item.ParentIndexNumber
	}
	if item.IndexNumber != nil {
		result["episodeNumber"] = *item.IndexNumber
	}
	if len(item.Genres) > 0 {
		result["genres"] = item.Genres
	}
	for _, ms := range item.MediaStreams {
		switch ms.Type {
		case "Video":
			result["videoCodec"] = ms.Codec
			result["width"] = ms.Width
			result["height"] = ms.Height
			result["bitrate"] = ms.BitRate
		case "Audio":
			if result["audioCodec"] == nil {
				result["audioCodec"] = ms.Codec
			}
		}
	}
	raw, _ := json.Marshal(item)
	result["rawMetadata"] = datatypes.JSON(raw)
	return result
}

func (c *JellyfinClient) ListLibraries(ctx context.Context) ([]map[string]interface{}, error) {
	viewsResp, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Views", c.userID), nil)
	if err != nil {
		return nil, err
	}
	var views jfUserViewsResponse
	if err := json.Unmarshal(viewsResp, &views); err != nil {
		return nil, err
	}
	var libs []map[string]interface{}
	for _, v := range views.Items {
		libs = append(libs, map[string]interface{}{
			"id":       v.ID,
			"name":     v.Name,
			"type":     v.ItemType,
			"itemCount": 0,
		})
	}
	return libs, nil
}

func (c *JellyfinClient) GetLibrary(ctx context.Context, id string) (map[string]interface{}, error) {
	viewsResp, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Views", c.userID), nil)
	if err != nil {
		return nil, err
	}
	var views jfUserViewsResponse
	if err := json.Unmarshal(viewsResp, &views); err != nil {
		return nil, err
	}
	for _, v := range views.Items {
		if v.ID == id {
			return map[string]interface{}{
				"id":       v.ID,
				"name":     v.Name,
				"type":     v.ItemType,
				"itemCount": 0,
			}, nil
		}
	}
	return nil, fmt.Errorf("library %s not found", id)
}

func (c *JellyfinClient) ListItems(ctx context.Context, libraryID string, limit, offset int, sortBy, query string) ([]map[string]interface{}, int, error) {
	params := url.Values{}
	params.Set("ParentId", libraryID)
	params.Set("Limit", strconv.Itoa(limit))
	params.Set("StartIndex", strconv.Itoa(offset))
	params.Set("Recursive", "true")
	params.Set("Fields", "MediaStreams,Genres,Path")
	params.Set("IncludeItemTypes", "Series,Episode,Movie")
	if sortBy != "" {
		params.Set("SortBy", sortBy)
	}
	if query != "" {
		params.Set("SearchTerm", query)
	}
	body, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Items", c.userID), params)
	if err != nil {
		return nil, 0, err
	}
	var resp jfItemsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, err
	}
	items := make([]map[string]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, mapJFItemToMediaSourceItem(item))
	}
	return items, resp.TotalRecordCount, nil
}

func (c *JellyfinClient) GetItem(ctx context.Context, id string) (map[string]interface{}, error) {
	body, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Items/%s", c.userID, id), nil)
	if err != nil {
		return nil, err
	}
	var item jfItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	return mapJFItemToMediaSourceItem(item), nil
}

func (c *JellyfinClient) SearchItems(ctx context.Context, query string, limit int) ([]map[string]interface{}, error) {
	params := url.Values{}
	params.Set("SearchTerm", query)
	params.Set("Limit", strconv.Itoa(limit))
	params.Set("Recursive", "true")
	params.Set("IncludeItemTypes", "Series,Episode,Movie")
	params.Set("Fields", "MediaStreams,Genres")
	body, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Items", c.userID), params)
	if err != nil {
		return nil, err
	}
	var resp jfItemsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, mapJFItemToMediaSourceItem(item))
	}
	return items, nil
}

func (c *JellyfinClient) GetStreamURL(ctx context.Context, itemID string, static bool) (string, error) {
	params := url.Values{}
	if static {
		params.Set("static", "true")
	} else {
		params.Set("static", "true")
	}
	return c.baseURL + "/Videos/" + itemID + "/stream?" + params.Encode() + "&api_key=" + c.apiKey, nil
}

func (c *JellyfinClient) GetPlaybackInfo(ctx context.Context, itemID string) (map[string]interface{}, error) {
	body, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/Users/%s/Items/%s", c.userID, itemID), nil)
	if err != nil {
		return nil, err
	}
	var item jfItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"duration": ticksToSeconds(item.RunTimeTicks),
		"container": item.Container,
		"subtitles": []map[string]interface{}{},
	}
	for _, ms := range item.MediaStreams {
		switch ms.Type {
		case "Video":
			result["videoCodec"] = ms.Codec
			result["width"] = ms.Width
			result["height"] = ms.Height
			result["bitrate"] = ms.BitRate
		case "Audio":
			if result["audioCodec"] == nil {
				result["audioCodec"] = ms.Codec
			}
		case "Subtitle":
			subs := result["subtitles"].([]map[string]interface{})
			result["subtitles"] = append(subs, map[string]interface{}{
				"language": ms.Language,
				"format":   ms.Codec,
				"title":    ms.DisplayTitle,
			})
		}
	}
	streamURL, err := c.GetStreamURL(ctx, itemID, true)
	if err != nil {
		return nil, err
	}
	result["streamUrl"] = streamURL
	return result, nil
}

func (c *JellyfinClient) ReportPlaybackProgress(ctx context.Context, itemID string, positionTicks int64, stopped bool) error {
	params := url.Values{}
	params.Set("ItemId", itemID)
	params.Set("PositionTicks", strconv.FormatInt(positionTicks, 10))
	params.Set("IsPaused", strconv.FormatBool(stopped))
	_, err := c.do(ctx, http.MethodPost, "/Sessions/Playing/Progress", params)
	return err
}
