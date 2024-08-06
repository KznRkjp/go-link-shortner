package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/filesio"

	"github.com/stretchr/testify/assert"
)

// полчаем ссылку и отдаем ответ типа text/plain с кодом 201
func Test_mainPage_1(t *testing.T) {
	type args struct {
		code        int
		url         string
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test URL 1",
			args: args{
				code:        409,
				url:         `https://yandex.ru`,
				contentType: "text/plain",
			},
		},
		{
			name: "Test URL 2",
			args: args{
				code:        409,
				url:         `https://google.com`,
				contentType: "text/plain",
			},
		},
		{
			name: "Test URL 3",
			args: args{
				code:        409,
				url:         `https://www.google.com/search?q=golang+tests+best+practices`,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.args.url))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			app.GetURL(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			assert.Equal(t, test.args.contentType, res.Header.Get("Content-Type"))
		})
	}
}

// главаная страница должна возвращать 400 при GET запросе
func Test_mainPage_2(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 400",
			args: args{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			app.GetURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}

func Test_returnURL(t *testing.T) {
	type args struct {
		code     int
		urlPart  string
		location string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test URL Return",
			args: args{
				code:     307,
				urlPart:  "/9JSpJWH612",
				location: "https://test-pass-ok.com",
			},
		},
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Populate DB
			// app.URLDb["9JSpJWH612"] = "https://test-pass-ok.com"
			app.URLDb["9JSpJWH612"] = filesio.URLRecord{ID: 1, ShortURL: "9JSpJWH612", OriginalURL: "https://test-pass-ok.com"}

			request := httptest.NewRequest(http.MethodGet, test.args.urlPart, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			app.ReturnURL(w, request)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)
			// проверяем ответную ссылку
			assert.Equal(t, test.args.location, res.Header.Get("Location"))
		})
	}
}

func TestAPIGetURL(t *testing.T) {
	type args struct {
		code        int
		data        string
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test URL 1",
			args: args{
				code:        409,
				data:        `{"url":"http://mail.ru"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Test URL 2",
			args: args{
				code:        409,
				data:        `{"url":"https://google.com"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Test URL 3",
			args: args{
				code:        409,
				data:        `{"url":"https://www.google.com/search?q=golang+tests+best+practices"}`,
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// req := models.Request{
			// 	URL: test.args.data,
			// }
			// req1, _ := json.Marshal(req)

			// request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(string(req1)))
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.args.data))

			w := httptest.NewRecorder()
			app.APIGetURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}

func BenchmarkMain(b *testing.B) {

}
