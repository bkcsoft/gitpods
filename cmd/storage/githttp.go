package main

// Much of this code originates from https://github.com/AaronO/go-git-http
// Licensed under Apache-2.0

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/sourcepods/sourcepods/pkg/api"
)

var allowedServices = [2]string{"upload-pack", "receive-pack"}

type GitHTTP struct {
	root   string
	git    string
	Logger log.Logger
}

func NewGitHTTP(root string) *GitHTTP {
	return &GitHTTP{
		git:    "/usr/bin/git",
		root:   root,
		Logger: log.NewNopLogger(),
	}
}

func (gh *GitHTTP) Handler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(api.NewRequestLogger(gh.Logger))

	r.Get("/{owner}/{name}/HEAD", noCaching(gh.headHandler))
	r.Get("/{owner}/{name}/info/refs", noCaching(serviceAllowed(gh.infoRefsHandler)))
	r.Get("/{owner}/{name}/objects/{folder:[0-9a-f]{2}}/{file:[0-9a-f]{38}}", cacheForever(gh.looseObjectHandler))
	r.Get("/{owner}/{name}/objects/info/{file}", noCaching(gh.infoHandler))
	r.Get("/{owner}/{name}/objects/info/alternates", noCaching(gh.alternatesHandler))
	r.Get("/{owner}/{name}/objects/info/http-alternates", noCaching(gh.httpAlternatesHandler))
	r.Get("/{owner}/{name}/objects/info/packs", cacheForever(gh.infoPacksHandler))
	r.Get("/{owner}/{name}/objects/pack/pack-{hash:[0-9a-f]{40}}.idx", cacheForever(gh.idxHandler))
	r.Get("/{owner}/{name}/objects/pack/pack-{hash:[0-9a-f]{40}}.pack", cacheForever(gh.packHandler))
	r.Post("/{owner}/{name}/git-{service}", serviceAllowed(gh.serviceHandler))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		level.Debug(gh.Logger).Log(
			"msg", "not found",
			"path", r.URL.String(),
		)
	})

	return r
}

func serviceAllowed(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := chi.URLParam(r, "service")
		if service == "" {
			service = serviceQuery(r)
		}

		for _, v := range allowedServices {
			if v == service {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, fmt.Sprintf("invalid service %q", service), http.StatusBadRequest)
		// TODO: inject logging?
		//level.Warn(logger).Log("msg", "invalid service", "err", err)
	}
}

func noCaching(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
		next.ServeHTTP(w, r)
	}
}

func cacheForever(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		expires := now.AddDate(1, 0, 0)
		w.Header().Set("Date", fmt.Sprintf("%d", now.Unix()))
		w.Header().Set("Expires", fmt.Sprintf("%d", expires.Unix()))
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	}
}

func (gh *GitHTTP) headHandler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.headHandler")
	defer span.Finish()

	owner, name := ownerName(r)
	path := filepath.Join(gh.root, owner, name, "HEAD")

	h := gh.textFileHandler(path, "text/plain")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) infoRefsHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	service := serviceQuery(r)

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.infoRefsHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("service", service)
	defer span.Finish()

	logger := log.With(gh.Logger,
		"owner", owner,
		"name", name,
		"service", service,
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	args := []string{service, "--stateless-rpc", "--advertise-refs", "."}
	cmd := exec.CommandContext(ctx, gh.git, args...)
	cmd.Dir = filepath.Join(gh.root, owner, name)

	refs, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to get refs", "err", err)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	w.Write(packetWrite(fmt.Sprintf("# service=git-%s\n", service)))
	w.Write(packetFlush())
	w.Write(refs)
}

func (gh *GitHTTP) looseObjectHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	folder := chi.URLParam(r, "folder")
	file := chi.URLParam(r, "file")
	path := filepath.Join(gh.root, owner, name, folder, file)

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.looseObjectHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)

	h := gh.textFileHandler(path, "application/x-git-loose-object")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) infoHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	file := chi.URLParam(r, "file")
	path := filepath.Join(gh.root, owner, name, "objects", "info", file)

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.infoHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "text/plain")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) alternatesHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	path := filepath.Join(gh.root, owner, name, "/objects/info/alternates")

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.alternatesHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "text/plain")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) httpAlternatesHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	path := filepath.Join(gh.root, owner, name, "/objects/info/http-alternates")

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.httpAlternatesHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "text/plain")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) infoPacksHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	path := filepath.Join(gh.root, owner, name, "/objects/info/packs")

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.infoPacksHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "text/plain; charset=utf-8")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) idxHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	hash := chi.URLParam(r, "hash")
	path := filepath.Join(gh.root, owner, name, fmt.Sprintf("objects/pack/pack-%s.idx", hash))

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.idxHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "application/x-git-packed-objects-toc")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) packHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	hash := chi.URLParam(r, "hash")
	path := filepath.Join(gh.root, owner, name, fmt.Sprintf("objects/pack/pack-%s.pack", hash))

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.packHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("path", path)
	defer span.Finish()

	h := gh.textFileHandler(path, "application/x-git-packed-objects")
	h.ServeHTTP(w, r.WithContext(ctx))
}

func (gh *GitHTTP) textFileHandler(path string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.textFileHandler")
		span.SetTag("path", path)
		defer span.Finish()

		f, err := os.Stat(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", f.Size()))
		w.Header().Set("Last-Modified", f.ModTime().Format(http.TimeFormat))
		http.ServeFile(w, r.WithContext(ctx), path)
	}
}

func (gh *GitHTTP) serviceHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	service := chi.URLParam(r, "service")

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "githttp.serviceHandler")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("service", service)
	defer span.Finish()

	logger := log.With(gh.Logger,
		"owner", owner,
		"name", name,
		"service", service,
	)

	defer r.Body.Close()

	var body io.Reader
	var err error
	switch r.Header.Get("content-encoding") {
	case "gzip":
		body, err = gzip.NewReader(r.Body)
	case "deflate":
		body = flate.NewReader(r.Body)
	default:
		body = r.Body
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to create body reader", "err", err)
		return
	}

	rpcSpan, rpcCtx := opentracing.StartSpanFromContext(ctx, "githttp.serviceHandler.rpc")
	defer rpcSpan.Finish()

	args := []string{service, "--stateless-rpc", "."}
	cmd := exec.CommandContext(rpcCtx, gh.git, args...)
	cmd.Dir = filepath.Join(gh.root, owner, name)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to create pipe to git's stdin", "err", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to create pipe to git's stdout", "err", err)
		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to start git", "err", err)
		return
	}

	if _, err := io.Copy(stdin, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to copy request to git's stdin", "err", err)
		return
	}
	stdin.Close()

	if _, err := io.Copy(w, stdout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to copy git's stdout to response", "err", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		level.Warn(logger).Log("msg", "failed to wait for git", "err", err)
		return
	}

	// TODO: Fire events to channel

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))
}

func ownerName(r *http.Request) (string, string) {
	return chi.URLParam(r, "owner"), chi.URLParam(r, "name")
}

func serviceQuery(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Query().Get("service"), "git-")
}

func packetFlush() []byte {
	return []byte("0000")
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)

	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}

	return []byte(s + str)
}
