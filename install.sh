#!/bin/sh
INSTALL_PATH="/home/McmAPI"
PRENT_PATH="/home/"

WORK_DIR=$(cd $(dirname $0) && pwd )
cd $WORK_DIR

if [ -d ${INSTALL_PATH} ]
then
	version="`cat ${INSTALL_PATH}/release-note| grep "^Version" | head -n 1 | cut -f 2 -d '='`"
	[ "$version" == "" ] && version="x.x.x"
	echo "===> Backup old version:$version ..."
	mv -v ${INSTALL_PATH} ${INSTALL_PATH}_${version}_`date +%s`
	mkdir -p ${PRENT_PATH}
fi

mkdir -p ${PRENT_PATH}

cd $WORK_DIR
echo "===> Install McmAPI ..."
cp -v -r McmAPI ${INSTALL_PATH}
scp -v McmAPI.conf ${INSTALL_PATH}/conf/

#echo "Modify Configure ..."

if [ `rpm -aq | grep supervisor | wc -l ` -eq 0 ]
then
	echo "===> Install supervisor ..."
	yum -y install supervisor
	systemctl start supervisord.service
	systemctl enable supervisord.service
fi

echo "===> Configure supervisor ..."
scp -v ${INSTALL_PATH}/conf/McmAPI.ini /etc/supervisord.d/
supervisorctl update
supervisorctl restart McmAPI
echo "===> McmAPI status :"
supervisorctl status



