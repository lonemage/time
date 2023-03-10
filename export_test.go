package time

type Value = dValue

func MakeValue(monotonic,
	now int64,
	location string) *Value {
	return &Value{monotonic, now, location}
}

func GetValue(v *Value) (monotonic,
	now int64,
	location string) {
	return v.monotonic, v.now, v.location
}
