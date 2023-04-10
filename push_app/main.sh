#!/bin/bash

#Variables
declare -i choice=0
declare -i aws_client_cnt=0
declare -i chameleon_client_cnt=0
declare -i azure_client_cnt=0
declare -i gcp_client_cnt=0
declare -i client_cnt_total=0

while [ "$choice" != 5 ]
do
    echo "Please enter the cloud provider: "
    echo "1: AWS"
    echo "2: Chameleon"
    echo "3: Azure (not supported)"
    echo "4: GCP (not supported)"
    echo "5: exit"
    read choice 
    case $choice in
        1)
            echo "selected AWS" ;;
        2) 
            echo "selected Chameleon" ;;
        3) 
            echo "selected Azure" ;;
        4)
            echo "selected GCP" ;;
        5)
            echo "exit" ;;
        *)
            echo "wrong input"
    esac
    if [ $choice -eq 5 ]
    then
        break
    fi
    echo "Enter the number of clients to spawn, MAKE SURE YOU HAVE ENOUGH RESOURCE QUOTA"
    if [ $choice -eq 1 ]
    then
        read aws_client_cnt
        for((i=0;i<$aws_client_cnt;i++));  
        do     
            bash scripts/rosa_cluster_depl.sh $i $client_cnt_total
            ((client_cnt_total=client_cnt_total+1))
            echo "client id: $client_cnt_total"
        done  
    elif [ $choice -eq 2 ]
    then
        read chameleon_client_cnt
        echo "not supported"
    elif [ $choice -eq 3 ]
    then
        read azure_client_cnt
        echo "not supported"
    elif [ $choice -eq 4 ]
    then
        read gcp_client_cnt
        echo "not supported"
    else
        echo "wrong input"
    fi
    echo ""
done
