package xss

// XXX TODO
// bluemonday!
// also have option to accept XSS into the database and filter it out on display

import (
	//"errors"
	"github.com/gin-gonic/gin"
	//"net/http/httputil" // debugging
	//"net/http"
	//"strings"
	"bytes"
	//"time"
	"encoding/json"
	"fmt"
	"html"
	//"html"
	//"io"
	"io/ioutil"
	//"net/url"
	//"os"
	"strconv"
)

//type nopCloser struct {
//	io.Reader
//}
//
//func (nopCloser) Close() error { return nil }

// GinXSSMiddleware provides an 'auto' remove XSS malicious from all submitted user input.
// e.g. POST and PUT
// it's highly configurable.
// uses HTML sanitizer https://github.com/microcosm-cc/bluemonday

type GinXSSMiddleware struct {
	// List of tables to not filter any fields on
	TableWhitelist []byte

	// List of fields to not filter . i.e. created_on, created_at, etc
	FieldWhitelist []byte

	// Hash of table->field combinations to skip filtering on
	TableFieldWhitelist []byte

	// Config options - how much filtering? Regular | Lite | Custom

	//PayloadFunc func(userID string) map[string]interface{}
	// User can define own Unauthorized func.
	//Unauthorized func(*gin.Context, int, string)

}

//// MiddlewareInit initialize jwt configs.
//func (mw *GinXSSMiddleware) MiddlewareInit() error {
//
//	if mw.TokenLookup == "" {
//		mw.TokenLookup = "header:Authorization"
//	}
//
//	if mw.Authorizator == nil {
//		mw.Authorizator = func(userID string, c *gin.Context) bool {
//			return true
//		}
//	}
//
//	if mw.Unauthorized == nil {
//		mw.Unauthorized = func(c *gin.Context, code int, message string) {
//			c.JSON(code, gin.H{
//				"code":    code,
//				"message": message,
//			})
//		}
//	}
//
//	if mw.Key == nil {
//		return errors.New("secret key is required")
//	}
//
//	return nil
//}

// MiddlewareFunc makes GinXSSMiddleware implement the Middleware interface.
func (mw *GinXSSMiddleware) MiddlewareFunc() gin.HandlerFunc {
	//if err := mw.MiddlewareInit(); err != nil {
	//	return func(c *gin.Context) {
	//		mw.unauthorized(c, http.StatusInternalServerError, err.Error())
	//		return
	//	}
	//}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *GinXSSMiddleware) middlewareImpl(c *gin.Context) {
	// remove xss
	err := mw.filterData(c)

	if err != nil {
		//mw.exception(c, http.StatusUnauthorized, err.Error())
		//c.Abort()
		//return
		c.Abort()
		return
	}

	//token, err := mw.parseToken(c)

	//if err != nil {
	//	mw.unauthorized(c, http.StatusUnauthorized, err.Error())
	//	return
	//}

	//claims := token.Claims.(jwt.MapClaims)

	//id := claims["id"].(string)
	//c.Set("JWT_PAYLOAD", claims)
	//c.Set("userID", id)

	//if !mw.Authorizator(id, c) {
	//	mw.unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
	//	return
	//}

	c.Next()
}

//func (mw *GinXSSMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
//	authHeader := c.Request.Header.Get(key)
//
//	if authHeader == "" {
//		return "", errors.New("auth header empty")
//	}
//
//	parts := strings.SplitN(authHeader, " ", 2)
//	if !(len(parts) == 2 && parts[0] == "Bearer") {
//		return "", errors.New("invalid auth header")
//	}
//
//	return parts[1], nil
//}
//
//func (mw *GinXSSMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
//	token := c.Query(key)
//
//	if token == "" {
//		return "", errors.New("Query token empty")
//	}
//
//	return token, nil
//}
//
//func (mw *GinXSSMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
//	cookie, _ := c.Cookie(key)
//
//	if cookie == "" {
//		return "", errors.New("Cookie token empty")
//	}
//
//	return cookie, nil
//}
//
//func (mw *GinXSSMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
//	var token string
//	var err error
//
//	parts := strings.Split(mw.TokenLookup, ":")
//	switch parts[0] {
//	case "header":
//		token, err = mw.jwtFromHeader(c, parts[1])
//	case "query":
//		token, err = mw.jwtFromQuery(c, parts[1])
//	case "cookie":
//		token, err = mw.jwtFromCookie(c, parts[1])
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
//			return nil, errors.New("invalid signing algorithm")
//		}
//
//		return mw.Key, nil
//	})
//}
//
//func (mw *GinXSSMiddleware) unauthorized(c *gin.Context, code int, message string) {
//
//	if mw.Realm == "" {
//		mw.Realm = "gin jwt"
//	}
//
//	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
//	c.Abort()
//
//	mw.Unauthorized(c, code, message)
//
//	return
//}

