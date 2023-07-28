#!/bin/bash

chmod a+x init-controller.sh init-middle-platform.sh init-ui.sh
# 启动脚本1
./init-controller.sh &

# 启动脚本2
./init-middle-platform.sh &

# 启动脚本3
./init-ui.sh &

read -rn1