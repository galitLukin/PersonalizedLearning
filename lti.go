// Package lti provides support for working with LTI
// more info can be checked at:
//
// https://www.imsglobal.org/activity/learning-tools-interoperability
//
// Basically it can sign http requests and also it can
// verify incoming LMI requests when acting as a Provider
// This package is WIP. More features will be added when
// needed. Will try to mantain a compatibility API.
package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
  "crypto/sha1"
  "encoding/base64"
  "bytes"

	"github.com/jordic/lti/oauth"
)

const (
	oAuthVersion = "1.0"
	SigHMAC      = "HMAC-SHA1"
	Version      = "0.1"
)

// Provider is an app, that can consume LTI messages,
// also a provider could be used, to construct messages and sign them
//
//  p := lti.NewProvider("secret", "http://url.com")
//  p.Add("param_name", "vale").
//    Add("other_param", "param2")
//
//  sig, err := p.Sign()
//
// will sign, the request, and add the needed fields to the
// Provider.values > Can access it throught p.Params()
// It also can be used to Verify and handle, incoming LTI requests.
//
//  p.IsValid(requesto)
//
// A Provider also holds a internal params url.Values, that can
// be accessed via Get, or Add.
type Provider struct {
	Secret      string
	URL         string
	ConsumerKey string
	Method      string
	values      url.Values
	r           *http.Request
	Signer      oauth.OauthSigner
}

// NewProvider is a provider configured with sensible defaults
// as a signer the HMACSigner is used... (seems that is the most used)
func NewProvider(secret, urlSrv string) *Provider {
	sig := oauth.GetHMACSigner(secret, urlSrv)
	return &Provider{
		Secret: secret,
		Method: "POST",
		values: url.Values{},
		Signer: sig,
		URL:    urlSrv,
	}
}

// HasRole checks if a LTI request, has a provided role
func (p *Provider) HasRole(role string) bool {
	ro := strings.Split(p.Get("roles"), ",")
	roles := strings.Join(ro, " ") + " "
	if strings.Contains(roles, role+" ") {
		return true
	}
	return false
}

// Get a value from the Params map in provider
func (p *Provider) Get(k string) string {
	return p.values.Get(k)
}

// Params returns the map of values stored on the LTI request
func (p *Provider) Params() url.Values {
	return p.values
}

// SetParams for a provider
func (p *Provider) SetParams(v url.Values) *Provider {
	p.values = v
	return p
}

// Add a new param to a LTI request
func (p *Provider) Add(k, v string) *Provider {
	if p.values == nil {
		p.values = url.Values{}
	}
	p.values.Set(k, v)
	return p
}

// Empty checks if a key is defined (or has something)
func (p *Provider) Empty(key string) bool {
	if p.values == nil {
		p.values = url.Values{}
	}
	return p.values.Get(key) == ""
}

