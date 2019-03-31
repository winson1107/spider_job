#!/bin/bash
set -e
#初始化
os_name=`uname -s`
machine=`uname -m`
PWD_PATH=`pwd`
#chromium_verison 版本号
chromium_verison='599821'
local_chromium_dir=$PWD_PATH"/node_modules/puppeteer/.local-chromium"
if [[ "$os_name" == "Linux" ]]; then
    download_url="https://npm.taobao.org/mirrors/chromium-browser-snapshots/Linux_x64/${chromium_verison}/chrome-linux.zip"
    #依赖库
	yum install pango.x86_64 libXcomposite.x86_64 libXcursor.x86_64 libXdamage.x86_64 libXext.x86_64 libXi.x86_64 libXtst.x86_64 cups-libs.x86_64 libXScrnSaver.x86_64 libXrandr.x86_64 GConf2.x86_64 alsa-lib.x86_64 atk.x86_64 gtk3.x86_64 wget -y
	#字体
	yum install ipa-gothic-fonts xorg-x11-fonts-100dpi xorg-x11-fonts-75dpi xorg-x11-utils xorg-x11-fonts-cyrillic xorg-x11-fonts-Type1 xorg-x11-fonts-misc -y
	local_chromium_dir=${local_chromium_dir}"/linux-"${chromium_verison}"/"
elif [[ "$os_name" == "Darwin" ]]; then
    download_url="https://npm.taobao.org/mirrors/chromium-browser-snapshots/Mac/${chromium_verison}/chrome-mac.zip"
    local_chromium_dir=${local_chromium_dir}"/mac-"${chromium_verison}"/"
fi

read -p "是否需要自动下载Chromium (1,自动下载,2,手动下载)：" INPUT_STRING
case $INPUT_STRING in
	1)
		npm install --registry https://registry.npm.taobao.org
		;;
	2)
		echo "正在从 npm.taobao.org 下载 Chromium……"
		wget -O chrome.zip ${download_url}
		npm install  --ignore-scripts --registry https://registry.npm.taobao.org
		if [[ ! -d ${local_chromium_dir} ]]; then
            mkdir -p ${local_chromium_dir}
        fi
		unzip -d ${local_chromium_dir} chrome.zip
		rm chrome.zip
		;;
	*)
		echo "选择不正确"
		;;
esac
