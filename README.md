# Structured logging with func names


## Usage

```golang

package main

import (
	"errors"

	"github.com/sirupsen/logrus"
	logerr "github.com/timurguseynov/logrus-logerr"
)

func main() {
	logFields := logrus.Fields{"func": "main"}
	err := outer()
	logerr.Entry(err, logFields).Error(err)
}

func outer() error {
	logFields := logrus.Fields{"func": "outer"}
	return logerr.WithFields(inner(), logFields)
}

func inner() error {
	logFields := logrus.Fields{"func": "inner"}
	return logerr.WithFields(errors.New("something went wrong"), logFields)
}

```

## Output

```
ERRO[0000] something went wrong                          func="inner,outer,main"
```
