package dataflow

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/artisancloud/httphelper/client"
	"github.com/artisancloud/httphelper/driver/nethttp"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	http2 "net/http"
	"strings"
	"testing"
	"time"
)

func InitBaseDataflow() *Dataflow {
	c, err := nethttp.NewHttpClient(&client.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	df := NewDataflow(c, nil, &Option{
		BaseUrl: "https://www.baidu.com",
	})
	return df
}

func TestDataflow_WithContext(t *testing.T) {
	df := InitBaseDataflow()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{}, 1)

	go func() {
		time.Sleep(time.Second * 1)
		_, err := df.WithContext(ctx).Request()
		if !errors.Is(err, context.Canceled) {
			t.Error("cancel failed")
		}
		done <- struct{}{}
	}()

	cancel()
	select {
	case <-done:
	}
}

func TestDataflow_Method(t *testing.T) {
	df := InitBaseDataflow()

	df.Method(http2.MethodGet)

	if df.request.Method != http2.MethodGet {
		t.Error("method diff")
	}
}

func TestDataflow_Header(t *testing.T) {
	df := InitBaseDataflow()

	df.Header("content-type", "application/json")

	if df.request.Header.Get("content-type") != "application/json" {
		t.Error("set header failed")
	}
}

func TestDataflow_Json(t *testing.T) {
	df := InitBaseDataflow()

	var data = map[string]interface{}{
		"a": "b",
		"c": map[string]interface{}{
			"c1": "c2",
		},
	}
	df.Json(data)

	jsonBytes, _ := json.Marshal(data)
	bodyBytes, _ := io.ReadAll(df.request.Body)

	// trim body 控制字符
	if string(jsonBytes) != strings.TrimSpace(string(bodyBytes)) {
		t.Error("json body failed")
	}
}

type CaseXmlNode struct {
	A string   `xml:"a"`
	B []string `xml:"b"`
}

type CaseXmlDoc struct {
	Node1 CaseXmlNode `xml:"node1"`
	Node2 CaseXmlNode `xml:"node2"`
}

func TestDataflow_Xml(t *testing.T) {
	df := InitBaseDataflow()

	data := CaseXmlDoc{
		Node1: CaseXmlNode{
			A: "1",
			B: []string{"1", "2"},
		},
		Node2: CaseXmlNode{
			A: "3",
			B: []string{"3", "4"},
		},
	}
	df.Xml(data)

	xmlBytes, _ := xml.Marshal(data)
	bodyBytes, _ := io.ReadAll(df.request.Body)

	// trim body 控制字符
	if string(xmlBytes) != strings.TrimSpace(string(bodyBytes)) {
		t.Error("xml body failed")
	}
}

func TestDataflow_Multipart(t *testing.T) {
	df := InitBaseDataflow()

	df.Multipart(func(multipart MultipartDataflow) error {
		mpDataflow := NewMultipartHelper()
		mpDataflow.Boundary("test-boundary")
		mpDataflow.FieldValue("param1", "value1")
		data := strings.NewReader("it's a string reader")
		mpDataflow.Field("data", data)
		if mpDataflow.Err() != nil {
			t.Error(mpDataflow.Err())
		}
		return nil
	})

	if df.Err() != nil {
		t.Error(df.Err())
	}
}

type TestStruct struct {
	Name  string `form:"name" query:"name"`
	Email string `form:"email" query:"email"`
}

func TestBindQuery(t *testing.T) {
	df := InitBaseDataflow()

	testStruct := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	df.BindQuery(testStruct)

	query := df.request.URL.Query()

	assert.Equal(t, "John Doe", query.Get("name"))
	assert.Equal(t, "john.doe@example.com", query.Get("email"))

	testMap := map[string]string{
		"name":  "Jane Doe",
		"email": "jane.doe@example.com",
	}

	df.BindQuery(testMap)

	query = df.request.URL.Query()

	assert.Equal(t, "Jane Doe", query.Get("name"))
	assert.Equal(t, "jane.doe@example.com", query.Get("email"))
}
