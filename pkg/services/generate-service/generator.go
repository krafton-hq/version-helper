package generate_service

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.krafton.com/sbx/version-helper/assets"
	metadata_resolver "github.krafton.com/sbx/version-helper/pkg/modules/metadata-resolver"
	path_utils "github.krafton.com/sbx/version-helper/pkg/modules/path-utils"
	"github.krafton.com/sbx/version-helper/pkg/modules/templates"
	"github.krafton.com/sbx/version-helper/pkg/modules/versions"
	"go.uber.org/zap"
)

type Service struct {
	TemplateFilePath string
	TemplateFileType TemplateFileType
	GenerateDir      string
	GenerateFilePath string
}

func NewService(templateFileUrl string, generateDir string, generateFilePath string) (*Service, error) {
	u, err := url.Parse(templateFileUrl)
	if err != nil {
		return nil, err
	}

	if u.Path == "" {
		return nil, fmt.Errorf("InvalidARgumentError, File Path is Null, url: %#v", u)
	}
	var fileType TemplateFileType
	var filePath string
	switch u.Scheme {
	case "embed":
		fileType = Embed
		filePath = strings.TrimPrefix(u.Path, "/")
		break
	case "file":
	case "":
		fileType = File
		filePath = u.Path
		break
	default:
		return nil, fmt.Errorf("InvalidArgumentError, templateFile scheme is not valid %s", u.Scheme)
	}

	if generateDir == "" {
		generateDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	s := &Service{
		TemplateFilePath: filePath,
		TemplateFileType: fileType,
		GenerateDir:      generateDir,
		GenerateFilePath: generateFilePath,
	}
	zap.S().Debugf("Initialize Generate Service, %+v", s)
	return s, nil
}

type TemplateFileType string

const (
	Embed TemplateFileType = "embed"
	File  TemplateFileType = "file"
)

func (s *Service) GenerateAndSave(version versions.Version, metadata *metadata_resolver.BuildMetadata, project string) error {
	tmpl, err := s.getTemplate()
	if err != nil {
		return err
	}

	values := &templates.TemplateRoot{
		Version:     version.String(),
		Project:     project,
		BaseVersion: version.BaseVersion(),
		Revision:    version.Revision(),
		FileName:    s.TemplateFilePath,
		Major:       version.Major(),
		Minor:       version.Minor(),
		Patch:       version.Patch(),
		Git: &templates.TemplateGit{
			Repository: metadata.Repository,
			Commit:     metadata.CommitSha,
			Branch:     metadata.Branch,
		},
	}

	genData, err := templates.Template(tmpl, values)
	if err != nil {
		return err
	}

	genPath, err := path_utils.ResolvePathToAbs(filepath.Join(s.GenerateDir, s.GenerateFilePath))
	if err != nil {
		zap.S().Debugf("Resolve Absolute Path Failed, dir: %s, file: %s, error: %s", s.GenerateDir, s.GenerateFilePath, err.Error())
		return err
	}

	err = ioutil.WriteFile(genPath, []byte(genData), 0644)
	if err != nil {
		zap.S().Infof("Write File Failed, error: %s", err.Error())
		return err
	}
	zap.S().Debugf("Write Generated File to %s", genPath)
	return nil
}

func (s *Service) getTemplate() (string, error) {
	switch s.TemplateFileType {
	case Embed:
		tmpl, err := assets.GetFile(s.TemplateFilePath)
		if err != nil {
			zap.S().Infof("File Not Found from embedded storage, name: %s error: %s", s.TemplateFilePath, err.Error())
			return "", err
		}
		return tmpl, nil

	case File:
		buf, err := ioutil.ReadFile(s.TemplateFilePath)
		if err != nil {
			wd, _ := os.Getwd()
			zap.S().Infof("File Not Found from volume, name: %s, workdir: %s, error: %s", s.TemplateFilePath, wd, err.Error())
			return "", err
		}
		return string(buf), nil

	default:
		zap.S().Errorf("InvalidEnumError: Template File Type is not valid %s", s.TemplateFileType)
		return "", fmt.Errorf("InvalidEnumError: Template File Type is not valid %s", s.TemplateFileType)
	}
}
