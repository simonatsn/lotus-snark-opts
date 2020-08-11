This is a self contained benchmark for snark proof verification.

The tests directory contains serialized proofs from the Lotus testnet. The test and benchmark simply deserialize and verify those proofs. 

# Build Lotus

env RUSTFLAGS="-C target-cpu=native -g" FFI_BUILD_FROM_SOURCE=1 make clean deps bench

# Test

Run tests:
go run main.go

By default it will validate every prf file in tests. 

# Bench

For stability we set up for running benchmarks by turning off turbo. For AMD:
'''
echo "0" | sudo tee /sys/devices/system/cpu/cpufreq/boost
for i in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
do
  echo "performance" | sudo tee $i
done
'''

For Intel:
'''
echo "1" | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo
echo "100" | sudo tee /sys/devices/system/cpu/intel_pstate/min_perf_pct
'''

Run all benchmarks:
go test -bench=.

Run WindowPoST benchmark (10 full sized window proofs):
go test -bench=VerifyWindow -benchtime 20x

