// Package api 提供 Steam API 通用请求与数据解析的核心方法
// Package api provides core methods for Steam API general request and data parsing
package api

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
	ue "github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
)

// GetRawBytes 执行API请求并返回序列化后的字节数组
// GetRawBytes executes an API request and returns the serialized byte array
//
// 参数说明 (Parameters):
//
//	c - 客户端实例, 用于发起HTTP请求 (Client instance for initiating HTTP requests)
//	method - HTTP请求方法（如GET/POST/PUT等） (HTTP request method (e.g. GET/POST/PUT))
//	url - 请求的目标URL地址 (Target URL address for the request)
//	params - URL查询参数 (URL query parameters)
//
// 返回值 (Returns):
//
//	respBytes - API响应序列化后的字节数组 (Serialized byte array of API response)
//	err - 执行过程中的错误，包含：
//	      1. 客户端请求错误(如网络错误、超时)
//	      2. 响应数据序列化错误(包装ue.ErrAPIResponse)
//	err - Error during execution, including:
//	      1. Client request errors (e.g. network error, timeout)
//	      2. Response data serialization error (wrapped with ue.ErrAPIResponse)
func GetRawBytes(c *client.Client, method, url string, params url.Values) (respBytes []byte, err error) {
	// 执行请求
	resp, err := c.DoRequest(method, url, params)
	if err != nil {
		return respBytes, err
	}

	// 转换为字节
	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", ue.ErrAPIResponse, err)
	}
	return respBytes, nil
}

// GetRawModel 通用API请求方法, 将响应数据反序列化为指定类型的结构体
// GetRawModel is a generic API request method that deserializes response data into a struct of the specified type
//
// 泛型参数 (Generic Parameters):
//
//	T - 目标反序列化的结构体类型，需满足：
//	    1. 字段首字母大写(可导出)
//	    2. 字段需匹配API响应的JSON结构(可通过json tag映射)
//	T - Target struct type for deserialization, which must satisfy:
//	    1. Field names are uppercase (exportable)
//	    2. Fields match the JSON structure of API response (mappable via json tag)
//
// 参数说明 (Parameters):
//
//	c - 客户端实例，用于发起HTTP请求 (Client instance for initiating HTTP requests)
//	method - HTTP请求方法（如GET/POST/PUT等） (HTTP request method (e.g. GET/POST/PUT))
//	reqUrl - 请求的目标URL地址 (Target URL address for the request)
//	params - URL查询参数 (URL query parameters)
//
// 返回值 (Returns):
//
//	T - 反序列化后的目标类型实例，错误时返回该类型的零值
//	err - 执行过程中的错误，包含：
//	      1. GetRawBytes返回的所有错误类型
//	      2. 响应数据反序列化错误(包装ue.ErrAPIResponse)
//	T - Instance of the target type after deserialization, returns zero value of the type on error
//	err - Error during execution, including:
//	      1. All error types returned by GetRawBytes
//	      2. Response data deserialization error (wrapped with ue.ErrAPIResponse)
func GetRawModel[T any](c *client.Client, method, reqUrl string, params url.Values) (T, error) {
	// 定义零值
	var zero T

	// 获取原始字节数据
	bytes, err := GetRawBytes(c, method, reqUrl, params)
	if err != nil {
		return zero, err
	}

	// 反序列化为传入的泛型类型
	var resp T
	if err = sonic.Unmarshal(bytes, &resp); err != nil {
		return zero, fmt.Errorf("%w: unmarshal %T resp failed: %v", ue.ErrAPIResponse, resp, err)
	}

	return resp, nil
}
