echo "Trying to release port 8080."
echo "Allow? (y/n)"

read symbol

if [ $symbol = 'y' ] || [ $symbol = 'Y' ] 
then
    pid="$(sudo lsof -i :8080 | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        sudo kill -9 $pid
        pid="$(sudo lsof -i :8080 | awk '{print $2}' | head -2 | tail -1)"
    done
fi