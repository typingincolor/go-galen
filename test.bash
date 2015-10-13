#!/bin/bash

ERROR=0
echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d );
do
        if ls $Dir/*.go &> /dev/null;
        then
            godep go test -coverprofile=profile.out $Dir
            if [ "$?" -ne "0" ]
            then
                ERROR=1
            fi
            
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "mode: set" >> acc.out
            fi
fi
done
goveralls -coverprofile=acc.out -service=travis-ci $COVERALLS
rm -rf ./profile.out
rm -rf ./acc.out

exit $ERROR
