#jobd

Jobd implements *cron* as a 9p file system. 9p is [a simple file protocol to build sane distributed systems](http://9p.cat-v.org/). It's the cornerstone of [the Styx Architecture for Distributed Systems](http://doc.cat-v.org/inferno/4th_edition/styx) which represents the system's resources as a form of file system.

##Usage
```
jobd --help
Usage of jobd:
  -alsologtostderr=false: log to standard error as well as files
  -dbdir="/var/lib/jobd": Location of the jobd jobs database
  -debug=false: 9p debugging to stderr
  -fsaddr="0.0.0.0:5640": Address where job file service listens for connections
  -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
  -log_dir="": If non-empty, write log files in this directory
  -logtostderr=false: log to standard error instead of files
  -stderrthreshold=0: logs at or above this threshold go to stderr
  -v=0: log level for V logs
  -vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```

##Design

*cron* is a time-based job scheduler. Jobd represents jobs as subdirectories of  a *jobs* directory. Each *job* subdirectory contains four files:

* the **ctl** file which is used to start and stop the job
* the **cmd** file that records the command the job executes
* the **log** file that is used to retrieve the job's execution history
* the **schedule** file that records the job's schedule and its next scheduled execution time

To start a job, write the string **start** to the *ctl* file. Similarly to stop it write the string **stop** to the *ctl* file. Read from the *cmd*, *log*, or *schedule* file to retrieve the information they provide.

Jobd jobs are created via the *clone* file. The *clone* file is a peer of the *jobs* directory in the jobd name space.
