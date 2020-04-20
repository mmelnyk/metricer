package metricer

const (
	logname = "metric"
)

const (
	rtmetricgoroutines     = "_goroutines"
	rtmetricgoroutineshelp = "Number of goroutines running by app [internal]"
	rtmetricos             = "_os"
	rtmetricoshelp         = "Application platform [internal]"
	rtmetricnumcpu         = "_num_cpu"
	rtmetricnumcpuhelp     = "Number of CPU [internal]"
	rtmetricmemalloc       = "_mem_alloc"
	rtmetricmemallochelp   = "Number of allocated memory for whole app in bytes [internal]"
	healthcheckfailed      = "_failed_healthchecks"
	healthcheckfailedhelp  = "Number of failed health checks [internal]"
	rtuptime               = "uptime"
	rtuptimehelp           = "Application uptime in nanosec [internal]"
)

const (
	acceptJSON      = "application/json"
	acceptText      = "text/plain"
	charsetUTF8     = "charset=utf-8"
	contenttypeJSON = acceptJSON + "; " + charsetUTF8
	contenttypeText = acceptText + "; " + charsetUTF8
	defaultPort     = 9110
)

const (
//envLocalhostOnly = "METRICER_LOCALHOST_ONLY"
//envPort          = "METRICER_PORT"
//envDebug         = "METRICER_DEBUG"
)
