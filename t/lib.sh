ok() {
    status="not ok"
    if [ $3 -eq 0 ]; then
       status="ok"
    fi
    echo $status $1 - $2
}

not_ok() {
    status="not ok"
    if [ $3 -eq 1 ]; then
       status="ok"
    fi
    echo $status $1 - $2
}