func (mw *GinXSSMiddleware) exception(c *gin.Context, code int, message string) {

	//c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	//mw.Unauthorized(c, code, message)

	return
}

// Remove XSS
//<nil>"PUT /api/v1/projects/1 HTTP/1.1\r\nHost: local.hubtones.com:11062\r\n
//Accept: application/json\r\nAccept-Encoding: gzip, deflate, sdch\r\n
// Accept-Language: en-US,en;q=0.8\r\nAuthorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODA2NzM2NDQsImlkIjoiVGVzdCBVc2VyIMKpIiwib3JpZ19pYXQiOjE0ODA0MDc2MzEsInVzZXJfaWQiOjUzNywidXNlcm5hbWUiOiJUZXN0IFVzZXIgwqkifQ.vYIdRkO6c4oatRZ4-gUL068RlWam-XO1XZrU6SGA8U0\r\n
// Connection: keep-alive\r\n
// Content-Length: 283\r\n
// Content-Type: application/json\r\n
// Origin: http://local.hubtones.com\r\n
// Referer: http://local.hubtones.com/project/1/edit\r\n
// User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36\r\n\r\n
// {\"id\":1,\"name\":\"Project ß£áçkqùë Jâçqùë ¥ \",\"description\":\"Iñtërnâtiônàlizætiøn project\",\"status\":\"Recording\",\"genre\":\"7\",\"sub_genre\":\"77\",\"bpm\":\"117\",\"key\":\"E\",\"visibility\":\"Public\",\"created_by\":537,\"created_at\":1474448233,\"updated_by\":537,\"updated_at\":1480493923}"
// [GIN] 2016/11/30 - 00:25:38 | 200 |    2.717269ms | 127.0.0.1 |   PUT     /api/v1/projects/1
func (mw *GinXSSMiddleware) filterData(c *gin.Context) error {
	//dump, derr := httputil.DumpRequest(c.Request, true)
	//fmt.Print(derr)
	//fmt.Printf("%q", dump)

	//ReqHeader := c.Request.Header
	//fmt.Printf("%v Header\n", ReqHeader)

	//// https://golang.org/pkg/net/http/#Request

	ReqMethod := c.Request.Method
	//fmt.Printf("%v Method\n", ReqMethod)

	//ReqURL := c.Request.URL
	//fmt.Printf("%v URL\n", ReqURL)

	ReqBody := c.Request.Body
	//fmt.Printf("%v URL\n", ReqBody)

	ct_hdr := c.Request.Header.Get("Content-Type") // [application/json]
	//fmt.Printf("%v\n", ct_hdr)                     // -> application/json

	cts_len := c.Request.Header.Get("Content-Length")
	//fmt.Printf("%v\n", cts_len)
	ct_len, _ := strconv.Atoi(cts_len)

	// https://golang.org/src/net/http/request.go
	// set expected application type
	if ct_hdr == "application/json" && ct_len > 1 && (ReqMethod == "POST" || ReqMethod == "PUT") {
		var jsonBod interface{}
		//jsnErr := json.NewDecoder(ReqBody).Decode(&jsonBod)
		d := json.NewDecoder(ReqBody)
		d.UseNumber()
		jsnErr := d.Decode(&jsonBod)
		if jsnErr == nil {
			var buff bytes.Buffer
			buff.WriteString(`{`)

			m := jsonBod.(map[string]interface{})
			for k, v := range m {
				fmt.Println(k, v)
				buff.WriteString(`"`)
				buff.WriteString(k)
				buff.WriteString(`":`)

				// FYI, json is string or float
				switch vv := v.(type) {
				case string:
					fmt.Println(k, "is string", vv)
					buff.WriteString(`"`)
					// TODO
					// XXX need to escape [ "`{},: ]
					//buff.WriteString(vv)
					// XXX to do  bluemonday!
					buff.WriteString(html.EscapeString(vv))
					buff.WriteString(`",`)
				case float64:
					fmt.Println(k, "is float", vv)
					//buff.WriteString(strconv.FormatFloat(vv, 'g', 0, 64))
					buff.WriteString(html.EscapeString(strconv.FormatFloat(vv, 'g', 0, 64)))
					buff.WriteString(`,`)
				default:
					// XXX talent_ids [1] is an array of values (handle it!)
					// talent_ids is of a type I don't know how to handle

					fmt.Println(k, "is of a type I don't know how to handle")
					fmt.Println("%#v", vv)
					//buff.WriteString(fmt.Sprintf("%v", vv))
					buff.WriteString(html.EscapeString(fmt.Sprintf("%v", vv)))
					buff.WriteString(`,`)
				}
			}
			buff.Truncate(buff.Len() - 1) // remove last ','
			buff.WriteString(`}`)

			////encBuf, merr := json.Marshal(buff.String(), encOpts{quoted: true})
			//encBuf, merr := json.Marshal(buff.String())
			//if merr != nil {
			//	fmt.Printf("%v", merr)
			//}
			//////fmt.Println(string(data))
			////c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(string(encBuf))))

			bodOut := buff.String()

			enc := json.NewEncoder(ioutil.Discard)
			if merr := enc.Encode(&bodOut); merr != nil {
				fmt.Printf("%v", merr)
			}

			fmt.Printf("ReqBody PRE: %v\n", ReqBody)
			//bf := `{"genre":"7","created_at":88812334,"updated_by":534,"updated_at":12344,"bpm":"117","key":"E","visibility": "Public","id":1,"name":"Project ß£áçkqùë Jâçqùë ¥  - value asdfasdfadfs","description": "Iñtërnâtiônàlizætiøn project  asdfasdf","status":"Recording","sub_genre":"77","created_by":534}`
			//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(bf)))
			//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(buff.String())))
			//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(string(encBuf))))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(bodOut)))

			fmt.Printf("ReqBody Post: %v\n", c.Request.Body)
			fmt.Printf("ReqBody Post: %#v\n", c.Request.Body)
		} else {
			fmt.Println("Failed")
		}

	}
	//return errors.New("Filter error")
	return nil

}

