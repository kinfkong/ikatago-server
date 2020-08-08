#!/bin/bash
ssh-keygen -R [192.168.1.222]:2222
ssh -o StrictHostKeyChecking=no -p 2222 kinfkong@192.168.1.222 run-katago --name katago-1.5.0 --config 1po --weight 20b