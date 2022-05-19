package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

const headRef = "HEAD"

func GetCommit() (string, error) {
	stdout, _, _, err := execute(context.TODO(), "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return trim(stdout), nil
}

func GetBranch() (string, error) {
	stdout, _, _, err := execute(context.TODO(), "git", "rev-parse", "--abbrev-ref", "HEAD", "--quiet")
	if err != nil {
		return "", err
	}
	stdout = trim(stdout)

	if stdout == headRef {
		return "", errors.New("current Git Repository HEAD does not Point Branch-Ref")
	}

	return stdout, nil
}

func GetTag() (string, error) {
	stdout, _, _, err := execute(context.TODO(), "git", "describe", "--tags", "--abbrev=0", "--match=v*")
	if err != nil {
		return "", err
	}
	return trim(stdout), nil
}

func GetRepository() (string, error) {
	url, err := GetRemoteUrl()
	if err != nil {
		return "", err
	}
	return ParseRepositoryName(url)
}

func NormalizeBranch(branch string) string {
	if strings.Contains(branch, "refs/heads/") {
		return strings.TrimPrefix(branch, "refs/heads/")
	} else if strings.Contains(branch, "refs/pull/") {
		return strings.TrimPrefix(branch, "refs/pull/")
	} else {
		return branch
	}
}

func GetRemoteUrl() (string, error) {
	stdout, _, _, err := execute(context.TODO(), "git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	return trim(stdout), nil
}

func GetRepositoryFromEnv(env string) (string, error) {
	url, exists := os.LookupEnv(env)
	if !exists {
		return "", errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", env))
	}
	return ParseRepositoryName(url)
}

// ParseRepositoryName
// Input: https://github.krafton.com/example/exrepo.git
// Output: exrepo
func ParseRepositoryName(url string) (string, error) {
	args := strings.Split(url, "/")
	rawRepo := args[len(args)-1]
	if !strings.HasSuffix(rawRepo, ".git") {
		return "", fmt.Errorf("UnexpectedRepoUrl, Url: %s", url)
	}

	return strings.TrimSuffix(rawRepo, ".git"), nil
}

func trim(output string) string {
	return strings.ReplaceAll(strings.Split(output, "\n")[0], "'", "")
}

func execute(ctx context.Context, bin string, args ...string) (stdout string, stderr string, exit int, err error) {
	exit = 0

	cmd := exec.CommandContext(ctx, bin, args...)
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	cmdErr := cmd.Run()
	if cmdErr != nil {
		if exiterr, ok := cmdErr.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exit = status.ExitStatus()
			} else {
				err = fmt.Errorf("ExitError: %s", cmdErr.Error())
			}
		} else if exiterr2, ok := cmdErr.(*exec.Error); ok {
			err = fmt.Errorf("ExecError: PATH=%s, %s", os.Getenv("PATH"), exiterr2.Error())
		} else {
			err = fmt.Errorf("ExecUnknownError: %s", cmdErr.Error())
		}
	}

	stdout = bufOut.String()
	stderr = bufErr.String()
	zap.S().Debugf("Command: %s %v, Exitcode: %d", bin, args, exit)
	if stdout != "" {
		zap.S().Debugf("Stdout: %s", strings.TrimSuffix(stdout, "\n"))
	}
	if stderr != "" {
		zap.S().Debugf("Stderr: %s", strings.TrimSuffix(stderr, "\n"))
	}
	return
}
