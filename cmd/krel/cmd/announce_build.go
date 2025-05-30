/*
Copyright 2020 The Kubernetes Authors.

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

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/release-utils/command"
)

const (
	branchFlag            = "branch"
	changeLogFilePathFlag = "changelog-file-path"
	changeLogHTMLFlag     = "changelog-html-file"
	workDirFlag           = "workdir"
)

const semVerRegex string = `^?(\d+)(\.\d+)?(\.\d+)`

const branchCreationMsg = `Kubernetes Community,
<p>Kubernetes' {{ .Branch }} branch has been created.</p>
<p>The release owner will be sending updates on how to interact with this branch shortly.  The <a href=https://git.k8s.io/community/contributors/devel/sig-release/cherry-picks.md target="_blank">Cherrypick Guide</a> has some general guidance on how things will proceed.</p>
<p>Announced by your <a href=https://git.k8s.io/website/content/en/releases/release-managers.md target="_blank">Kubernetes Release Managers</a>.</p>
`

const releaseAnnouncementMsg = `Kubernetes Community,
<p>Kubernetes <b>{{ .Tag }}</b> has been built and pushed using Golang version <b>{{ .GoVersion }}</b> .</p>
<p>The release notes have been updated in <a href=https://git.k8s.io/kubernetes/{{ .ChangelogFilePath }}/#{{ .StrippedTag }} target="_blank">{{ .ChangelogFileName }}</a>, with a pointer to them on <a href=https://github.com/kubernetes/kubernetes/releases/tag/{{ .Tag }} target="_blank">github</a>:</p>
<p><hr>{{ .ChangelogHTML }}<hr></p>

<p><br>Contributors, the <a href=https://git.k8s.io/kubernetes/{{ .ChangelogFilePath }}/#{{ .StrippedTag }} target="_blank">{{ .ChangelogFileName }}</a> has been bootstrapped with {{ .Tag }} release notes and you may edit now as needed.</p>
<p><br><br>Published by your <a href=hhttps://git.k8s.io/website/content/en/releases/release-managers.md href target="_blank">Kubernetes Release Managers</a>.</p>
`

// buildAnnounceCmd represents the subcommand for `krel announce build`.
var buildAnnounceCmd = &cobra.Command{
	Use:           "build",
	Short:         "Build the announcement Kubernetes releases",
	Long:          "krel announce build",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var buildBranchAnnounceCmd = &cobra.Command{
	Use:           "branch",
	Short:         "Build the announcement Kubernetes for branch creation",
	Long:          "krel announce build branch",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuildBranchAnnounce(buildBranchAnnounceOpts, buildAnnounceOpts)
	},
}

var buildReleaseAnnounceCmd = &cobra.Command{
	Use:           "release",
	Short:         "Build the announcement Kubernetes for new releases",
	Long:          "krel announce build release",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuildReleaseAnnounce(buildReleaseAnnounceOpts, buildAnnounceOpts, announceOpts)
	},
}

type buildAnnounceOptions struct {
	workDir string
}

type buildBranchAnnounceOptions struct {
	branch string
}

type buildReleaseAnnounceOptions struct {
	changelogFilePath string
	changelogHTML     string
}

var (
	buildAnnounceOpts        = &buildAnnounceOptions{}
	buildBranchAnnounceOpts  = &buildBranchAnnounceOptions{}
	buildReleaseAnnounceOpts = &buildReleaseAnnounceOptions{}
)

func init() {
	buildBranchAnnounceCmd.PersistentFlags().StringVarP(
		&buildBranchAnnounceOpts.branch,
		branchFlag,
		"b",
		"",
		"set this flag when need to build the announcement for the branch creation, ie. release-1.19",
	)

	buildReleaseAnnounceCmd.PersistentFlags().StringVarP(
		&buildReleaseAnnounceOpts.changelogFilePath,
		changeLogFilePathFlag,
		"",
		"",
		"changelog path including the filename",
	)

	buildReleaseAnnounceCmd.PersistentFlags().StringVarP(
		&buildReleaseAnnounceOpts.changelogHTML,
		changeLogHTMLFlag,
		"",
		"",
		"contents of the changelog",
	)

	buildAnnounceCmd.PersistentFlags().StringVarP(
		&buildAnnounceOpts.workDir,
		workDirFlag,
		"",
		"",
		"working directory to store the announcement files",
	)

	if err := buildAnnounceCmd.MarkPersistentFlagRequired(workDirFlag); err != nil {
		logrus.Fatal(err)
	}

	// Check flags for krel announcement build branch command
	if err := buildBranchAnnounceCmd.MarkPersistentFlagRequired(branchFlag); err != nil {
		logrus.Fatal(err)
	}

	// Check flags for krel announcement build release command
	if err := buildReleaseAnnounceCmd.MarkPersistentFlagRequired(changeLogFilePathFlag); err != nil {
		logrus.Fatal(err)
	}

	if err := buildReleaseAnnounceCmd.MarkPersistentFlagRequired(changeLogHTMLFlag); err != nil {
		logrus.Fatal(err)
	}

	buildAnnounceCmd.AddCommand(buildBranchAnnounceCmd)
	buildAnnounceCmd.AddCommand(buildReleaseAnnounceCmd)
	announceCmd.AddCommand(buildAnnounceCmd)
}

// runBuildBranchAnnounce build the announcement file when creating the Kubernetes release branch.
func runBuildBranchAnnounce(opts *buildBranchAnnounceOptions, buildOpts *buildAnnounceOptions) error {
	logrus.Info("Building release announcement for branch creation")

	t, err := template.New("announcement-branch").Parse(branchCreationMsg)
	if err != nil {
		return err
	}

	announcement := bytes.Buffer{}
	if err := t.Execute(&announcement, struct {
		Branch string
	}{opts.branch}); err != nil {
		return fmt.Errorf("generating the announcement html file: %w", err)
	}

	return buildOpts.saveAnnouncement(announcement)
}

// runBuildReleaseAnnounce build the announcement file when creating a new Kubernetes release.
func runBuildReleaseAnnounce(opts *buildReleaseAnnounceOptions, buildOpts *buildAnnounceOptions, announceOpts *announceOptions) error {
	if err := announceOpts.Validate(); err != nil {
		return fmt.Errorf("validating announcement send options: %w", err)
	}

	logrus.Info("Building release announcement for new release")

	t, err := template.New("announcement-release").Parse(releaseAnnouncementMsg)
	if err != nil {
		return err
	}

	goVersion, err := getGoVersion()
	if err != nil {
		return err
	}

	changelogHTML, err := os.ReadFile(opts.changelogHTML)
	if err != nil {
		return err
	}

	announcement := bytes.Buffer{}
	if err := t.Execute(&announcement, struct {
		Tag               string
		StrippedTag       string
		GoVersion         string
		ChangelogFilePath string
		ChangelogFileName string
		ChangelogHTML     string
	}{
		announceOpts.tag,
		strings.ReplaceAll(announceOpts.tag, ".", ""),
		goVersion,
		opts.changelogFilePath,
		filepath.Base(opts.changelogFilePath),
		string(changelogHTML),
	}); err != nil {
		return fmt.Errorf("generating the announcement html file: %w", err)
	}

	return buildOpts.saveAnnouncement(announcement)
}

func (opts *buildAnnounceOptions) saveAnnouncement(announcement bytes.Buffer) error {
	logrus.Info("Creating announcement files")

	absOutputPath := filepath.Join(opts.workDir, "announcement.html")
	logrus.Infof("Writing HTML file to %s", absOutputPath)

	err := os.WriteFile(absOutputPath, announcement.Bytes(), os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("saving announcement.html: %w", err)
	}

	logrus.Info("Kubernetes Announcement created.")

	return nil
}

func getGoVersion() (string, error) {
	cmdStatus, err := command.New(
		"go", "version").
		RunSilentSuccessOutput()
	if err != nil {
		return "", fmt.Errorf("get go version: %w", err)
	}

	versionRegex := regexp.MustCompile(semVerRegex)

	return versionRegex.FindString(strings.TrimSpace(cmdStatus.Output())), nil
}
