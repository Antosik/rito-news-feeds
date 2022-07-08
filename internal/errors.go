package internal

import "fmt"

type ErrorCollector []error

func (c *ErrorCollector) Collect(e error) { *c = append(*c, e) }

func (c *ErrorCollector) CollectMany(e []error) { *c = append(*c, e...) }

func (c *ErrorCollector) CollectFrom(e ErrorCollector) { *c = append(*c, e...) }

func (c *ErrorCollector) Error() (err string) {
	err = "Collected errors:\n"
	for i, e := range *c {
		err += fmt.Sprintf("\tError %d: %s\n", i, e.Error())
	}

	return err
}

func (c *ErrorCollector) Size() int {
	return len(*c)
}

func NewErrorCollector() *ErrorCollector {
	return new(ErrorCollector)
}
