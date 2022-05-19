package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"go.uber.org/zap"
)

type TemplateRoot struct {
	Version *TemplateVersion
	Git     *TemplateGit
	Project string
}

type TemplateVersion struct {
	FullVersion string
	BaseVersion string
	PreRelease  string

	Major    uint
	Minor    uint
	Patch    uint
	Branch   string
	Revision uint
	Commit   string
}

type TemplateGit struct {
	Repository string
	Commit     string
	Branch     string
}

func Template(tmpl string, templateName string, values *TemplateRoot) (string, error) {
	zap.S().Debug("Start Template with Sprig Function Map")
	zap.S().Debugf("template: %s", tmpl)
	zap.S().Debugw("", "values", values)

	t := template.New(templateName).Funcs(sprig.TxtFuncMap())
	t, err := t.Parse(tmpl)
	if err != nil {
		zap.S().Debugf("Create Go Template Failed, error: %s", err.Error())
		return "", fmt.Errorf("CreateGoTemplateFailed, error: %s", err.Error())
	}
	t.Option("missingkey=error")

	var buf bytes.Buffer
	err = t.Execute(&buf, values)
	if err != nil {
		zap.S().Debugf("Templating Failed, error: %s", err.Error())
		return "", fmt.Errorf("TemplatingFailed, error: %s", err.Error())
	}

	return buf.String(), nil
}
