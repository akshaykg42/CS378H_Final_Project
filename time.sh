
# The following function is taken from https://stackoverflow.com/questions/1527049/join-elements-of-an-array
function join_by {
  local IFS="$1"; shift; echo "$*";
}

THREADS=(1 10 100 1000)

declare -a TIMES
for t in "${THREADS[@]}"
do
  echo $(go run ctr_parallel.go --input inputs/$1 --workers ${t})

done


