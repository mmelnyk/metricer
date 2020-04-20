package metricer

type health struct {
	name    string
	help    string
	checker HealthcheckFunc
}

func (metric *health) Name() string {
	return metric.name
}

func (metric *health) Help() string {
	return metric.help
}

func (metric *health) Check() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errHealthCheckerPanic
		}
	}()

	if metric.checker != nil {
		err = metric.checker()
	}

	return
}
