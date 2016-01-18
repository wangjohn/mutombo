from fabric.api import *

def deploy():
    with cd("~/golang/src/github.com/wangjohn/mutombo"):
        with shell_env(GOPATH="/home/zincadm/golang", PATH="$PATH:/usr/local/go/bin"):
            run("git pull origin master")
            run("/usr/local/go/bin/go get")
            run("/usr/local/go/bin/go install")

    sudo("gem install einhorn")
    run("einhornsh -e upgrade")