// Sign a request, adding, required fields,
// A request, can be drilled on a template, iterating, over p.Prams()
func (p *Provider) Sign() (string, error) {
	if p.Empty("oauth_version") {
		p.Add("oauth_version", oAuthVersion)
	}
	if p.Empty("oauth_timestamp") {
		p.Add("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	}
	if p.Empty("oauth_nonce") {
		p.Add("oauth_nonce", nonce())
	}
	if p.Empty("oauth_signature_method") {
		p.Add("oauth_signature_method", p.Signer.GetMethod())
	}
	p.Add("oauth_consumer_key", p.ConsumerKey)

	signature, err := Sign(p.values, p.URL, p.Method, p.Signer)
	if err == nil {
		p.Add("oauth_signature", signature)
	}
	return signature, err
}

// IsValid returns if lti request is valid, currently only checks
// if signature is correct
func (p *Provider) IsValid(r *http.Request) (bool, error) {
	r.ParseForm()
	p.values = r.Form
  mybody := "<?xml version = \"1.0\" encoding = \"UTF-8\"?><imsx_POXEnvelopeRequest xmlns = \"http://www.imsglobal.org/services/ltiv1p1/xsd/imsoms_v1p0\"><imsx_POXHeader><imsx_POXRequestHeaderInfo><imsx_version>V1.0</imsx_version><imsx_messageIdentifier>999999123</imsx_messageIdentifier></imsx_POXRequestHeaderInfo></imsx_POXHeader><imsx_POXBody><replaceResultRequest><resultRecord><sourcedGUID><sourcedId>course-v1:MITx+15.071x+1T2019:courses.edx.org-a855518774854399b79abee373351e3c:6987787dd79cf0aecabdca8ddae95b4a</sourcedId></sourcedGUID><result><resultScore><language>en-us</language><textString>0.92</textString></resultScore></result></resultRecord></replaceResultRequest></imsx_POXBody></imsx_POXEnvelopeRequest>"
  myreq, err := http.NewRequest("POST", "https://courses.edx.org/courses/course-v1:MITx+15.071x+1T2019/xblock/block-v1:MITx+15.071x+1T2019+type@lti_consumer+block@a855518774854399b79abee373351e3c/handler_noauth/outcome_service_handler", bytes.NewBuffer([]byte(mybody)))
  var body []byte
  if myreq.Body != nil {
    var err error
    body, err = getBody(myreq)
    if err != nil {
      return false, fmt.Errorf("Failed to get body of request ...")
    }
  }
  hasher := sha1.New()
  hasher.Write(body)
  sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
  p = p.Add("oauth_body_hash", sha)

	ckey := p.values.Get("oauth_consumer_key")
	if ckey != p.ConsumerKey {
		return false, fmt.Errorf("Invalid consumer key provided")
	}

  bodyHash := p.values.Get("oauth_body_hash")
	if bodyHash == "" {
		return false, fmt.Errorf("Does not have oauth_body_hash")
	}

	if p.values.Get("oauth_signature_method") != p.Signer.GetMethod() {
		return false, fmt.Errorf("wrong signature method %s",
			r.Form.Get("oauth_signature_method"))
	}
	signature := p.values.Get("oauth_signature")
	// log.Printf("REQuest URLS %s", r.RequestURI)
	sig, err := Sign(p.values, p.URL, r.Method, p.Signer)
	if err != nil {
		return false, err
	}
	if sig == signature {
		return true, nil
	}
	return false, fmt.Errorf("Invalid signature, %s, expected %s", sig, signature)
}

// SetSigner defines the signer that want to use.
func (p *Provider) SetSigner(s oauth.OauthSigner) {
	p.Signer = s
}

// Sign a lti request using HMAC containing a u, url, a http method,
// and a secret. ts is a tokenSecret field from the oauth spec,
// that in this case must be empty.
func Sign(form url.Values, u, method string, firm oauth.OauthSigner) (string, error) {
  // for k, v := range form {
  //   if !string.Contains(k, "oauth"){
  //     delete(form, k);
  //   }
  // }

	fmt.Println("My logging!!!!!")
	fmt.Println(form, u, method, firm)
	str, err := getBaseString(method, u, form)
	fmt.Println("My logging2!!!!!")
	if err != nil {
		return "", err
	}
	fmt.Println(str)
	// log.Printf("Base string: %s", str)
	sig, err := firm.GetSignature(str)
	if err != nil {
		return "", err
	}
	return sig, nil
}

func getBaseString(m, u string, form url.Values) (string, error) {

	var kv []oauth.KV
	for k := range form {
		if k != "oauth_signature" {
			s := oauth.KV{k, form.Get(k)}
			kv = append(kv, s)
		}
	}

	str, err := oauth.GetBaseString(m, u, kv)
	if err != nil {
		return "", err
	}
	// ugly patch for formatting string as expected.
	str = strings.Replace(str, "%2B", "%2520", -1)
	return str, nil
}

var nonceCounter uint64

// nonce returns a unique string.
func nonce() string {
	n := atomic.AddUint64(&nonceCounter, 1)
	if n == 1 {
		binary.Read(rand.Reader, binary.BigEndian, &n)
		n ^= uint64(time.Now().UnixNano())
		atomic.CompareAndSwapUint64(&nonceCounter, 1, n)
	}
	return strconv.FormatUint(n, 16)
}
