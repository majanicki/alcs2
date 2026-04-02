eps = 1e-6
def calulate_quality(init, optimum):
    return (init - optimum) / optimum * 100

def caluculate_efficiency(quality, time_measure):
    return quality / (time_measure + eps) * 100

tab = []
gready_sko_init = []
greedy_tai_init = []
steepest_sko_init = []
steepest_tai_init = []

gready_sko_end = []
greedy_tai_end = []
steepest_sko_end = []
steepest_tai_end = []
showfliers = False

data = {}
with open('./result/ex2.csv', 'r') as file:

    for line in file:
        tab.append(line.strip().split(','))

for i in range(1, len(tab)):
    type_alg, instance = tab[i][0].split('/')
    if instance not in data:
        data[instance] = {}
        if type_alg not in data[instance]:
            data[instance][type_alg] = []
    else:
        if type_alg not in data[instance]:
            data[instance][type_alg] = []
    
    data[instance][type_alg].append({
        'obj': float(tab[i][2]),
        'iterations': float(tab[i][3]),
        'evaluations': float(tab[i][4]),
        'time': float(tab[i][5]),
        'optimum': float(tab[i][8]),
    })

results = {}
for instance, algs in data.items():
    results[instance] = {}
    for alg, runs in algs.items():
        results[instance][alg] = []
        for run in runs:
            q = calulate_quality(run['obj'], run['optimum'])
            e = caluculate_efficiency(q, run['time'])
            results[instance][alg].append({
                'quality': q,
                'efficiency': e,
                'time': run['time'],
                'iterations': run['iterations'],
                'evaluations': run['evaluations'],
            })

import matplotlib.pyplot as plt
import numpy as np
import re

metrics = ['quality', 'efficiency', 'time', 'iterations', 'evaluations']
ylabels = {
    'quality': 'Quality (%)',
    'efficiency': 'Efficiency (%)',
    'time': 'Time (s)',
    'iterations': 'Steps',
    'evaluations': 'Evaluations'
}
step_algs = ["greedy", "steepest"]

def extract_number(s):
    match = re.search(r'\d+', s)
    return int(match.group()) if match else float('inf')
instances = sorted(results.keys(), key=extract_number)

for metric in metrics:
    plt.figure(figsize=(10,6))
    all_algs = set()
    for instance in instances:
        if metric == "iterations" or metric == "time":
            all_algs.update([alg for alg in results[instance].keys() if alg in step_algs])
        else:
            all_algs.update(results[instance].keys())
    all_algs = sorted(all_algs)
    
    for alg in all_algs:
        means = []
        stds = []
        mins = []
        maxs = []
        for instance in instances:
            if alg in results[instance]:
                values = [run[metric] for run in results[instance][alg]]
                if metric == "iterations" and alg == "heuristic":
                    print(values)
                means.append(np.mean(values))
                stds.append(np.std(values))
                mins.append(np.min(values))
                maxs.append(np.max(values))


        plt.errorbar(instances, means, yerr=[mins, maxs], label=alg, marker='o', capsize=5)

    plt.title(f'{ylabels[metric]} Across Instances')
    plt.xlabel('Instance')
    plt.ylabel(ylabels[metric])
    plt.xticks(rotation=45)
    plt.grid(True, linestyle='--', alpha=0.5)
    plt.legend()
    if metric == 'quality':
        plt.yscale('symlog', linthresh=1e1)
        plt.ylim(bottom=0)

    if metric == "time":
        plt.yscale('symlog',  linthresh=1e5)
        plt.ylim(bottom=0)

    if metric == "iterations":
        plt.yscale('log')
    if metric == "evaluations":
        plt.yscale('symlog',  linthresh=1e2)
        plt.ylim(bottom=0)

    plt.tight_layout()
    plt.savefig(f'./ex2/{metric}_aggregated.png')
    plt.close()