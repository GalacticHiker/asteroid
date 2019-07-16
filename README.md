# asteroid

# Program that send log formatted usin a golang template

git clone https://github.com/GalacticHiker/asteroid.git

# to build
cd asteroid/logdna-kitchen/src

go build -x -o ../runtime/bin/logdna-feeder cmd/logdna-feeder/main.go

# should not be exposed
export LOGDNA_API_KEY=your_LOGDNA_API_KEY

../runtime/bin/logdna-feeder --hostname='logdna-feeder-host' --log-file-name='logdna-feeder-filename'


https://github.com/GalacticHiker/asteroid.git
