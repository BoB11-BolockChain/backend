#!/bin/sh
#virsh qemu 자동 생성 쉘스크립트 - H4uN

# vm 이름
echo -n "Input VM name : "
read VM_NAME

# 지정한 경로에 폴더가 존재하지 않으면 폴더를 생성
CreateDIR=./VM/${VM_NAME}
if [ ! -d $CreateDIR ]; then
    mkdir ${CreateDIR}
    NewVM=${CreateDIR}/${VM_NAME}.qcow2
    cp windows_10_x64_comp.qcow2 ${NewVM}

    #vm 생성 virsh 명령어
    virt-install --name=${VM_NAME} --ram=4096 --cpu=host --vcpus=1 --os-type=windows --os-variant=win10 --disk path=${NewVM},device=disk,bus=virtio,format=qcow2 --network network=default,model=virtio --graphics vnc,password=test,listen=0.0.0.0 --import --wait 0 --check all=off
    virsh list --all
fi