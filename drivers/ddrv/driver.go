package ddrv

import (
	"context"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"

	"encoding/json"
	"time"

	"errors"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/go-resty/resty/v2"
)

type Ddrv struct {
	model.Storage
	Addition
}

func (d *Ddrv) Config() driver.Config {
	return config
}

func (d *Ddrv) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Ddrv) Init(ctx context.Context) error {
	// TODO login / refresh token
	//op.MustSaveDriverStorage(d)
	return nil
}

func (d *Ddrv) Drop(ctx context.Context) error {
	return nil
}

func (d *Ddrv) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	var url string
	if strings.Contains(d.Addition.Address, ",") {
		urlList := strings.Split(d.Addition.Address, ",")
		randomIndex := rand.Intn(len(urlList))
		url = urlList[randomIndex] + "/api/directories/" + dir.GetID()
	} else {

		url = d.Addition.Address + "/api/directories/" + dir.GetID()
	}

	client := resty.New()
	client.SetAuthToken(d.Addition.Token)

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	var res []model.Obj
	for _, item := range response.Data.Files {
		if !item.IsDir {
			res = append(res, &model.Object{
				ID:       item.ID,
				Name:     item.Name,
				Path:     item.Parent,
				Size:     int64(item.Size),
				IsFolder: false,
				Modified: time.Now(),
			})
		} else {
			res = append(res, &model.Object{
				ID:       item.ID,
				Name:     item.Name,
				Path:     item.Parent,
				Size:     0,
				IsFolder: true,
				Modified: time.Now(),
			})
		}
	}
	return res, nil
}

func (d *Ddrv) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	var url string
	if strings.Contains(d.Addition.Address, ",") {
		urlList := strings.Split(d.Addition.Address, ",")
		randomIndex := rand.Intn(len(urlList))
		url = urlList[randomIndex]
	} else {

		url = d.Addition.Address
	}

	return &model.Link{
		URL: url + "/files/" + file.GetID(),
	}, nil
}

func (d *Ddrv) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	var url string
	if strings.Contains(d.Addition.Address, ",") {
		urlList := strings.Split(d.Addition.Address, ",")
		randomIndex := rand.Intn(len(urlList))
		url = urlList[randomIndex] + "/api/directories/"
	} else {

		url = d.Addition.Address + "/api/directories/"
	}

	method := "POST"

	payload := `{"name": "` + dirName + `", "parent": "` + parentDir.GetID() + `"}`

	client := resty.New()
	client.SetAuthToken(d.Addition.Token)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Execute(method, url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return errors.New(resp.String())
	}
	utils.Log.Debug(resp.String())
	return nil
}

func (d *Ddrv) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	if srcObj.IsDir() {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + srcObj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + srcObj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + srcObj.GetID()
		method := "PUT"

		payload := `{"name": "` + srcObj.GetName() + `", "parent": "` + dstDir.GetID() + `"}`

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Execute(method, url)

		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.String())
		}
		utils.Log.Debug(resp.String())
		return nil
	} else {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		method := "PUT"

		payload := `{"name": "` + srcObj.GetName() + `", "parent": "` + dstDir.GetID() + `"}`

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
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
}

func (d *Ddrv) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	if srcObj.IsDir() {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + srcObj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + srcObj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + srcObj.GetID()
		method := "PUT"

		payload := `{"name": "` + newName + `", "parent": "` + srcObj.GetPath() + `"}`

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Execute(method, url)

		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.String())
		}
		utils.Log.Debug(resp.String())
		return nil
	} else {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + srcObj.GetPath() + "/files/" + srcObj.GetID()
		method := "PUT"

		payload := `{"name": "` + newName + `", "parent": "` + srcObj.GetPath() + `"}`

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
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
}

func (d *Ddrv) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	// TODO copy obj, optional
	return errs.NotSupport
}

func (d *Ddrv) Remove(ctx context.Context, obj model.Obj) error {

	if obj.IsDir() {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + obj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + obj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + obj.GetID()
		method := "DELETE"

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			Execute(method, url)

		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.String())
		}
		utils.Log.Debug(resp.String())
	} else {
		var url string
		if strings.Contains(d.Addition.Address, ",") {
			urlList := strings.Split(d.Addition.Address, ",")
			randomIndex := rand.Intn(len(urlList))
			url = urlList[randomIndex] + "/api/directories/" + obj.GetPath() + "/files/" + obj.GetID()
		} else {

			url = d.Addition.Address + "/api/directories/" + obj.GetPath() + "/files/" + obj.GetID()
		}
		// url := d.Addition.Address + "/api/directories/" + obj.GetPath() + "/files/" + obj.GetID()
		method := "DELETE"

		client := resty.New()
		client.SetAuthToken(d.Addition.Token)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			Execute(method, url)

		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.String())
		}
		utils.Log.Debug(resp.String())
	}
	return nil
}

func (d *Ddrv) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	const chunkSize = 20 * 1024 * 1024
	var url string
	if strings.Contains(d.Addition.Address, ",") {
		urlList := strings.Split(d.Addition.Address, ",")
		randomIndex := rand.Intn(len(urlList))
		url = urlList[randomIndex] + "/api/directories/" + dstDir.GetID() + "/files"
	} else {

		url = d.Addition.Address + "/api/directories/" + dstDir.GetID() + "/files"
	}
	// url := d.Addition.Address + "/api/directories/" + dstDir.GetID() + "/files"

	// Create the pipe
	pr, pw := io.Pipe()

	// Create a new multipart writer
	bodyWriter := multipart.NewWriter(pw)

	// Create a goroutine to copy the file data to the pipe in chunks
	go func() {
		defer pw.Close()

		// Create a new form file field
		formFile, err := bodyWriter.CreateFormFile("file", stream.GetName())
		if err != nil {
			utils.Log.Info(err)
			return
		}

		// Read and copy the file data in chunks
		buffer := make([]byte, chunkSize)
		for {
			n, _ := stream.Read(buffer)
			if n == 0 {
				break
			}
			formFile.Write(buffer[:n])
		}

		// Close the multipart writer after writing the form data
		bodyWriter.Close()
	}()

	// Create a new HTTP request
	request, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return err
	}

	// Set the Content-Type header to the multipart form data
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	// Set the Authorization header
	request.Header.Set("Authorization", "Bearer "+d.Addition.Token)

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response body
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return errors.New(string(data))
	}

	return nil
}

//func (d *Ddrv) Other(ctx context.Context, args model.OtherArgs) (interface{}, error) {
//	return nil, errs.NotSupport
//}

var _ driver.Driver = (*Ddrv)(nil)
