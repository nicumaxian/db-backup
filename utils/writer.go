package utils

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }
