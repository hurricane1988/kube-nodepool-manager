/*
Copyright 2025 CodeFuture Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"github.com/fatih/color"
	"github.io/codeFuthure/kube-nodepool-manager/pkg/constants"
	"runtime"
	"sigs.k8s.io/yaml"
	"strings"
)

// These are set during build time via -ldflags
var (
	Version   = "N/A"
	GitCommit = "N/A"
	BuildDate = "N/A"
	Author    = constants.DefaultAuthor
)

var (
	Yellow       = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	YellowItalic = color.New(color.FgHiYellow, color.Bold, color.Italic).SprintFunc()
	Green        = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	Blue         = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	Cyan         = color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	Red          = color.New(color.FgHiRed, color.Bold).SprintFunc()
	White        = color.New(color.FgWhite).SprintFunc()
	WhiteBold    = color.New(color.FgWhite, color.Bold).SprintFunc()
	forceDetail  = "yaml"
)

// Info holds the version information of the driver
type Info struct {
	Author       string `json:"Author"`
	Version      string `json:"Version"`
	GitCommit    string `json:"Git Commit"`
	BuildDate    string `json:"Build Date"`
	GoVersion    string `json:"Go Version"`
	Compiler     string `json:"Compiler"`
	Platform     string `json:"Platform"`
	KubeVersion  string `json:"KubernetesVersion"`
	RuntimeCores int    `json:"RuntimeCores"`
	TotalMem     int    `json:"TotalMem"`
}

// GetVersion returns the version information of the driver
func GetVersion(kubeVersion string) Info {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Info{
		Author:       Author,
		Version:      Version,
		GitCommit:    GitCommit,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		KubeVersion:  kubeVersion,
		RuntimeCores: runtime.NumCPU(),
		TotalMem:     int(memStats.TotalAlloc / 1024),
	}
}

// GetVersionYAML returns the version information of the driver
// in YAML format
func GetVersionYAML(kubeVersion string) (string, error) {
	info := GetVersion(kubeVersion)
	marshalled, err := yaml.Marshal(&info)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(marshalled)), nil
}

// Print the version information.
func Print(kubeVersion string) {
	v := GetVersion(kubeVersion)
	fmt.Printf(`
----------------------------------------------
 Author: %s
 Version: %s
 GitCommit: %s
 BuildDate: %s
 GoVersion: %s
 Compiler: %s
 Platform: %s
 RuntimeCores: %d cores
 TotalMem: %d KB
 KubeVersion: %s
----------------------------------------------
`, v.Author, v.Version, v.GitCommit, v.BuildDate, v.GoVersion, v.Compiler, v.Platform, v.RuntimeCores, v.TotalMem, v.KubeVersion)
}

// Term print the terminal logo information.
func Term() string {
	return fmt.Sprint(Blue(`
╭╮╭━╮╱╱╭╮╱╱╱╱╱╱╱╭━╮╱╭╮╱╱╱╱╭╮╱╱╱╱╱╱╱╱╱╱╱╭╮╱╱╭━╮╭━╮
┃┃┃╭╯╱╱┃┃╱╱╱╱╱╱╱┃┃╰╮┃┃╱╱╱╱┃┃╱╱╱╱╱╱╱╱╱╱╱┃┃╱╱┃┃╰╯┃┃
┃╰╯╯╭╮╭┫╰━┳━━╮╱╱┃╭╮╰╯┣━━┳━╯┣━━┳━━┳━━┳━━┫┃╱╱┃╭╮╭╮┣━━┳━╮╭━━┳━━┳━━┳━╮
┃╭╮┃┃┃┃┃╭╮┃┃━╋━━┫┃╰╮┃┃╭╮┃╭╮┃┃━┫╭╮┃╭╮┃╭╮┃┣━━┫┃┃┃┃┃╭╮┃╭╮┫╭╮┃╭╮┃┃━┫╭╯
┃┃┃╰┫╰╯┃╰╯┃┃━╋━━┫┃╱┃┃┃╰╯┃╰╯┃┃━┫╰╯┃╰╯┃╰╯┃╰┳━┫┃┃┃┃┃╭╮┃┃┃┃╭╮┃╰╯┃┃━┫┃
╰╯╰━┻━━┻━━┻━━╯╱╱╰╯╱╰━┻━━┻━━┻━━┫╭━┻━━┻━━┻━╯╱╰╯╰╯╰┻╯╰┻╯╰┻╯╰┻━╮┣━━┻╯
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭━╯┃
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰╯╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰━━╯
`))
}
