#!/bin/bash

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE

module=start
app=broker
conf=broker-a.toml
sourceDir=../../conf/
targetDir=conf
logsDir=logs
pidfile=$logsDir/broker.pid
logfile=$logsDir/broker.log

mkdir -p $logsDir $targetDir

function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running already, pid="
        cat $pidfile
        return 1
    fi

    nohup ./$app -c $conf >> $logfile 2>&1 &
    echo $! > $pidfile
    echo "$app started..., pid=$!"
}

function stop() {
    pid=`cat $pidfile`
    kill -9 $pid
    echo "$app stoped..."
}

function restart() {
    stop
    sleep 1
    start
}

function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running, pid="
        cat $pidfile
    else
        echo "$app is stoped"
    fi
}

function tailf() {
    tail -f $logfile
}

function build() {
    go build
    if [ $? -ne 0 ]; then
        exit $?
    fi
    mv $module $app
    ./$app -v
}

function pack() {
    build
    git log -1 --pretty=%h > gitversion
    version=`./$app -v`
    cp -rf $sourceDir .
    packName=$app-$version.tar.gz
    rm -f $packName
    tar -zcvf $packName control.sh $targetDir $app gitversion
    rm -rf $targetDir
}

function packbin() {
    build
    git log -1 --pretty=%h > gitversion
    version=`./$app -v`
    tar -zcvf $app-bin-$version.tar.gz $app gitversion
}

function help() {
    echo "$0 build|pack|packbin|start|stop|restart|status|tail"
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "status" ];then
    status
elif [ "$1" == "tail" ];then
    tailf
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "pack" ];then
    pack
elif [ "$1" == "packbin" ];then
    packbin
else
    help
fi
