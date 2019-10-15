# Site-Health
Use to check the health of a website. Checks the response status and response's duration.

## Compile and Run
```
go build health.go
./health.exe <uri> [-c <count> | -r | -g]
```
* ```-c <count>``` runs the query ```<count>``` times
* ```-r``` repeats the query indefinitely
* ```-g``` repeats the query indefinitely with a graphic display
