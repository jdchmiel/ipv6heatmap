# ipv6heatmap
This is my playground for learning Golang, you probably do not want this.

To run this you need to have docker installed on you destination. 
`docker-compose up` will start the code.
Note this consumes port 8080 so that needs to be available.  If it is not available please change 
```
ports: 
  - "8080:8080"
```
to 
```
ports: 
  - "<your_port>:8080"
```

you can reach the website via the ip of the machine you installed this at, specifying the port if you did not choose 80.
I ran this under docker-machine which had an ip address of 192.168.99.100 so the url was 
`http://192.168.99.100:8080`

If you chose port 80 and were running this on your localmachine you would use
`http://127.0.0.1`

You are free to copy / distribute / do whatever you want with this, no support or guarantee of any nature is implied.


Known bugs:
- the json encoding is terribly inefficient
- wrong scope for heat var, re adds the same points when the map is only shifted a tiny bit, needs to clear map before adding points.
- no design / styling / css at all yet
- no data tweaking, still using the raw db dumb, with only the Latitude column indexed.
- - no zoom level support
- hard coded 10k limit to the points returned vs a more designed approach as to which points to return.


