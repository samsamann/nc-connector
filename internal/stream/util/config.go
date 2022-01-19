package util

import "errors"

var ErrConfigNotFound = errors.New("config not found")

type ConfigMap struct {
	config map[string]interface{}
	errs   []error
}

func NewConfigMap(config map[string]interface{}) *ConfigMap {
	return &ConfigMap{
		config: config,
		errs:   make([]error, 0),
	}
}

func (m *ConfigMap) Get(name string) ConfigValidator {
	return &validationMetadata{configMap: m, name: name}
}

func (m *ConfigMap) Error() error {
	if len(m.errs) > 0 {
		// TODO add error
		return nil
	}
	return nil
}

type ConfigValidator interface {
	Required() ConfigValidator
	String() string
	StringWithDefault(string) string
	Map() map[string]string
	Slice() []string
}

type validationMetadata struct {
	configMap *ConfigMap
	name      string
	required  bool
}

func (meta *validationMetadata) Required() ConfigValidator {
	meta.required = true
	return meta
}

func (meta *validationMetadata) String() string {
	val, hasErr := meta.getVal()
	if hasErr || val == nil {
		return ""
	}

	if str, ok := val.(string); ok {
		return str
	}
	meta.configMap.errs =
		append(meta.configMap.errs, errors.New("value must be a string"))
	return ""
}

func (meta *validationMetadata) StringWithDefault(str string) string {
	meta.required = false
	if ret := meta.String(); len(ret) != 0 {
		return ret
	}
	return str
}

func (meta *validationMetadata) Map() map[string]string {
	val, hasErr := meta.getVal()
	if hasErr || val == nil {
		return make(map[string]string)
	}

	var ok bool
	var infMap map[interface{}]interface{}
	if infMap, ok = val.(map[interface{}]interface{}); !ok {
		meta.configMap.errs = append(
			meta.configMap.errs,
			errors.New("value must be a list of key value pairs"),
		)
		return make(map[string]string)
	}
	m := make(map[string]string)
	for k, v := range infMap {
		m[k.(string)] = v.(string)
	}
	return m
}

func (meta *validationMetadata) Slice() []string {
	val, hasErr := meta.getVal()
	if hasErr || val == nil {
		return make([]string, 0)
	}

	var ok bool
	var infSlice []interface{}
	if infSlice, ok = val.([]interface{}); !ok {
		meta.configMap.errs = append(
			meta.configMap.errs,
			errors.New("value must be a list"),
		)
		return make([]string, 0)
	}
	s := make([]string, 0)
	for _, v := range infSlice {
		s = append(s, v.(string))
	}
	return s
}

func (meta *validationMetadata) getVal() (interface{}, bool) {
	if val, ok := meta.configMap.config[meta.name]; ok {
		if val == nil && meta.required {
			meta.configMap.errs =
				append(meta.configMap.errs, errors.New("config is required"))
			return nil, true
		}
		return val, false
	}
	if meta.required {
		meta.configMap.errs = append(meta.configMap.errs, ErrConfigNotFound)
		return nil, true
	}
	return nil, false
}
