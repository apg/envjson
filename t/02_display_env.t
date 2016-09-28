. ./t/lib.sh

echo "1..2"

FOO=bar ./envjson | grep '"FOO":"bar"' > /dev/null
ok 1 'Display environment in JSON when no command present' $?

./envjson -i | grep '{}' > /dev/null
ok 2 'Display ignored environment is empty' $?