//func ConstructRequest(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println(r.Method, "-", r.RequestURI)
//		cookie, _ := r.Cookie("username")
//		if cookie != nil {
//			//Add data to context
//			ctx := context.WithValue(r.Context(), "Username", cookie.Value)
//			next.ServeHTTP(w, r.WithContext(ctx))
//		} else {
//			next.ServeHTTP(w, r)
//		}
//	})
//}

// func AddContext(next http.Handler) http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     log.Println(r.Method, "-", r.RequestURI)
//     cookie, _ := r.Cookie("username")
//     if cookie != nil {
//       //Add data to context
//       ctx := context.WithValue(r.Context(), "Username", cookie.Value)
//       next.ServeHTTP(w, r.WithContext(ctx))
//     } else {
//       next.ServeHTTP(w, r)
//     }
//   })
// }

// var postCustomHandle map[string]interface{}
// jsnErr := json.NewDecoder(c.Request.Body).Decode(&postCustomHandle)
// if jsnErr == nil {
//         proj_id := postCustomHandle["project_id"]
//         talent_ids := postCustomHandle["talent_ids"]
//         if proj_id == nil || talent_ids == nil {
//                 c.JSON(http.StatusBadRequest, gin.H{"message": "required paramaters cannot be empty."})
//                 return
//         } else {
//                 //project_id := postCustomHandle["proj_id"].(string)
//                 project_id = proj_id.(string)
//         }
// } else {
//         cglog.CGlog.Println(c.Request.Body)
//         cglog.CGlog.Printf("%T", c.Request.Body)
//         cglog.CGlog.Print(jsnErr)
//         c.JSON(http.StatusBadRequest, gin.H{"message": "error handling parameters"})
//         return
// }

//func (mw *GinXSSMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
//	authHeader := c.Request.Header.Get(key)
//
//	if authHeader == "" {
//		return "", errors.New("auth header empty")
//	}
//
//	parts := strings.SplitN(authHeader, " ", 2)
//	if !(len(parts) == 2 && parts[0] == "Bearer") {
//		return "", errors.New("invalid auth header")
//	}
//
//	return parts[1], nil
//}
//
//func (mw *GinXSSMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
//	token := c.Query(key)
//
//	if token == "" {
//		return "", errors.New("Query token empty")
//	}
//
//	return token, nil
//}
//
