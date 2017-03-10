package parseUtil

func Int32ToPointer(input int32) *int32 {
	tmp := new(int32)
	*tmp = input
	return tmp
}

func BoolToPointer(input bool) *bool {
	tmp := new(bool)
	*tmp = input
	return tmp
}
