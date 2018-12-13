THREADS=(1 10 100 1000)

declare -a TIMES
for t in "${THREADS[@]}"
do
  echo $(go run ctr_parallel.go --input inputs/$1 --workers ${t})
done


