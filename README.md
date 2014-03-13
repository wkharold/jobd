#jobd

Jobd implements *cron* as a 9p file system. 9p is [a simple file protocol to build sane distributed systems](http://9p.cat-v.org/). It's the cornerstone of [the Styx Architecture for Distributed Systems](http://doc.cat-v.org/inferno/4th_edition/styx) which represents the system's resources as a form of file system. Jobd illustrates this approach to system design using the [go9p](https://code.google.com/p/go9p) library.

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

Once jobd is started the file system it provides can be mounted via
```
$ mount -t 9p -o protocol=tcp,port=5640 <addr> <mountpoint>
```
Where **addr** is the IP address of the box running jobd. Note that it needn't be mounted on the machine running jobd, any box running a Linux 3.x kernel should suffice. 

##Design

*cron* is a time-based job scheduler, it has two primary concerns: *jobs* which are commands to be executed, and *schedules* that determine when a job is run. The design of a 9p-based application or system service generally begins with the creation of a *name space*, think file system subtree, that represents the application's resources in terms of files and directories. 

Jobd represents jobs as subdirectories of  a *jobs* directory. Each *job* subdirectory contains four files:

* the **ctl** file which is used to start and stop the job
* the **cmd** file that records the command the job executes
* the **log** file that is used to retrieve the job's execution history
* the **schedule** file that records the job's schedule and its next scheduled execution time

To start a job, write the string **start** to the *ctl* file
```
$ echo -n start > <mountpoint>/jobs/<job>/ctl
```
Similarly to stop it write the string **stop** to the *ctl* file
```
$ echo -n stop > <mountpoint>/jobs/<job>/ctl
```
Read from the *cmd*, *log*, or *schedule* file to retrieve the information they provide
```
$ cat <mountpoint>/jobs/<job>/cmd
echo hello world
$ cat <mountpoint>/jobs/<job>/schedule
0 0/5 * * * ? *
$ cat <mountpoint>/jobs/<job>/log
2014-02-11 09:42:33.454707331 -0600 CST:started
2014-02-11 09:42:35.004655691 -0600 CST:hello world
2014-02-11 09:42:40.003579265 -0600 CST:hello world
2014-02-11 09:42:45.003220637 -0600 CST:hello world
2014-02-11 09:42:50.00294003 -0600 CST:hello world
...
```

Jobd jobs are created via the *clone* file. The *clone* file is a peer of the *jobs* directory in the jobd name space. To create a job write a string of the form: <jobname>:<cronexpr>:<cmd> to the clone file
```
$ echo -n 'hello:0 0/5 * * * ? *:echo hello world' > <mountpoint>/clone
```

##TODO

* support deleting jobs
