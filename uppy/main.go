package main

import (
    "os"
	"io"
    "fmt"
    "path"
	"mime"
	"bytes"
	"embed"
	"strings"
	"net/http"
	"math/rand"
	"html/template"
    "github.com/u1f408/x"
)

//go:embed template static
var embedFS embed.FS

var (
    Version string
	Config UppyConfig

    flags = x.NewFlagSet(path.Base(os.Args[0]), Version)
    debugFlag = flags.F.Bool("log.debug", false, "enable debug logging")
	configFile = flags.F.String("config.file", "uppy.yml", "path to configuration file")
	configEnv = flags.F.Bool("config.env", false, "load configuration from environment")

	embedTemplate = template.Must(template.ParseFS(embedFS, "template/*.tmpl"))
)

func Logf(format string, a ...any) (int, error) {
    format = fmt.Sprintf("[%s] %s\n", path.Base(os.Args[0]), format)
    return fmt.Fprintf(os.Stderr, format, a...)
}

func Debugf(format string, a ...any) (int, error) {
    if *debugFlag {
        format = fmt.Sprintf("DEBUG: %s", format)
        return Logf(format, a...)
    }

    return 0, nil
}

func GenIdent(ident string) string {
	return x.EncBase32Lower.EncodeString(x.ShortHashf("$Ident$%s", ident))
}

func GenFilename(ext string) string {
	r := make([]byte, 18)
	rand.Read(r)

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return x.EncBase62Modified.EncodeString(r) + ext
}

func UploadFile(content *bytes.Buffer, mimeType, dir, filename string) (string, error) {
	if Config.BunnyCdn.IsValid() {
		return Config.BunnyCdn.UploadFile(content, mimeType, dir, filename)
	}

	return "", fmt.Errorf("couldn't find a suitable upload destination!")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// 10MB max upload size
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		Logf("error getting form file: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	io.Copy(&buf, file)
	file.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	ext := ".bin"
	if idx := strings.LastIndex(fileHeader.Filename, "."); idx > 0 {
		ext = fileHeader.Filename[idx:]
	} else {
		possible, _ := mime.ExtensionsByType(mimeType)
		if len(possible) > 0 {
			ext = possible[0]
		}
	}

	newDir := "+" + GenIdent(r.Header.Get("X-Uppy-User"))
	newFile := GenFilename(ext)

	Debugf("UploadFile(<len:%d>), %s, %s, %s)", buf.Len(), mimeType, newDir, newFile)
	path, err := UploadFile(&buf, mimeType, newDir, newFile)
	if err != nil {
		Logf("error uploading file to destination: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusFound)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	embedTemplate.ExecuteTemplate(w, "index.tmpl", map[string]interface{}{
		"User": r.Header.Get("X-Uppy-User"),
		"Meow": "meow!",
	})
}

func main() {
	var err error
	flags.ParseFromOS()

	if *configEnv {
		err = ParseConfigFromEnv(&Config)
	} else {
		err = ParseConfigFromFile(&Config, *configFile)
	}

	if err != nil {
		Logf("failed to load config: %s", err.Error())
        os.Exit(2)
	}

	Debugf("Config.Auth: %+v", Config.Auth)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(embedFS)))
	mux.HandleFunc("POST /upload", handleUpload)
	mux.HandleFunc("GET /", handleIndex)

	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		Debugf("%s %s %s from %s", r.Proto, r.Method, r.URL.Path, r.RemoteAddr)

		r.Header.Set("X-Uppy-User", r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")])
		if Config.Auth.Enable && len(Config.Auth.HttpHeader) > 0 {
			hval := r.Header.Get(Config.Auth.HttpHeader)
			if len(hval) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			hdomain := strings.TrimPrefix(hval[strings.LastIndex(hval, "@"):], "@")
			domainOk := false
			if len(Config.Auth.AllowedDomains) == 0 {
				domainOk = true
			}
			for _, d := range Config.Auth.AllowedDomains {
				if hdomain == d {
					domainOk = true
				}
			}

			if !domainOk {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-Uppy-User", hval)
		}

		mux.ServeHTTP(w, r)
	})

	Logf("listening on %s, port %d", Config.ListenHost, Config.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", Config.ListenHost, Config.ListenPort), handler)
	Logf("%s", err.Error())
}
