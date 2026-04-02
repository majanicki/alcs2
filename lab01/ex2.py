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

metrics_aggregated = {}

for metric in metrics:
    metrics_aggregated[metric] = {}
    # plt.figure(figsize=(10,6))
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
        metrics_aggregated[metric][alg] = {
            "means": means,
            "stds": stds,
            "maxs": maxs,
            "mins": mins,
        }
        # plt.errorbar(instances, means, yerr=[mins, maxs], label=alg, marker='o', capsize=5)


# plot quality
quality_metrics = metrics_aggregated["quality"]


fig, (outlier_ax, normal_ax) = plt.subplots(2, 1, sharex=True)
fig.subplots_adjust(hspace=0.05)  # adjust space between Axes

x = np.arange(len(instances))
# jitter strength
jitter = 0.1
# create jittered x positions

for alg in quality_metrics:
# plot the same data on both Axes
    x_jittered = x + np.random.uniform(-jitter, jitter, size=len(x))
    outlier_ax.errorbar(x_jittered ,quality_metrics[alg]["means"], yerr=quality_metrics[alg]["stds"], label=alg, marker='o', capsize=5)
    normal_ax.errorbar(x_jittered, quality_metrics[alg]["means"],  yerr=quality_metrics[alg]["stds"], label=alg, marker='o', capsize=5)


bottom_outlier_range = 200
top_normal_range = 125
# zoom-in / limit the view to different portions of the data
outlier_ax.set_ylim(bottom_outlier_range, 3400)  # outliers only
normal_ax.set_ylim(0, top_normal_range)  # most of the data
outlier_ax.hlines([bottom_outlier_range], xmin=[0], xmax=[len(instances)], color="k", linewidth=4)
normal_ax.hlines ([top_normal_range], xmin=[0],     xmax=[len(instances)], color="k", linewidth=4)

# hide the spines between ax and normal_ax
outlier_ax.spines.bottom.set_visible(False)
normal_ax.spines.top.set_visible(False)
outlier_ax.xaxis.tick_top()
outlier_ax.tick_params(labeltop=False)  # don't put tick labels at the top
# outlier_ax.set_xticks(range(0,120,1))
normal_ax.xaxis.tick_bottom()

# Now, let's turn towards the cut-out slanted lines.
# We create line objects in axes coordinates, in which (0,0), (0,1),
# (1,0), and (1,1) are the four corners of the Axes.
# The slanted lines themselves are markers at those locations, such that the
# lines keep their angle and position, independent of the Axes size or scale
# Finally, we need to disable clipping.

d = .85  # proportion of vertical to horizontal extent of the slanted line
kwargs = dict(marker=[(-1, -d), (1, d)], markersize=12,
              linestyle="none", color='k', mec='k', mew=1, clip_on=False)
outlier_ax.plot([0, 1], [0, 0], transform=outlier_ax.transAxes, **kwargs)
normal_ax.plot([0, 1], [1, 1], transform=normal_ax.transAxes, **kwargs)

plt.show()



    # plt.title(f'{ylabels[metric]} Across Instances')
    # plt.xlabel('Instance')
    # plt.ylabel(ylabels[metric])
    # plt.xticks(rotation=45)
    # plt.grid(True, linestyle='--', alpha=0.5)
    # plt.legend()
    # if metric == 'quality':
    #     plt.yscale('symlog', linthresh=1e1)
    #     plt.ylim(bottom=0)

    # if metric == "time":
    #     plt.yscale('symlog',  linthresh=1e5)
    #     plt.ylim(bottom=0)

    # if metric == "iterations":
    #     plt.yscale('log')
    # if metric == "evaluations":
    #     plt.yscale('symlog',  linthresh=1e2)
    #     plt.ylim(bottom=0)

    # plt.tight_layout()
    # plt.savefig(f'./ex2/{metric}_aggregated.png')
    # plt.close()