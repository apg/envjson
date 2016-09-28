. ./t/lib.sh

echo "1..1"

FOO=bar echo '{"FOO": "baz"}' | ./envjson --stdin t/required-foo.json env | grep 'FOO=baz' > /dev/null
ok 1 'STDIN overwrites process environment, and FOO=baz appears in output of env command.' $?
