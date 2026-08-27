// Command genresources rebuilds internal/resource/bundled.json.gz from the
// official CDN. Run it when the game ships new config tables:
//
//	go run ./tools/genresources -version 59 -mainscene eab4a -delayres 30c5b
//
// The bundle hashes come from the packaged game's settings.*.json.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alttab8520/qqfarm-sdk/internal/resource"
)

func main() {
	out := flag.String("out", "internal/resource/bundled.json.gz", "output path")
	version := flag.String("version", "", "game pack version, recorded in the bundle")
	mainscene := flag.String("mainscene", "", "mainscene bundle hash")
	delayRes := flag.String("delayres", "", "delayRes bundle hash")
	flag.Parse()

	vers := map[string]string{}
	for k, v := range resource.DefaultBundleVers {
		vers[k] = v
	}
	if *mainscene != "" {
		vers["mainscene"] = *mainscene
	}
	if *delayRes != "" {
		vers["delayRes"] = *delayRes
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	blob, err := resource.BuildBundle(ctx, vers, *version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "生成失败:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, blob, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "写入失败:", err)
		os.Exit(1)
	}
	fmt.Printf("%s %d bytes\n", *out, len(blob))
}
