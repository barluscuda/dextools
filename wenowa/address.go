package wenowa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Options carries pluginKey and optional query flags.
type Options struct {
	PluginKey string
	KW        string
	BaseURL   string
	Lang      string
}

func assertPluginKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("pluginKey is required")
	}
	return nil
}

func npmQuery(pluginKey string, opts Options) map[string]string {
	q := map[string]string{"key": pluginKey}
	if strings.TrimSpace(opts.KW) != "" {
		q["kw"] = opts.KW
	}
	if strings.TrimSpace(opts.Lang) != "" {
		q["lang"] = opts.Lang
	}
	return q
}

func getJSON(ctx context.Context, url string, query map[string]string) (any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ErrorFromResponseBody(b, resp.StatusCode, "Request failed")
	}
	var out any
	if len(b) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}

func assertPositiveID(name string, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%s must be a positive number", name)
	}
	return nil
}

// GetProvinces GET /npm-provinces?key=&kw=&lang=
func GetProvinces(ctx context.Context, opts Options) (any, error) {
	return Wenowa{}.GetProvinces(ctx, opts)
}

func (Wenowa) GetProvinces(ctx context.Context, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-provinces"
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}

// GetProvinceById GET /npm-provinces/:id?key=&lang=
func GetProvinceById(ctx context.Context, id int64, opts Options) (any, error) {
	return Wenowa{}.GetProvinceById(ctx, id, opts)
}

func (Wenowa) GetProvinceById(ctx context.Context, id int64, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	if err := assertPositiveID("id", id); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-provinces/" + strconv.FormatInt(id, 10)
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}

// GetDistrictsByProvince GET /npm-districts/by-province/:provinceId?key=&kw=&lang=
func GetDistrictsByProvince(ctx context.Context, provinceID int64, opts Options) (any, error) {
	return Wenowa{}.GetDistrictsByProvince(ctx, provinceID, opts)
}

func (Wenowa) GetDistrictsByProvince(ctx context.Context, provinceID int64, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	if err := assertPositiveID("provinceId", provinceID); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-districts/by-province/" + strconv.FormatInt(provinceID, 10)
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}

// GetDistrictById GET /npm-districts/:id?key=&lang=
func GetDistrictById(ctx context.Context, id int64, opts Options) (any, error) {
	return Wenowa{}.GetDistrictById(ctx, id, opts)
}

func (Wenowa) GetDistrictById(ctx context.Context, id int64, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	if err := assertPositiveID("id", id); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-districts/" + strconv.FormatInt(id, 10)
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}

// GetVillagesByDistrict GET /npm-vilages/by-district/:districtId?key=&kw=&lang=
func GetVillagesByDistrict(ctx context.Context, districtID int64, opts Options) (any, error) {
	return Wenowa{}.GetVillagesByDistrict(ctx, districtID, opts)
}

func (Wenowa) GetVillagesByDistrict(ctx context.Context, districtID int64, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	if err := assertPositiveID("districtId", districtID); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-vilages/by-district/" + strconv.FormatInt(districtID, 10)
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}

// GetVillageById GET /npm-vilages/:id?key=&lang=
func GetVillageById(ctx context.Context, id int64, opts Options) (any, error) {
	return Wenowa{}.GetVillageById(ctx, id, opts)
}

func (Wenowa) GetVillageById(ctx context.Context, id int64, opts Options) (any, error) {
	if err := assertPluginKey(opts.PluginKey); err != nil {
		return nil, err
	}
	if err := assertPositiveID("id", id); err != nil {
		return nil, err
	}
	base := ResolveBaseURL(opts.BaseURL)
	url := base + "/npm-vilages/" + strconv.FormatInt(id, 10)
	return getJSON(ctx, url, npmQuery(opts.PluginKey, opts))
}
