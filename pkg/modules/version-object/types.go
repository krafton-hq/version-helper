package version_object

import "encoding/json"

type VersionObj struct {
	ApiVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Metadata   *Metadata `json:"metadata"`
	Spec       *Spec     `json:"spec"`
	Status     *Status   `json:"status"`
}

type Metadata struct {
	Name        string            `json:"name"`
	Project     string            `json:"project"`
	BaseVersion string            `json:"baseVersion"`
	Revision    uint              `json:"revision"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type Spec struct {
	GitRef *GitRef `json:"gitRef"`
}

type GitRef struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Branch     string `json:"branch"`
}

type Status struct {
	Artifact []*Artifact `json:"artifact"`
}

type Artifact struct {
	Platform     string `json:"platform"`
	Target       string `json:"target,omitempty"`
	ArtifactType string `json:"artifactType"`
	Uri          string `json:"uri"`
	Description  string `json:"description,omitempty"`
}

func Merge(obj1 *VersionObj, obj2 *VersionObj) (*VersionObj, error) {
	m1, err := toMap(obj1)
	if err != nil {
		return nil, err
	}
	m2, err := toMap(obj2)
	if err != nil {
		return nil, err
	}

	m3 := mergeMaps(m1, m2)
	return fromMap(m3)
}

func toMap(obj *VersionObj) (map[string]any, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	m := map[string]any{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func fromMap(m map[string]any) (*VersionObj, error) {
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	obj := &VersionObj{}
	err = json.Unmarshal(buf, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for key, value := range a {
		out[key] = value
	}
	for key, rawValue := range b {
		if value, ok := rawValue.(map[string]interface{}); ok {
			if mergedRawV, ok := out[key]; ok {
				if mergedV, ok := mergedRawV.(map[string]interface{}); ok {
					out[key] = mergeMaps(mergedV, value)
					continue
				}
			}
		}

		if value, ok := rawValue.([]interface{}); ok {
			if mergedRawV, ok := out[key]; ok {
				if mergedV, ok := mergedRawV.([]interface{}); ok {
					out[key] = append(mergedV, value...)
					continue
				}
			}
		}
		out[key] = rawValue
	}
	return out
}

func EqualKey(obj1 *VersionObj, obj2 *VersionObj) bool {
	return obj1.Kind == obj2.Kind &&
		obj1.Metadata.Name == obj2.Metadata.Name &&
		obj1.Metadata.Project == obj2.Metadata.Project
}
