import matplotlib.pyplot as plt


# ex4
# # open fiel .\result\ext4.csv and read it
# tab = []
# gready_sko = []
# greedy_tai = []
# steepest_sko = []
# steepest_tai = []
# with open('.\\result\\ex4.csv', 'r') as file:

#     for line in file:
#         tab.append(line.strip().split(','))

# for i in range(1, len(tab)):
#     # if tab[i][0] contains 'greedy'
#     if 'greedy' in tab[i][0]:
#         if 'sko' in tab[i][0]:
#             gready_sko.append(float(tab[i][2]))
#         elif 'tai' in tab[i][0]:
#             greedy_tai.append(float(tab[i][2]))
#     elif 'steepest' in tab[i][0]:
#         if 'sko' in tab[i][0]:
#             steepest_sko.append(float(tab[i][2]))
#         elif 'tai' in tab[i][0]:
#             steepest_tai.append(float(tab[i][2]))


# greedy_sko_best = []
# greedy_tai_best = []
# steepest_sko_best = []
# steepest_tai_best = []
# for i in range(len(gready_sko)):
#     greedy_sko_best.append(min(gready_sko[:i + 1]))
#     greedy_tai_best.append(min(greedy_tai[:i + 1]))
#     steepest_sko_best.append(min(steepest_sko[:i + 1]))
#     steepest_tai_best.append(min(steepest_tai[:i + 1]))

# greedy_sko_avg = []
# greedy_tai_avg = []
# steepest_sko_avg = []
# steepest_tai_avg = []
# for i in range(len(gready_sko)):
#     greedy_sko_avg.append(sum(gready_sko[:i + 1]) / (i + 1))
#     greedy_tai_avg.append(sum(greedy_tai[:i + 1]) / (i + 1))
#     steepest_sko_avg.append(sum(steepest_sko[:i + 1]) / (i + 1))
#     steepest_tai_avg.append(sum(steepest_tai[:i + 1]) / (i + 1))



# # plot greedy for sko
# plt.plot(greedy_sko_best, label='best found so far')
# plt.plot(greedy_sko_avg, label='average found so far')
# plt.legend()
# plt.title('Greedy Sko90')
# plt.xlabel('multi-random start')
# plt.ylabel('Objective Value')
# plt.show()

# # plot greedy for tai
# plt.plot(greedy_tai_best, label='best found so far')
# plt.plot(greedy_tai_avg, label='average found so far')
# plt.legend()
# plt.title('Greedy Tai256c')
# plt.xlabel('multi-random start')
# plt.ylabel('Objective Value')
# plt.show()

# # plot steepest for sko
# plt.plot(steepest_sko_best, label='best found so far')
# plt.plot(steepest_sko_avg, label='average found so far')
# plt.legend()
# plt.title('Steepest Sko90')
# plt.xlabel('multi-random start')
# plt.ylabel('Objective Value')
# plt.show()

# # plot steepest for tai
# plt.plot(steepest_tai_best, label='best found so far')
# plt.plot(steepest_tai_avg, label='average found so far')
# plt.legend()
# plt.title('Steepest Tai256c')
# plt.xlabel('multi-random start')
# plt.ylabel('Objective Value')
# plt.show()


# ex 5
tab = []
def similarity(list1, list2):
    count = 0
    for i in range(len(list1)):
        if list1[i] == list2[i]:
            count += 1
    return count / len(list1)


def parse_permutation(value):
    value = value.strip()
    if value.startswith('[') and value.endswith(']'):
        value = value[1:-1].strip()
    if not value:
        return []
    return [int(x) for x in value.replace(',', ' ').split()]

greedy_sko_permutations = []
greedy_tai_permutations = []
steepest_sko_permutations = []
steepest_tai_permutations = []

greedy_sko_optimal = []
greedy_tai_optimal = []
steepest_sko_optimal = []
steepest_tai_optimal = []

greedy_sko_quality = []
greedy_tai_quality = []
steepest_sko_quality = []
steepest_tai_quality = []

with open('.\\result\\ex4.csv', 'r') as file:

    for line in file:
        tab.append(line.strip().split(','))

