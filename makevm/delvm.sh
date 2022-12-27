#!/bin/sh
#virsh qemu 자동 삭제 쉘스크립트 - H4uN

#삭제할 VM 이름
virsh list --all
echo -n "Input del VM name : "
read DEL_VM_NAME

# VM안에 생성된 .qcow2를 삭제
DeleteDIR=./VM/${DEL_VM_NAME}
if [  -d $DeleteDIR ]; then
    rm -rf ${DeleteDIR}

    #vm 삭제 virsh 명령어
    virsh destroy ${DEL_VM_NAME}
    virsh undefine ${DEL_VM_NAME}
    echo "Deleting..."
    sleep 5
    virsh list --all
fi