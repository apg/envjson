. ./t/lib.sh

echo "1..2"

FOO=bar ./envjson t/required-foo.json env | grep 'FOO=bar' > /dev/null
ok 1 'env runs successfully and FOO=bar appears in output.' $?

./envjson t/required-foo.json env 2> /dev/null
not_ok 2 "required FOO not present, so no command run" $?
