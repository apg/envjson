. ./t/lib.sh

echo "1..2"

./envjson -i t/required-foo.json 2> /dev/null
not_ok 1 'Required FOO fails with no value' $?

FOO=bar ./envjson t/required-foo.json 2>&1 > /dev/null
ok 2 'Required FOO succeeds when FOO is set and in environment' $?
