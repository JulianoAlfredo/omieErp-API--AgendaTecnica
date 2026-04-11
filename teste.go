package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	url := "https://click.omie.com/pdfnfse-2404j76xjimr"
	filepath := "nf99.pdf"

	fmt.Println("Iniciando download...")
	err := DownloadFile(url, filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Download concluído com sucesso!")
}
