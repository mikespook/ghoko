#/bin/bash

usage() {
    echo "Usage: $0 [-n] [-p] [-u] [-o]"
	printf "\t-n: SSH Host\n"
	printf "\t-p: SSH Port\n"
	printf "\t-u: SSH User\n"
	printf "\t-o: Output file\n"
    exit 1
}

while getopts 'p:n:u:o:' o &>> /dev/null; do
    case "$o" in
	p)
		PORT="$OPTARG";;
    n)
        HOST="$OPTARG";;
	u)
		USR="$OPTARG";;
	o)
		OUTPUT="$OPTARG";;    
    *)
        usage
    esac
done

APP=ghoko
base=`dirname -- "$0"`
if [ "$base" != "$0" ]; then
	cd $base
	base=./
fi

target=`mktemp -d`

# build binary
echo "[INFO] Building ..."
go build -o $target/$APP/usr/bin/$APP ../
if [ $? -ne 0 ] ; then
	exit 1
fi

# set version info
echo "[INFO] Generating version info ..."
release_date=`date +'%Y-%m-%d %T %z'`
ver=`date -d "$release_date" +'%s'`
echo "$ver"> $target/$APP/VERSION
echo "$release_date" >> $target/$APP/VERSION
git log --pretty=format:"[%ai] %h (%an) '%s'" --max-count=10 HEAD >> $target/$APP/VERSION
echo "" >> $target/$APP/VERSION

# copy
echo "[INFO] Copy files ..."
cp -R $base/usr $target/$APP
cp -R $base/etc $target/$APP
cp $base/install.sh $target/$APP
cp $base/uninstall.sh $target/$APP

echo "[INFO] Packing ..."
f=`mktemp -u`
pushd . > /dev/null
cd $target
tar zcvf $f $APP > /dev/null
popd > /dev/null

if [ "$OUTPUT" != "" ]; then
	echo "[INFO] Outputting ..."
	cp -f $f $OUTPUT
fi

if [ "$HOST" != "" ]; then
	[ "$PORT" == "" ] && PORT=22
	[ "$USR" == "" ] && USR=$USER
	echo "[INFO] Uploading ..."
	scp -P $PORT $f $USR@$HOST:$f
	read -p "Unpack? [y/n]" unpack
	if [ "$unpack" == "y" ] ; then
		ssh -p $PORT $USR@$HOST "tar zxvf $f && rm $f"
	else
		ssh -p $PORT $USR@$HOST "mv $f ~/$APP.tar.gz"
	fi
fi

rm -rf $target
rm -f $f

echo "[INFO] Complated!"
