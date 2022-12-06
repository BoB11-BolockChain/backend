from io import TextIOWrapper
from typing import List
import argparse

parser = argparse.ArgumentParser(description='docker module')

class CustomAction(argparse.Action):
    def __call__(self, parser, namespace, values, option_string=None):
        if not 'order' in namespace:
            setattr(namespace, 'order', [])
        # previous = namespace.order
        namespace.order.append(self.dest)
        # setattr(namespace, 'order', previous)
        
        if getattr(namespace, self.dest) == None:
            setattr(namespace, self.dest, [])
        getattr(namespace, self.dest).append(values)

parser.add_argument("-f", "--FROM", dest="os", action="store")
parser.add_argument("-r", "--RUN", dest="run", nargs='*', action=CustomAction)
parser.add_argument("-e", "--ENV", dest="env", nargs='*', action=CustomAction)
parser.add_argument("-x", "--EXPOSE", dest="expose", nargs='*', action=CustomAction)
parser.add_argument("-w", "--WORKDIR", dest="workdir", nargs='*', action=CustomAction)
parser.add_argument("-p", "--COPY", dest="copy", nargs='*', action=CustomAction)
parser.add_argument("-c", "--CMD", dest="cmd", nargs='*', action=CustomAction)
parser.add_argument("-a", "--ADD", dest="add", nargs='*', action=CustomAction)

DF_CMD = {
    'os' : "FROM ",
    'run' : 'RUN ',
    'copy' : 'COPY ',
    'add' : 'ADD ',
    'env' : 'ENV ',
    'expose' : 'EXPOSE ',
    'cmd' : 'CMD '
}

def write_line(cmd:str, value:str|List[str], file:TextIOWrapper):
    file.write(DF_CMD[cmd] + value)
    
class Insts:
    def os(val):
        print(val)
        # return DF_CMD["os"] + val
    def run(val):
        print(val)
        # return DF_CMD["run"] + val
    def copy(val):
        print(val)
        # return DF_CMD["copy"] + val
    def add(val):
        print(val)
        # return DF_CMD["add"] + val
    def env(val):
        print(val)
        # return DF_CMD["env"] + val
    def expose(val):
        print(val)
        # return DF_CMD["expose"] + val
    def cmd(val):
        print(val)
        # return DF_CMD["cmd"] + val

# def create_dockerfile(order:List[str], values:List[List[str]]):
#     with open("Dockerfile", "w") as df:
#         for cmd in order:
#             pass


if (__name__ == "__main__"):
    args = parser.parse_args()
    print(args.os)
    # print(args.order)
    for inst in args.order:
        for val in getattr(args, inst):
            getattr(Insts, inst)(val)
    
    
    
# def default_docker_file(filename):
#     with open(filename, 'w') as dockerfile_f:
#         dockerfile_f.write(Web_Docker.FROM())
#         dockerfile_f.write(Web_Docker.DEFAULT_RUN())
#         dockerfile_f.write(Web_Docker.RUN("apt update && apt instll -y net-tools"))
# if __name__ == '__main__':
#     default_docker_file("test_dockerfile")