# asteroid

A program that sends logs formatted using a golang template to a rsyslogd:tcp:514, logdna server

git clone https://github.com/GalacticHiker/asteroid.git

# Build (on mac workstation)
cd ~/git_repos/asteroid/src

go build -x -o ../runtime/bin/logmill cmd/logmill/main.go


# Build for docker (all libraries statically linked)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../runtime/bin/logmill-linux cmd/logdna-feeder/main.go

# Build on my Centos7 VM
/usr/local/go/bin/go build -x -o ../runtime/bin/logmill cmd/logmill/main.go/home/tharnett/asteroid/src

# Design Intents
should always be a stand-alone(NO! required support files) executable

../runtime/ is the directory footprint that is assumed when using optional file artifacts (e.g log template)

if a file path starts with a '/' it is assumed to be an absolute path

# To send log to logdna (via logdna api)
export LOGDNA_API_KEY=your_LOGDNA_API_KEY

../runtime/bin/logmill --hostname='logdna-feeder-host' --log-file-name='logdna-feeder-filename' --logsPerTick=1 --nLogsToSend=10 --tick=1s

../runtime/bin/logmill --hostname='logdna-feeder-host' --log-file-name='logdna-feeder-filename' --logsPerTick=1 --nLogsToSend=10 --tick=1s --template=defaultKVP

https://github.com/GalacticHiker/asteroid.git

# unimplemented cli arguments

tick  - the tick rate of the send clock, hardcoded to 500ms 
logsPerTick - how many logs to send per tick, hardcoded to 1; so send 2 log per second
nLogsToSend - how many log to send, hardcoded to 10

