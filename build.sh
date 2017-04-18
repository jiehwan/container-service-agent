echo "****************************"
if [ "$1" = "arm" ]; then
        echo "Target Binary arch is ARM"
        export GOARCH=arm GOARM=7
        export CC="arm-linux-gnueabi-gcc"
else
        echo "Target Binary arch is amd64"
        export GOARCH=amd64
        export CC="gcc"
fi

echo make clean
make clean

echo make build
make build