for i in range(1, len(tab)):
    # if tab[i][0] contains 'greedy'
    if 'greedy' in tab[i][0]:
        if 'sko' in tab[i][0]:
            greedy_sko_permutations.append(parse_permutation(tab[i][6]))
            greedy_sko_optimal.append(parse_permutation(tab[i][7]))
            greedy_sko_quality.append(float(tab[i][2]))
        elif 'tai' in tab[i][0]:
            greedy_tai_permutations.append(parse_permutation(tab[i][6]))
            greedy_tai_optimal.append(parse_permutation(tab[i][7]))
            greedy_tai_quality.append(float(tab[i][2]))
    elif 'steepest' in tab[i][0]:
        if 'sko' in tab[i][0]:
            steepest_sko_permutations.append(parse_permutation(tab[i][6]))
            steepest_sko_optimal.append(parse_permutation(tab[i][7]))
            steepest_sko_quality.append(float(tab[i][2]))
        elif 'tai' in tab[i][0]:
            steepest_tai_permutations.append(parse_permutation(tab[i][6]))
            steepest_tai_optimal.append(parse_permutation(tab[i][7]))
            steepest_tai_quality.append(float(tab[i][2]))



sim_greedy_sko = []
sim_greedy_tai = []
sim_steepest_sko = []
sim_steepest_tai = []

for i in range(len(greedy_sko_permutations)):
    sim_greedy_sko.append([])
    sim_greedy_tai.append([])
    sim_steepest_sko.append([])
    sim_steepest_tai.append([])
    for j in range(len(greedy_sko_permutations)):
        if i != j:
            sim_greedy_sko[i].append(similarity(greedy_sko_permutations[i], greedy_sko_permutations[j]))
            sim_greedy_tai[i].append(similarity(greedy_tai_permutations[i], greedy_tai_permutations[j]))
            sim_steepest_sko[i].append(similarity(steepest_sko_permutations[i], steepest_sko_permutations[j]))
            sim_steepest_tai[i].append(similarity(steepest_tai_permutations[i], steepest_tai_permutations[j]))

sim_greedy_sko_avg = []
sim_greedy_tai_avg = []
sim_steepest_sko_avg = []
sim_steepest_tai_avg = []

for i in range(len(greedy_sko_permutations)):
    sim_greedy_sko_avg.append(sum(sim_greedy_sko[i]) / len(sim_greedy_sko[i]))
    sim_greedy_tai_avg.append(sum(sim_greedy_tai[i]) / len(sim_greedy_tai[i]))
    sim_steepest_sko_avg.append(sum(sim_steepest_sko[i]) / len(sim_steepest_sko[i]))
    sim_steepest_tai_avg.append(sum(sim_steepest_tai[i]) / len(sim_steepest_tai[i]))



sim_greedy_sko_optimal = []
sim_greedy_tai_optimal = []
sim_steepest_sko_optimal = []
sim_steepest_tai_optimal = []

for i in range(len(greedy_sko_optimal)):
    sim_greedy_sko_optimal.append(similarity(greedy_sko_permutations[i], greedy_sko_optimal[i]))
    sim_greedy_tai_optimal.append(similarity(greedy_tai_permutations[i], greedy_tai_optimal[i]))
    sim_steepest_sko_optimal.append(similarity(steepest_sko_permutations[i], steepest_sko_optimal[i]))
    sim_steepest_tai_optimal.append(similarity(steepest_tai_permutations[i], steepest_tai_optimal[i]))


#plot quality vs similarity for greedy sko
plt.scatter(greedy_sko_quality, sim_greedy_sko_avg, label='Average similarity to other solutions')
plt.scatter(greedy_sko_quality, sim_greedy_sko_optimal, label='Similarity to the optimal solution')
plt.title('Greedy Sko90')
plt.xlabel('Objective Value')
plt.ylabel('Similarity')
plt.legend()
plt.show()


#plot quality vs similarity for greedy tai
plt.scatter(greedy_tai_quality, sim_greedy_tai_avg, label='Average similarity to other solutions')
plt.scatter(greedy_tai_quality, sim_greedy_tai_optimal, label='Similarity to the optimal solution')
plt.title('Greedy Tai256c')
plt.xlabel('Objective Value')
plt.ylabel('Similarity')
plt.legend()
plt.show()

#plot quality vs similarity for steepest sko
plt.scatter(steepest_sko_quality, sim_steepest_sko_avg, label='Average similarity to other solutions')
plt.scatter(steepest_sko_quality, sim_steepest_sko_optimal, label='Similarity to the optimal solution')
plt.title('Steepest Sko90')
plt.xlabel('Objective Value')
plt.ylabel('Similarity')
plt.legend()
plt.show()

#plot quality vs similarity for steepest tai
plt.scatter(steepest_tai_quality, sim_steepest_tai_avg, label='Average similarity to other solutions')
plt.scatter(steepest_tai_quality, sim_steepest_tai_optimal, label='Similarity to the optimal solution')
plt.title('Steepest Tai256c')
plt.xlabel('Objective Value')
plt.ylabel('Similarity')
plt.legend()
plt.show()






