package presentation

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func logging(handler http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		recorder := &responseRecorder{
			ResponseWriter: w,
			Status:         200,
			ContentLength:  0,
		}
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		handler.ServeHTTP(recorder, r)

		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info("Request info")

		log.WithFields(logrus.Fields{
			"status":         recorder.Status,
			"content length": recorder.ContentLength,
		}).Info("Response info")

	}
	return http.HandlerFunc(logFn)
}

func compress(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			compressWriter := newCompressWriter(w)
			originalWriter = compressWriter
			defer compressWriter.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		handler.ServeHTTP(originalWriter, r)
	})
}

func auth(handlerFn AuthenticatedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		UserID := jose.ParseUserID(token)
		if UserID == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handlerFn(w, r, *UserID)
	})
}
