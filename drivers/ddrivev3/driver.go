package ddrivev3

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"errors"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/go-resty/resty/v2"
)

type Ddrivev3 struct {
	model.Storage
	Addition
}

func (d *Ddrivev3) Config() driver.Config {
	return config
}

func (d *Ddrivev3) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Ddrivev3) Init(ctx context.Context) error {
	return nil
}

func (d *Ddrivev3) Drop(ctx context.Context) error {
	return nil
}

func (d *Ddrivev3) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	url := d.Addition.Address + dir.GetPath()

	client := resty.New()

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}

	var response DirectoryResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	var res []model.Obj
	for _, item := range response.Directory {
		res = append(res, &model.Object{
			ID:       item.ID,
			Name:     item.Name,
			Path:     item.Path,
			Size:     0,
			IsFolder: true,
			Modified: time.Now(),
		})
	}
	for _, item := range response.Files {
		res = append(res, &model.Object{
			ID:       item.ID,
			Name:     item.Name,
			Path:     item.Path,
			Size:     int64(item.Size),
			IsFolder: false,
			Modified: time.Now(),
		})
	}

	return res, nil
}

func (d *Ddrivev3) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	return &model.Link{
		URL: d.Addition.Address + "/" + file.GetPath(),
	}, nil
}

func (d *Ddrivev3) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	url := d.Addition.Address + "/"
	var parentPath string
	if parentDirPath := parentDir.GetPath(); parentDirPath == "" {
		parentPath = dirName
	} else {
		parentPath = parentDirPath + "/" + dirName
	}
	method := "PUT"

	client := resty.New()

	resp, err := client.R().
		Execute(method, url+parentPath)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	utils.Log.Debug(resp.String())
	return nil

}

func (d *Ddrivev3) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errs.NotSupport
}

func (d *Ddrivev3) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	return errs.NotSupport
}

func (d *Ddrivev3) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errs.NotSupport
}

func (d *Ddrivev3) Remove(ctx context.Context, obj model.Obj) error {
	url := d.Addition.Address + obj.GetPath()

	method := "DELETE"

	client := resty.New()

	resp, err := client.R().
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	utils.Log.Debug(resp.String())
	return nil
}

func (d *Ddrivev3) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	url := d.Addition.Address + dstDir.GetPath() + "/" + stream.GetName()

	method := "POST"

	payload := stream

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

var _ driver.Driver = (*Ddrivev3)(nil)
