#!/bin/bash
ssh-keygen -R [120.53.123.43]:36040
ssh -o StrictHostKeyChecking=no -p 36040 testuser@120.53.123.43 run-katago