eps = 1e-6
def calulate_quality(init, optimum):
    return (init - optimum) / optimum * 100

def caluculate_efficiency(quality, time_measure):
    return quality * (time_measure + eps) * 100

tab = []
gready_sko_init = []
greedy_tai_init = []
steepest_sko_init = []
steepest_tai_init = []

gready_sko_end = []
greedy_tai_end = []
steepest_sko_end = []
steepest_tai_end = []

data = {}
with open('.\\result\\ex2_t2.csv', 'r') as file:

    for line in file:
        tab.append(line.strip().split(','))

for i in range(len(tab)):
    # if tab[i][0] contains 'greedy'
    type_alg, instance = tab[i][0].split('/')
    # if type_alg in ['heuristic']:
    #     continue

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
        'time': int(tab[i][5]),
        'optimum': float(tab[i][8]),
    })

# tab_2 = []

# with open('.\\result\\ex2_t.csv', 'r') as file:

#     for line in file:
#         tab_2.append(line.strip().split(','))

# for i in range(len(tab_2)):
#     # if tab[i][0] contains 'greedy'
#     type_alg, instance = tab_2[i][0].split('/')

#     if instance not in data:
#         data[instance] = {}
#         if type_alg not in data[instance]:
#             data[instance][type_alg] = []
#     else:
#         if type_alg not in data[instance]:
#             data[instance][type_alg] = []
    
#     data[instance][type_alg].append({
#         'obj': float(tab_2[i][2]),
#         'iterations': float(tab_2[i][3]),
#         'evaluations': float(tab_2[i][4]),
#         'time': int(tab_2[i][5]),
#         'optimum': float(tab_2[i][8]),
#     })

print(data['tai150b'])



results = {}
for instance, algs in data.items():
    results[instance] = {}
    for alg, runs in algs.items():
        results[instance][alg] = []
        for run in runs:
            q = calulate_quality(run['obj'], run['optimum'])
            e = caluculate_efficiency(q, run['time'])
            # e = caluculate_efficiency(run['obj'], run['time'])
            results[instance][alg].append({
                'quality': q,
                'efficiency': e,
                'time': run['time'],
                'iterations': run['iterations'],
                'evaluations': run['evaluations'],
            })

import matplotlib.pyplot as plt
from statistics import mean, stdev

folder = '.\\ex2_2\\'

# # for each instanece plot box plot for quality
# for instance, algs in results.items():
#     fig, ax = plt.subplots()
#     data_to_plot = []
#     labels = []
#     for alg, runs in algs.items():
#         data_to_plot.append([run['quality'] for run in runs])
#         labels.append(alg)
#     ax.boxplot(data_to_plot, labels=labels)
#     ax.set_title(f'Quality for {instance}')
#     ax.set_ylabel('Quality (%)')
#     plt.savefig(f'{folder}quality_{instance}.png')

# for each instanece plot box plot for efficiency
for instance, algs in results.items():
    fig, ax = plt.subplots()
    means = []
    std_devs = []
    labels = []
    for alg, runs in algs.items():
        efficiencies = [run['efficiency'] for run in runs]
        means.append(mean(efficiencies) if efficiencies else 0.0)
        std_devs.append(stdev(efficiencies) if len(efficiencies) > 1 else 0.0)
        labels.append(alg)
    positions = range(len(labels))
    ax.bar(positions, means, yerr=std_devs, capsize=5)
    ax.set_xticks(list(positions))
    ax.set_xticklabels(labels)
    ax.set_title(f'Mean Efficiency for {instance}')
    ax.set_ylabel('Efficiency (%)')
    fig.tight_layout()
    plt.savefig(f'{folder}efficiency_{instance}.png')
    plt.close(fig)

# # for each instanece plot box plot for time
# for instance, algs in results.items():
#     fig, ax = plt.subplots()
#     data_to_plot = []
#     labels = []
#     for alg, runs in algs.items():
#         data_to_plot.append([run['time'] for run in runs])
#         labels.append(alg)
#     ax.boxplot(data_to_plot, labels=labels)
#     ax.set_title(f'Time for {instance}')
#     ax.set_ylabel('Time (s)')
#     plt.savefig(f'{folder}time_{instance}.png')

# # for each instanece plot box plot for iterations
# for instance, algs in results.items():
#     fig, ax = plt.subplots()
#     data_to_plot = []
#     labels = []
#     for alg, runs in algs.items():
#         data_to_plot.append([run['iterations'] for run in runs])
#         labels.append(alg)
#     ax.boxplot(data_to_plot, labels=labels)
#     ax.set_title(f'Iterations for {instance}')
#     ax.set_ylabel('Iterations')
#     plt.savefig(f'{folder}iterations_{instance}.png')


# # for each instanece plot box plot for evaluations
# for instance, algs in results.items():
#     fig, ax = plt.subplots()
#     data_to_plot = []
#     labels = []
#     for alg, runs in algs.items():
#         data_to_plot.append([run['evaluations'] for run in runs])
#         labels.append(alg)
#     ax.boxplot(data_to_plot, labels=labels)
#     ax.set_title(f'Evaluations for {instance}')
#     ax.set_ylabel('Evaluations')
#     plt.savefig(f'{folder}evaluations_{instance}.png')