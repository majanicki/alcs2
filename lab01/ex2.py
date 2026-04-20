import csv
import numpy as np
import re
import matplotlib.pyplot as plt
# Map column index → target dictionary
columns = {
    3: {},  # fitnesses
    4: {},  # iterations
    5: {},  # evaluations
    6: {},  # times
    7: {},  # initial fitness
}

instances = set()
with open('measurements_5.csv', newline='') as csvfile:
    reader = csv.reader(csvfile)

    for row in reader:
        instance = row[0]
        if instance == "esc32f":
            continue
        alg = row[1]
        instances.add(instance)
        for col_idx, target_dict in columns.items():
            value = int(row[col_idx])
            if alg not in target_dict:
                # if alg == "heuristic":
                #     continue
                target_dict[alg] = {}
            if instance not in target_dict[alg]:
                target_dict[alg][instance] = []

            target_dict[alg][instance].append(value)

# Optional: unpack for readability
fitnesses = columns[3]
iterations = columns[4]
evaluations = columns[5]
times = columns[6]
initial_fitness = columns[7]


instances = list(instances)
instances.sort(key=lambda x: int(re.search(r'\d+', x).group()))
optima = {}
for instance in instances:
    with open(f"./data/{instance}.sln", "r") as file:
        line = file.readline()
        optimum = int(line.split()[1])
        optima[instance] = optimum


x = np.arange(len(instances))

jitter_strength = 0.1 
image_size = (8,5)
plt.figure(figsize=image_size)


for i, alg in enumerate(times):
    y = []
    stds = []
    for instance in instances:
        mean = np.mean(times[alg][instance]) / 1e6
        std = np.std(times[alg][instance]) / 1e6
        y.append(mean)
        stds.append(std)
    y = np.array(y)
    stds = np.array(stds)

    jitter = (i - len(times)/2) * jitter_strength
    x_jittered = x + jitter

    plt.errorbar(x_jittered, y, yerr=stds, label=alg, marker=".")

plt.xticks(x, instances)
plt.yscale("log")
plt.legend()
plt.title("Execution Time ↓")
plt.ylabel("Time [ms]")
plt.savefig("ex2/time_aggregate_4a.png")
plt.figure(figsize=image_size)

for i, alg in enumerate(times):
    y = []
    stds = []
    for instance in instances:
        improvment = (np.array(fitnesses[alg][instance]) - optima[instance]) / optima[instance]
        f = improvment
        # print(f)
        mean = np.mean(f)
        # print(alg, mean)
        std = np.std(f)
        y.append(mean)
        stds.append(std)
    y = np.array(y)
    stds = np.array(stds)

    jitter = (i - len(times)/2) * jitter_strength
    x_jittered = x + jitter

    plt.errorbar(x_jittered, y, yerr=stds, label=alg, marker=".")

plt.xticks(x, instances)
# plt.yscale("log")
plt.legend()
plt.title("Quality ↓")
plt.ylabel("Ratio to optimum [q]")
plt.savefig("ex2/quality_aggregate_4a.png")
plt.show()
plt.figure(figsize=image_size)

for i, alg in enumerate(times):
    y = []
    stds = []
    for instance in instances:
        improvment = (np.array(initial_fitness[alg][instance]) - np.array(fitnesses[alg][instance])) / optima[instance]
        f = improvment / (np.array(times[alg][instance]) / 1e6)
        # print(f)
        mean = np.mean(f)
        # print(alg, mean)
        std = np.std(f)
        y.append(mean)
        stds.append(std)
    y = np.array(y)
    stds = np.array(stds)

    jitter = (i - len(times)/2) * jitter_strength
    x_jittered = x + jitter

    plt.errorbar(x_jittered, y, yerr=stds, label=alg, marker=".")

plt.xticks(x, instances)
plt.yscale("symlog",  linthresh=0.0001)
plt.legend()
plt.title("Efficiency ↑")
plt.ylabel("Improvment to quality from initial solution over time [q/ms]")
plt.savefig("ex2/efficiency_aggregate_4a.png")
plt.show()

plt.figure(figsize=image_size)

for i, alg in enumerate(iterations):
    if alg not in ["steepest", "greedy"]:
        continue
    y = []
    stds = []
    for instance in instances:
        mean = np.mean(iterations[alg][instance])
        std = np.std(iterations[alg][instance])
        y.append(mean)
        stds.append(std)
    y = np.array(y)
    stds = np.array(stds)

    jitter = (i - len(iterations)/2) * jitter_strength
    x_jittered = x + jitter

    plt.errorbar(x_jittered, y, yerr=stds, label=alg, marker=".")

plt.xticks(x, instances)
plt.yscale("log")
plt.title("Iterations")
plt.legend()
plt.savefig("ex2/iterations_aggregate_4a.png")
plt.show()

plt.figure(figsize=image_size)

for i, alg in enumerate(evaluations):
    y = []
    stds = []
    for instance in instances:
        mean = np.mean(evaluations[alg][instance])
        std = np.std(evaluations[alg][instance])
        y.append(mean)
        stds.append(std)
    y = np.array(y)
    stds = np.array(stds)

    jitter = (i - len(evaluations)/2) * jitter_strength
    x_jittered = x + jitter

    plt.errorbar(x_jittered, y, yerr=stds, label=alg, marker=".")

plt.xticks(x, instances)
plt.yscale("log")
plt.title("Evaluations")
plt.legend()
plt.savefig("ex2/evaluations_aggregate_4a.png")
plt.show()
