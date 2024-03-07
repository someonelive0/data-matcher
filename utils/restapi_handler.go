package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type RestapiHandler struct {
	Name    string
	Runtime time.Time
}

func (p *RestapiHandler) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // should before w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, Version(p.Name))
}

func (p *RestapiHandler) DebugHandler(w http.ResponseWriter, r *http.Request) {
	s := ""
	if strings.ToUpper(r.Method) == "GET" {
		s = log.GetLevel().String()
	} else if strings.ToUpper(r.Method) == "PUT" {
		params := r.URL.Query()
		verbose, ok := params["verbose"]
		if ok {
			if verbose[0] == "true" {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}
		} else {
			log.SetLevel(log.DebugLevel)
		}
		s = log.GetLevel().String()
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s)
}

func (p *RestapiHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s := fmt.Sprintf(`{"app": "%s", `, p.Name)
	claims, ok := r.Context().Value(ContextKeyRequestID).(*MyClaims)
	if ok && claims != nil {
		b, _ := json.Marshal(claims)
		s += fmt.Sprintf(`"jwt": %s, `, b)
	}
	s += fmt.Sprintf(`"run_time": "%s"}`, p.Runtime.Format(time.RFC3339))
	fmt.Fprint(w, s)
}

func (p *RestapiHandler) ChangelogHandler(w http.ResponseWriter, r *http.Request) {
	if err := p.OutputFile(w, "Changelog.md"); err != nil {
		fmt.Fprintf(w, "get file [Changelog.md] failed: %v", err)
	}
}

func (p *RestapiHandler) LogHandler(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("log/%s.log.%s", p.Name, time.Now().Format("20060102"))
	if err := p.OutputFile(w, filename); err != nil {
		fmt.Fprintf(w, "get file [%s] failed: %v", filename, err)
	}
}

func (p *RestapiHandler) OutputFile(w http.ResponseWriter, filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(w, fp); err != nil {
		return err
	}

	return nil
}

func (p *RestapiHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		//fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
		if part.FileName() == "" { // this is FormData
			data, _ := io.ReadAll(part)
			log.Debugf("  formdata: [%s]=[%s]\n", part.FormName(), string(data))
		} else { // This is FileData
			filename := "./etc/" + path.Base(part.FileName()) + ".upload.tmp"
			// dir := path.Dir(part.FileName())
			// if dir != "." {
			// 	if err := os.MkdirAll(dir, 0755); err != nil {
			// 		log.Errorf("RestapiHandler mkdir [%s] failed: %s\n", dir, err)
			// 		http.Error(w, err.Error(), http.StatusInternalServerError)
			// 		break
			// 	}
			// }
			dst, err := os.Create(filename)
			if err != nil {
				log.Errorf("RestapiHandler create file [%s] failed: %s", filename, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			defer dst.Close()
			filelen, err := io.Copy(dst, part)
			if err != nil {
				log.Errorf("RestapiHandler write file [%s] failed: %s", filename, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			dst.Close()
			if err := os.Rename(filename, filename[0:len(filename)-4]); err != nil {
				log.Errorf("RestapiHandler rename file [%s] failed: %s", filename, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			filename = filename[0 : len(filename)-4]
			log.Infof("RestapiHandler receive file %s success, len %d", filename, filelen)
			fmt.Fprintf(w, "RestapiHandler receive file %s success, len %d\n", filename, filelen)
		}
	}
}

// JWT auth middleware for http router
// Context 类型保存 KV 对时, key 不能使用原生类型，而应该使用派生类型。
// 所以采用 ContextKeyRequestID 作为 context 的key, 而不是原来原生的字符串。
type ContextKey int

const ContextKeyRequestID ContextKey = iota

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			tokenString := authHeader[1]
			myjwt := NewMyJwt(nil)
			claims, err := myjwt.ParseToken(tokenString)
			if err != nil {
				log.Warnf("AuthMiddleware parse token failed: %s", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			} else {
				ctx := context.WithValue(r.Context(), ContextKeyRequestID, claims)
				// ctx := context.WithValue(r.Context(), "claims", claims)
				//log.Debug(claims)
				// Access context values in handlers like this
				//props, _ := r.Context().Value("claims").(*MyClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

	})
}
