package shellPlugin

type outputGetter func(command string) (string, string, error)
