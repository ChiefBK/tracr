package properties

type Props map[string]interface{}

func (self Props) GetString(key string) *string {
	value, hasKey := self[key]

	if !hasKey {
		// log error
		return nil
	}

	stringValue, isString := value.(string)

	if !isString {
		// log error
		return nil
	}

	return &stringValue
}

func (self Props) GetInt(key string) *int {
	value, hasKey := self[key]

	if !hasKey {
		// log error
		return nil
	}

	intValue, isInt := value.(int)

	if !isInt {
		// log error
		return nil
	}

	return &intValue
}