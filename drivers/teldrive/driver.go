package teldrive

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"

	"github.com/go-resty/resty/v2"
)

type Teldrive struct {
	model.Storage
	Addition
}

func (d *Teldrive) Config() driver.Config {
	return config
}

func (d *Teldrive) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Teldrive) Init(ctx context.Context) error {
	return nil
}

func (d *Teldrive) Drop(ctx context.Context) error {
	return nil
}

func (d *Teldrive) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	url := d.Addition.Address + "/api/files?order=asc&path="
	if dirPath := dir.GetPath(); dirPath == "" {
		url += d.GetRootPath()
	} else {
		url += dirPath
	}
	client := resty.New()

	client.SetHeader("Cookie", "user-session="+d.Addition.Cookies)
	client.SetHeader("Accept", "application/json, text/plain, */*")
	client.SetHeader("Accept-Language", "en-US,en;q=0.9")
	client.SetHeader("Connection", "keep-alive")

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}

	var response FileList
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	var res []model.Obj
	for _, item := range response.Results {
		if item.Type == "folder" {
			res = append(res, &model.Object{
				ID:       item.ID,
				Name:     item.Name,
				Path:     item.Path,
				Size:     item.Size,
				IsFolder: true,
				Modified: item.UpdatedAt,
			})
		} else {
			res = append(res, &model.Object{
				ID:       item.ID,
				Name:     item.Name,
				Path:     item.Path,
				Size:     item.Size,
				IsFolder: false,
				Modified: item.UpdatedAt,
			})
		}
	}

	return res, nil
}

func (d *Teldrive) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		Execute("GET", d.Addition.Address+"/api/auth/session")

	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal(resp.Body(), &session)
	if err != nil {
		return nil, err
	}

	return &model.Link{
		URL: d.Addition.Address + "/api/files/" + file.GetID() + "/" + file.GetName() + "?hash=" + session.Hash,
	}, nil
}

func (d *Teldrive) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	url := d.Addition.Address + "/api/files"
	method := "POST"

	var parentPath string
	if parentDirPath := parentDir.GetPath(); parentDirPath == "" {
		parentPath = d.GetRootPath()
	} else {
		parentPath = parentDirPath
	}

	payload := `{"name":"` + dirName + `","type":"folder","path":"` + parentPath + `"}`

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		SetBody(payload).
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func (d *Teldrive) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	url := d.Addition.Address + "/api/files/movefiles"
	method := "POST"

	payload := `{"files":["` + srcObj.GetID() + `"],"destination":"` + dstDir.GetPath() + `"}`

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		SetBody(payload).
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func (d *Teldrive) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	url := d.Addition.Address + "/api/files/" + srcObj.GetID()
	method := "PATCH"

	var payload string
	if srcObj.IsDir() {
		payload = `{"name":"` + newName + `","type":"folder"}`
	} else {
		payload = `{"name":"` + newName + `","type":"file"}`
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		SetBody(payload).
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func (d *Teldrive) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errs.NotImplement
}

func (d *Teldrive) Remove(ctx context.Context, obj model.Obj) error {
	url := d.Addition.Address + "/api/files/deletefiles"
	method := "POST"

	payload := `{"files":["` + obj.GetID() + `"]}`

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		SetBody(payload).
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	return nil
}

func (d *Teldrive) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	return errs.NotImplement
}

var _ driver.Driver = (*Teldrive)(nil)
