package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Complain struct {
	Id         string
	UserName   string
	UserMail   string
	Category   Category
	Issue      string
}

type Category struct {
	ID            string
	Proves        string
	File          string
	Amount        int
	BankName      string
	BankAccount   string
	Remark        string `validate:"required,max=255"`
	Name          string
	NumberAccount int
}

var (
	dbComplain = Database()
	emptyString = ""
)

func SubmitForm(w http.ResponseWriter, r *http.Request)  {
	complain := Complain{
		Id:       uuid.New().String(),
		UserName: r.FormValue("Uname"),
		UserMail: r.FormValue("Umail"),
		Issue:    r.FormValue("Uissue"),
	}
	kategori := complain.Category
	kategori.ID = r.FormValue("idCategory")
	if kategori.ID  == "categoryA" {
		file, name, errFile := CheckFile(r)
		if errFile != nil{
			http.Redirect(w, r, "/", http.StatusBadRequest)
		}
		kategori.File = file
		kategori.Proves = name
		Amount, _ := strconv.Atoi(r.FormValue("Amount"))
		kategori.Amount = Amount
		kategori.BankName = r.FormValue("BankName")
		kategori.BankAccount = r.FormValue("BankAccount")
		kategori.Remark = r.FormValue("Remark")
		complain.Category = kategori
	}else {
		kategori.Name = r.FormValue("NameB")
		kategori.BankName = r.FormValue("BankNameB")
		NumberAccount, _ := strconv.Atoi(r.FormValue("NumberAccountB"))
		kategori.NumberAccount = NumberAccount
		complain.Category = kategori
	}
	fmt.Println(complain)
	dbComplain.Write("complain", complain.Id, complain)

}
func CheckFile(r *http.Request ) (string, string, error) {
	file, header, err :=  r.FormFile("Proves")
	defer file.Close()
	if err != nil {
		return emptyString, emptyString, err
	}
	data, errs := ioutil.ReadAll(file)
	if errs != nil {
		log.Println(errs)
		return emptyString, emptyString, errs
	}
	if header.Size > (1 * 1024 * 1024) {
		log.Printf("invalid request ", "The maximum upload file is 1MB ", header.Size)
		return emptyString , emptyString, http.ErrMissingFile
	}
	contentType := http.DetectContentType(data)

	switch contentType {
	case "image/png":
		fmt.Println("Image type is already PNG.")
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return "", "", fmt.Errorf("unable to decode jpeg: %w", err)
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return "", "", fmt.Errorf("unable to encode png: %w", err)
		}
		data = buf.Bytes()
	default:
		return "", "", fmt.Errorf("unsupported content typo: %s", contentType)
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(data)
	return imgBase64Str, header.Filename, nil
}