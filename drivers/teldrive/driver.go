package teldrive

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"

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
		dirPath, err := d.sanitizeHTMLURL(dirPath)
		if err != nil {
			return nil, err
		}
		url += dirPath
	}

	utils.Log.Info(url)
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
	fileSize := stream.GetSize()
	fileName := stream.GetName()
	input := fmt.Sprintf("%s:%s:%d", fileName, dstDir, fileSize)

	hash := md5.Sum([]byte(input))
	hashString := hex.EncodeToString(hash[:])

	uploadURL := fmt.Sprintf("/api/uploads/%s", hashString)

	url := d.Addition.Address + uploadURL
	method := "POST"

	var uploadFile UploadFile
	var uploadedSize int64 = 0

	if len(uploadFile.Parts) != 0 {
		for _, part := range uploadFile.Parts {
			uploadedSize += part.Size
		}
	}

	client := resty.New()

	utils.Log.Info(uploadedSize)
	utils.Log.Info(stream.GetSize())
	if uploadedSize != stream.GetSize() {

		in := bufio.NewReader(stream)

		if uploadedSize > 0 {
			io.CopyN(io.Discard, in, uploadedSize)
		}

		left := stream.GetSize() - uploadedSize

		partNo := 1

		if len(uploadFile.Parts) > 0 {
			partNo = len(uploadFile.Parts) + 1
		}

		totalParts := int(math.Ceil(float64(stream.GetSize()) / float64(1024*1024*1024)))

		for {

			if _, err := in.Peek(1); err != nil {
				if left > 0 {
					return err
				}
				break
			}
			n := int64(1024 * 1024 * 1024)
			if stream.GetSize() != -1 {
				n = d.int64min(left, n)
				left -= n
			}
			partReader := io.LimitReader(in, n)

			name := fmt.Sprintf("%s.part.%03d", stream.GetName(), partNo)
			resp, err := client.R().
				SetHeader("Content-Type", "application/octet-stream").
				SetHeader("Cookie", "user-session="+d.Addition.Cookies).
				SetHeader("Accept-Language", "en-US,en;q=0.9").
				SetHeader("Connection", "keep-alive").
				SetBody(partReader).
				SetHeader("Content-Length", strconv.FormatInt(n, 10)).
				SetQueryParams(map[string]string{
					"fileName":   name,
					"partNo":     strconv.Itoa(partNo),
					"totalparts": strconv.FormatInt(int64(totalParts), 10),
				}).
				Execute(method, url)

			utils.Log.Info(resp.String())
			if err != nil {
				return err
			}
			if resp.StatusCode() != http.StatusOK {
				return errors.New(resp.String())
			}
			var response PartFile
			err = json.Unmarshal(resp.Body(), &response)
			if err != nil {
				return err
			}

			uploadFile.Parts = append(uploadFile.Parts, response)
			partNo++
		}
	}

	fileParts := []FilePart{}

	for _, part := range uploadFile.Parts {
		fileParts = append(fileParts, FilePart{ID: part.PartId})
	}

	// upload
	payload := FileUploadRequest{
		Name:     stream.GetName(),
		MimeType: stream.GetMimetype(),
		Type:     "file",
		Parts:    fileParts,
		Size:     int(stream.GetSize()),
		Path:     dstDir.GetPath(),
	}
	resp_file, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		SetBody(payload).
		Execute(method, d.Addition.Address+"/api/files")

	if err != nil {
		return err
	}
	if resp_file.StatusCode() != http.StatusOK {
		return errors.New(resp_file.String())
	}
	utils.Log.Info(resp_file.String())

	resp_del_temp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", "user-session="+d.Addition.Cookies).
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "keep-alive").
		Execute("DELETE", d.Addition.Address+"/api/uploads/"+hashString)
	if err != nil {
		return err
	}
	if resp_del_temp.StatusCode() != http.StatusOK {
		return errors.New(resp_del_temp.String())
	}

	return nil
}

var _ driver.Driver = (*Teldrive)(nil)
