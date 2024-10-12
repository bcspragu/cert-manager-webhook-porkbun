package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

type Client struct {
	secretAPIKey string
	apiKey       string
	http         *http.Client
}

const baseURL = "https://api.porkbun.com/api/json/v3/"

func New(secretAPIKey, apiKey string) *Client {
	return &Client{
		secretAPIKey: secretAPIKey,
		apiKey:       apiKey,
		http:         &http.Client{},
	}
}

type BaseRequest struct {
	SecretAPIKey string `json:"secretapikey"`
	APIKey       string `json:"apikey"`
}

type PingResponse struct {
	Status string `json:"status"`
	IP     string `json:"yourIp"`
}

func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	var req bytes.Buffer
	if err := json.NewEncoder(&req).Encode(c.baseRequest()); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := ctxhttp.Post(ctx, c.http, baseURL+"ping", "application/json", &req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit ping endpoint: %w", err)
	}
	defer httpResp.Body.Close()

	var resp *PingResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return resp, nil
}

func (c *Client) baseRequest() *BaseRequest {
	return &BaseRequest{
		SecretAPIKey: c.secretAPIKey,
		APIKey:       c.apiKey,
	}
}

type RetrieveDNSRecordsResponse struct {
	Status  string      `json:"status"`
	Records []DNSRecord `json:"records"`
}

type DNSRecord struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func (c *Client) RetrieveDNSRecordsByDomain(ctx context.Context, domain string) (*RetrieveDNSRecordsResponse, error) {
	if strings.Contains(domain, "/") {
		return nil, errors.New("invalid domain given")
	}

	return c.retrieveDNSRecords(ctx, "dns/retrieve/"+domain)
}

func (c *Client) RetrieveDNSRecordsByDomainSubdomainType(ctx context.Context, domain, subdomain, typ string) (*RetrieveDNSRecordsResponse, error) {
	if strings.Contains(domain, "/") {
		return nil, errors.New("invalid domain given")
	}
	if strings.Contains(subdomain, "/") {
		return nil, errors.New("invalid subdomain given")
	}
	if strings.Contains(typ, "/") {
		return nil, errors.New("invalid type given")
	}
	return c.retrieveDNSRecords(ctx, "dns/retrieveByNameType/"+domain+"/"+typ+"/"+subdomain)
}

func (c *Client) retrieveDNSRecords(ctx context.Context, path string) (*RetrieveDNSRecordsResponse, error) {
	var req bytes.Buffer
	if err := json.NewEncoder(&req).Encode(c.baseRequest()); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := ctxhttp.Post(ctx, c.http, baseURL+path, "application/json", &req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit retrieve DNS records endpoint: %w", err)
	}
	defer httpResp.Body.Close()

	var resp *RetrieveDNSRecordsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return resp, nil
}

type NewDNSRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
}

type CreateDNSRecordRequest struct {
	*BaseRequest
	*NewDNSRecord
}

type CreateDNSRecordResponse struct {
	Status string `json:"status"`
	ID     int    `json:"id"`
}

func (c *Client) CreateDNSRecord(ctx context.Context, domain string, record *NewDNSRecord) (*CreateDNSRecordResponse, error) {
	if strings.Contains(domain, "/") {
		return nil, errors.New("invalid domain given")
	}

	createReq := &CreateDNSRecordRequest{
		BaseRequest:  c.baseRequest(),
		NewDNSRecord: record,
	}
	var req bytes.Buffer
	if err := json.NewEncoder(&req).Encode(createReq); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := ctxhttp.Post(ctx, c.http, baseURL+"dns/create/"+domain, "application/json", &req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit reate DNS record endpoint: %w", err)
	}
	defer httpResp.Body.Close()

	var resp *CreateDNSRecordResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return resp, nil
}

type DeleteDNSRecordResponse struct {
	Status string `json:"status"`
}

func (c *Client) DeleteDNSRecordByDomainID(ctx context.Context, domain, id string) (*DeleteDNSRecordResponse, error) {
	if strings.Contains(domain, "/") {
		return nil, errors.New("invalid domain given")
	}
	if strings.Contains(id, "/") {
		return nil, errors.New("invalid id given")
	}

	var req bytes.Buffer
	if err := json.NewEncoder(&req).Encode(c.baseRequest()); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := ctxhttp.Post(ctx, c.http, baseURL+"dns/delete/"+domain+"/"+id, "application/json", &req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit delete DNS record endpoint: %w", err)
	}
	defer httpResp.Body.Close()

	var resp *DeleteDNSRecordResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return resp, nil
}
