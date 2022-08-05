package jsonrpc2

import "fmt"

type calls map[string]any

func (c calls) Add(name string, method any) error {
	if method == nil {
		return fmt.Errorf("failed to add nil method")
	}

	if _, err := c.Get(name); err != nil {
		c[name] = method
	} else {
		return fmt.Errorf("call [%s] already registered", name)
	}
	return nil
}

func (c calls) Get(name string) (any, error) {
	if method, ok := c[name]; ok {
		return method, nil
	} else {
		return nil, fmt.Errorf("call [%s] isn't registered", name)
	}
}

func (c *calls) Remove(name string) error {
	if _, err := c.Get(name); err != nil {
		return fmt.Errorf("%s", err)
	}
	delete(*c, name)
	return nil
}
