# uptime-koguma
remote agent for uptime-kuma

Currently in beta. This is a very simple remote push agent to work with Uptime Kuma. An example use case is to monitor linux boxes running behind a firewall. Let's say, you have your Uptime Kuma running in a cloud vps but you want to monitor your homelab without opening ports on your firewall. The push agent will periodically make an http GET request to Uptime Kuma reporting either status=up or status=down. At the moment, it is able to check cpu load average, free memory and available disk space. If cpu load average exceeds the configured limit or if the available free memory drops below the limit or if the avilable space on the file system is below the configured limit, then status=down will be reported. Otherwise if all three are within limits status=up will be sent.
### How to build:
1) Clone the repo
2) cd into uptime-koguma folder
3) Build:
```
$ go build -o koguma
```
### How to configure:
1) Obtain the Push URL from your running Uptime Kuma instance
2) Edit the config file configs/koguma.conf:
```
{
  "push_url": "https://cyclops.ddns.net/api/push/YDCd3KsEnz?status=up&msg=OK&ping=",
  "heartbeat_interval": 60,
  "cpu_threshold": 70,
  "cpu_load_average_type": 15,
  "memory_threshold": 10,
  "disks": [
    {
      "disk_path": "/",
      "threshold": 20
    },
    {
      "disk_path": "/home/",
      "threshold": 20
    }
  ]
}
```
Config parameters:
- push_url: String. The push URL obtained from your running Uptime Kuma instance
- heartbeat_interval: Time duration. It needs to match whatever is configured in Uptime Kuma
- cpu_threshold: Uint. A percentage, without the % sign. If the load average reaches this threshold it will triger an alarm. That is, koguma will simply report status=down.
- cpu_load_average_type: Uint. Possible values are: 1, 5 and 15. The well-known load average intervals from top command.
- memory_threshold: Uint. Percentage, without the % sign. If the available memory is below this threshold, it will triger status=down
- disks: List of 0 or more 'disk objects'
  - disk_path: String. The mount point of the disk you want to monitor.
  - threshold: Uint. Percentage, without the % sign. If the available disk space is below this, it will trigger status=down
  
### How to run:
1) On a box with systemd, copy configs/koguma.service file into /etc/systemd/system folder and modify the path to the koguma executeable and the path to koguma.conf files.
```
[mik@mik uptime-koguma] (master)$ cat configs/koguma.service
[Unit]
Description=Uptime Koguma
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=mik
ExecStart=/home/mik/src/uptime-koguma/koguma -f /home/mik/src/uptime-koguma/configs/koguma.conf

[Install]
WantedBy=multi-user.target
```
2) Enable koguma service
```
$ systemctl enable koguma.service
```
3) Start koguma service
```
$ systemctl start koguma.service
```


