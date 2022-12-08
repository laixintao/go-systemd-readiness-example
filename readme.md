This is an example of [systemd](https://systemd.io/) with
[`Type=Notify`](https://www.freedesktop.org/software/systemd/man/systemd.service.html).

Systemd will wait for your process until it said it was ready.

When use `Type=simple` like this:

```shell
# /lib/systemd/system/myservice.service
[Unit]
Description=Simple service

[Service]
Type=simple
ExecStart=/vagrant/go-systemd-readiness-example/main

[Install]
WantedBy=multi-user.target
```

When you run `systemctl start myservice`, systemctl command will exit
immediately. Logs:

```shell
root@vagrant:/vagrant/go-systemd-readiness-example# vim /lib/systemd/system/myservice.service
root@vagrant:/vagrant/go-systemd-readiness-example# systemctl daemon-reload
root@vagrant:/vagrant/go-systemd-readiness-example# systemctl start myservice
root@vagrant:/vagrant/go-systemd-readiness-example# journalctl -u myservice
-- Logs begin at Tue 2022-08-30 16:04:31 UTC, end at Thu 2022-12-08 09:45:02 UTC. --
Dec 08 09:44:55 vagrant systemd[1]: Started Simple service.
Dec 08 09:44:56 vagrant main[2543]: I am sleeping 0
Dec 08 09:44:57 vagrant main[2543]: I am sleeping 1
Dec 08 09:44:58 vagrant main[2543]: I am sleeping 2
Dec 08 09:44:59 vagrant main[2543]: I am sleeping 3
Dec 08 09:45:00 vagrant main[2543]: I am sleeping 4
Dec 08 09:45:00 vagrant main[2543]: I am up now!
```

So obviously systemd thinks your service is ready to work, but the truth is, it
need 5 more seconds to ready.

Let's change the systemd type to `Type=notify`, then restart the service:

```shell
root@vagrant:/vagrant/go-systemd-readiness-example# cat /lib/systemd/system/myservice.service
[Unit]
Description=Simple service

[Service]
Type=notify
ExecStart=/vagrant/go-systemd-readiness-example/main -systemd-type=notify

[Install]
WantedBy=multi-user.target

root@vagrant:/vagrant/go-systemd-readiness-example# systemctl daemon-reload
root@vagrant:/vagrant/go-systemd-readiness-example# systemctl restart myservice
```

You will find that `systemctl` blocks here for 5 seconds.

Check the logs:

```shell
Dec 08 09:48:09 vagrant systemd[1]: Stopped Simple service.
Dec 08 09:48:09 vagrant systemd[1]: Starting Simple service...
Dec 08 09:48:10 vagrant main[3288]: I am sleeping 0
Dec 08 09:48:11 vagrant main[3288]: I am sleeping 1
Dec 08 09:48:12 vagrant main[3288]: I am sleeping 2
Dec 08 09:48:13 vagrant main[3288]: I am sleeping 3
Dec 08 09:48:14 vagrant main[3288]: I am sleeping 4
Dec 08 09:48:14 vagrant main[3288]: I am up now!
Dec 08 09:48:14 vagrant systemd[1]: Started Simple service.
```

So now systemd didn't think your service is up, until 5 seconds later, when your
application said so: `daemon.SdNotify(false, daemon.SdNotifyReady)`

Code uses github.com/coreos/go-systemd/daemon
