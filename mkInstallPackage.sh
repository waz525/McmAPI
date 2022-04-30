#!/bin/sh
WORK_DIR=$(cd $(dirname $0) && pwd )
cd $WORK_DIR

echo "===> make McmAPI ..."

InstallTempPath="${WORK_DIR}/../Package/McmAPI_install"
[ -d ${InstallTempPath} ] && rm -rf ${InstallTempPath}

mkdir -p ${InstallTempPath}/McmAPI/bin/
mkdir -p ${InstallTempPath}/McmAPI/log/
cp -r conf ${InstallTempPath}/McmAPI/
cp -r www ${InstallTempPath}/McmAPI/
rm -fv ${InstallTempPath}/McmAPI/www/html/UploadFile/*
cp release-note ${InstallTempPath}/McmAPI/
cp install.sh ${InstallTempPath}/
chmod +x ${InstallTempPath}/install.sh
cp conf/McmAPI.conf ${InstallTempPath}/
Version=`cat release-note | grep "^Version" | head -n 1 | cut -f 2 -d '='`
echo "Version: $Version"

Author="Worden"
BuildTime=`date +'%Y-%m-%d %H:%M:%S'`
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'BinInfo.Version=${Version}' \
    -X 'BinInfo.Author=${Author}' \
    -X 'BinInfo.BuildTime=${BuildTime}' \
    -X 'BinInfo.BuildGoVersion=${BuildGoVersion}' \
"


cd bin/
rm -f McmAPI
go build -ldflags "$LDFlags" McmAPI.go DBProcess.go FilterProcess.go
mv McmAPI ${InstallTempPath}/McmAPI/bin/
cp uninstall.sh ${InstallTempPath}/McmAPI/bin/
chmod +x ${InstallTempPath}/McmAPI/bin/*.sh

cd ${WORK_DIR}/../Package/
rm -f McmAPI_install-${Version}-`date +"%Y%m%d"`.tar.gz
tar czf McmAPI_install-${Version}-`date +"%Y%m%d"`.tar.gz McmAPI_install
echo "=================> cd ${WORK_DIR}/../Package/ ;  ls -lh McmAPI_install-* :"
ls -lh McmAPI_install-*




