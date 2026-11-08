package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

// VersionInfo is the payload returned by /version. Every field is additive: a
// release may add fields but never removes or repurposes one, so an older
// client keeps decoding a newer response.
type VersionInfo struct {
	Schema      string `json:"schema"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildTime   string `json:"build_time"`
	Modified    bool   `json:"modified"`
	Environment string `json:"environment"`
}

// currentSchema identifies the response shape. Bump it only when a field is
// removed or changes meaning; adding a field keeps the schema version.
const currentSchema = "version.v1"

// BuildVersion assembles a VersionInfo from the build metadata that the Go
// toolchain embeds and from the runtime environment. buildInfo may be nil when
// the binary is built without module information.
func BuildVersion(env string, bi *debug.BuildInfo) VersionInfo {
	v := VersionInfo{
		Schema:      currentSchema,
		Version:     "(devel)",
		Commit:      "unknown",
		BuildTime:   "unknown",
		Environment: env,
	}
	if v.Environment == "" {
		v.Environment = "local"
	}
	if bi != nil {
		if bi.Main.Version != "" {
			v.Version = bi.Main.Version
		}
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if s.Value != "" {
					v.Commit = s.Value
				}
			case "vcs.time":
				if s.Value != "" {
					v.BuildTime = s.Value
				}
			case "vcs.modified":
				v.Modified = s.Value == "true"
			}
		}
	}
	return v
}

// versionHandler serves the running build so a rollout can be verified by asking
// the Service which revision answered the request.
func versionHandler(bi *debug.BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BuildVersion(os.Getenv("APP_ENV"), bi))
	}
}

func main() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		bi = nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/version", versionHandler(bi))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	fmt.Printf("listening on %s version=%s env=%s\n",
		addr, BuildVersion(os.Getenv("APP_ENV"), bi).Version, BuildVersion(os.Getenv("APP_ENV"), bi).Environment)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
