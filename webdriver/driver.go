package webdriver

import (
	"github.com/sclevine/agouti/webdriver/element"
)

type Driver struct {
	Session executable
}

type executable interface {
	Execute(endpoint, method string, body, result interface{}) error
}

type Element interface {
	GetText() (string, error)
	GetAttribute(attribute string) (string, error)
	Click() error
}

type Cookie struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	HTTPOnly bool        `json:"httpOnly"`
	Expiry   int64       `json:"expiry"`
}

func (d *Driver) Navigate(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return d.Session.Execute("url", "POST", request, &struct{}{})
}

func (d *Driver) GetElements(selector string) ([]Element, error) {
	request := struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{"css selector", selector}

	var results []struct{ Element string }

	if err := d.Session.Execute("elements", "POST", request, &results); err != nil {
		return nil, err
	}

	elements := []Element{}
	for _, result := range results {
		elements = append(elements, &element.Element{result.Element, d.Session})
	}

	return elements, nil
}

func (d *Driver) SetCookie(cookie *Cookie) error {
	request := struct {
		Cookie *Cookie `json:"cookie"`
	}{cookie}

	return d.Session.Execute("cookie", "POST", request, &struct{}{})
}

func (d *Driver) GetURL() (string, error) {
	var url string
	if err := d.Session.Execute("url", "GET", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}
