import matplotlib.pyplot as plt

def rank_correlation(x_vals, y_vals):
    if len(x_vals) != len(y_vals) or len(x_vals) < 2:
        return float('nan')

    n = len(x_vals)
    rank_x = sorted(range(n), key=lambda i: x_vals[i])
    rank_y = sorted(range(n), key=lambda i: y_vals[i])

    d_squared_sum = sum((rank_x[i] - rank_y[i]) ** 2 for i in range(n))
    return 1 - (6 * d_squared_sum) / (n * (n**2 - 1))


def pearson_corr(x_vals, y_vals):
    if len(x_vals) != len(y_vals) or len(x_vals) < 2:
        return float('nan')

    x_mean = sum(x_vals) / len(x_vals)
    y_mean = sum(y_vals) / len(y_vals)

    numerator = 0.0
    x_sq_sum = 0.0
    y_sq_sum = 0.0

    for x, y in zip(x_vals, y_vals):
        dx = x - x_mean
        dy = y - y_mean
        numerator += dx * dy
        x_sq_sum += dx * dx
        y_sq_sum += dy * dy

    denominator = (x_sq_sum * y_sq_sum) ** 0.5
    if denominator == 0:
        return float('nan')

    return numerator / denominator


def fit_line(x_vals, y_vals):
    if len(x_vals) != len(y_vals) or len(x_vals) < 2:
        return None, None

    x_mean = sum(x_vals) / len(x_vals)
    y_mean = sum(y_vals) / len(y_vals)

    cov = 0.0
    var_x = 0.0
    for x, y in zip(x_vals, y_vals):
        dx = x - x_mean
        cov += dx * (y - y_mean)
        var_x += dx * dx

    if var_x == 0:
        return None, None

    slope = cov / var_x
    intercept = y_mean - slope * x_mean
    return slope, intercept


def add_fit_line(x_vals, y_vals, color, label):
    slope, intercept = fit_line(x_vals, y_vals)
    if slope is None:
        return

    x1, x2 = min(x_vals), max(x_vals)
    y1 = slope * x1 + intercept
    y2 = slope * x2 + intercept
    plt.plot([x1, x2], [y1, y2], color=color, linewidth=1.5, label=label)

# open fiel .\result\ext4.csv and read it
tab = []
gready_sko_init = []
greedy_tai_init = []
steepest_sko_init = []
steepest_tai_init = []

gready_sko_end = []
greedy_tai_end = []
steepest_sko_end = []
steepest_tai_end = []
with open('.\\result\\ex3.csv', 'r') as file:

    for line in file:
        tab.append(line.strip().split(','))

for i in range(1, len(tab)):
    # if tab[i][0] contains 'greedy'
    if 'greedy' in tab[i][0]:
        if 'sko90' in tab[i][0]:
            gready_sko_init.append(float(tab[i][9]))
            gready_sko_end.append(float(tab[i][2]))
        elif 'tai256c' in tab[i][0]:
            greedy_tai_init.append(float(tab[i][9]))
            greedy_tai_end.append(float(tab[i][2]))
    elif 'steepest' in tab[i][0]:
        if 'sko90' in tab[i][0]:
            steepest_sko_init.append(float(tab[i][9]))
            steepest_sko_end.append(float(tab[i][2]))
        elif 'tai256c' in tab[i][0]:
            steepest_tai_init.append(float(tab[i][9]))
            steepest_tai_end.append(float(tab[i][2]))

# plot greedy for sko
#  use small points
min_val = min(min(gready_sko_init), min(gready_sko_end))
max_val = max(max(gready_sko_init), max(gready_sko_end))


sko_greedy_corr = rank_correlation(gready_sko_init, gready_sko_end)
sko_steepest_corr = rank_correlation(steepest_sko_init, steepest_sko_end)
print(f'SKO greedy rank correlation r: {sko_greedy_corr:.4f}')
print(f'SKO steepest rank correlation r: {sko_steepest_corr:.4f}')

plt.scatter(gready_sko_init, gready_sko_end, label='greedy points', s=6, color='tab:blue')
plt.scatter(steepest_sko_init, steepest_sko_end, label='steepest points', s=6, color='tab:orange')
add_fit_line(gready_sko_init, gready_sko_end, 'tab:blue', f'greedy fit (r={sko_greedy_corr:.3f})')
add_fit_line(steepest_sko_init, steepest_sko_end, 'tab:orange', f'steepest fit (r={sko_steepest_corr:.3f})')

plt.legend()
plt.title('sko90')
plt.xlabel('Initial Objective Value')
plt.ylabel('Final Objective Value')
# plt.xlim(min(gready_sko_end), max(gready_sko_init))
# plt.ylim(min(gready_sko_end), max(gready_sko_init))

plt.show()

tai_greedy_corr = rank_correlation(greedy_tai_init, greedy_tai_end)
tai_steepest_corr = rank_correlation(steepest_tai_init, steepest_tai_end)
print(f'TAI greedy rank correlation r: {tai_greedy_corr:.4f}')
print(f'TAI steepest rank correlation r: {tai_steepest_corr:.4f}')

plt.scatter(greedy_tai_init, greedy_tai_end, label='greedy points', s=6, color='tab:blue')
plt.scatter(steepest_tai_init, steepest_tai_end, label='steepest points', s=6, color='tab:orange')
add_fit_line(greedy_tai_init, greedy_tai_end, 'tab:blue', f'greedy fit (r={tai_greedy_corr:.3f})')
add_fit_line(steepest_tai_init, steepest_tai_end, 'tab:orange', f'steepest fit (r={tai_steepest_corr:.3f})')

plt.legend()
plt.title('tai256c')
plt.xlabel('Initial Objective Value')
plt.ylabel('Final Objective Value')

plt.show()


