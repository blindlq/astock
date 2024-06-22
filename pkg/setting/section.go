package setting

var sections = make(map[string]interface{})


type ScriptsSettingS struct {
	Url      string
	Token    string
	DailyHotUrl  string
}

// ReadSection 读取配置
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

// ReloadAllSection 重载配置
func (s *Setting) ReloadAllSection() error  {
	for k, v := range sections {
		err := s.ReadSection(k,v)
		if err != nil {
			return err
		}
	}

	return nil
}