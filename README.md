# asteroid

# Program that send log formatted using a golang template to a logdna server
git clone https://github.com/GalacticHiker/asteroid.git

# to build
cd ~/git_repos/asteroid/src

go build -x -o ../runtime/bin/logdna-feeder cmd/logdna-feeder/main.go

# to execute
export LOGDNA_API_KEY=your_LOGDNA_API_KEY

../runtime/bin/logdna-feeder --hostname='logdna-feeder-host' --log-file-name='logdna-feeder-filename'


https://github.com/GalacticHiker/asteroid.git

# unimplemented cli arguments

tick  - the tick rate of the send clock, hardcoded to 500ms 
logsPerTick - how many logs to send per tick, hardcoded to 1; so send 2 log per second
nLogsToSend - how many log to send, hardcoded to 10

