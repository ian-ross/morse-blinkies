from espresso import DONTCARE, POS, NEG

NOT = 0
AND = 1
OR = 2

# Gates are represented as [<type>, [<input nets>], <output net>]

net_idx = 0
andint_net_idx = 0
net_names = {}
gates = []
pos_nets = {}
neg_nets = {}

def place(exp, rules):
    global net_idx
    global net_names
    global gates
    global pos_nets
    global neg_nets

    # Assign net names to inputs (a, b, c, ... starting at LSB).
    ninputs = len(exp[0])
    for i in range(ninputs):
        n = input_name(i)
        net_names[n] = net_idx
        pos_nets[i] = n
        net_idx += 1

    # Identify variable positions with negations (to size number of
    # inverters).
    negs = {}
    for term in exp:
        for i in range(ninputs):
            if term[i] == NEG:
                negs[i] = True

    # Generate inverter gates and assign net names to negated inputs
    # (~a, ~b, ~c, ... starting at LSB).
    for i in sorted(negs.keys()):
        pos_net = input_name(i)
        neg_net = '~' + pos_net
        net_names[neg_net] = net_idx
        neg_nets[i] = neg_net
        gates.append([NOT, [pos_net], neg_net])
        net_idx += 1

    # Process terms.
    and_idx = 1
    and_nets = []
    for term in exp:
        # Use appropriate strategy to evaluate AND using only 4-, 3-,
        # or 2-input gates.
        out_net = 'and' + str(and_idx)
        and_nets.append(out_net)
        and_idx += 1
        net_names[out_net] = net_idx
        net_idx += 1
        make_and(term, out_net)

    out_net = 'output'
    net_names[out_net] = net_idx
    net_idx += 1
    make_or_tree(and_nets, out_net)

    return dict(nets=net_names, gates=gates)


def make_or_tree(in_nets, out_net):
    global net_idx
    global net_names
    global gates

    int_net_idx = 1
    while len(in_nets) > 2:
        int_net = 'orint' + str(int_net_idx)
        int_net_idx += 1
        net_names[int_net] = net_idx
        net_idx += 1
        gates.append([OR, [in_nets[0], in_nets[1]], int_net])
        in_nets = in_nets[2:]
        in_nets.append(int_net)
    gates.append([OR, [in_nets[0], in_nets[1]], out_net])

def make_and(term, out_net):
    global net_idx
    global pos_nets
    global neg_nets
    global gates

    # Count number of inputs to AND.
    in_nets = []
    andn = 0
    in_idx = -1
    for f in term:
        in_idx += 1
        if f == DONTCARE:
            continue
        andn += 1
        in_nets.append(pos_nets[in_idx] if f == POS else neg_nets[in_idx])

    make_and_tree(in_nets, out_net)

def make_and_tree(in_nets, out_net):
    global net_idx
    global andint_net_idx
    global net_names
    global gates

    take = []
    while len(in_nets) > 4:
        andint_net_idx += 1
        int_net = 'andint' + str(andint_net_idx)
        net_names[int_net] = net_idx
        net_idx += 1
        if len(take) == 0:
            take = [4]
        else:
            take = [4]
            # Treat 5 as (a & b & c) & d & e = 2 x 3.
            # Treat 6 as (a & b & c) & (d & e & f) = 2 x 3 + 1 x 2/3
            if len(in_nets) == 5 or len(in_nets) == 6:
                take = [3, 3]
        gate_in_nets = in_nets[:take[0]]
        in_nets = in_nets[take[0]:]
        take = take[1:]
        gates.append([AND, gate_in_nets, int_net])
        in_nets.append(int_net)

    gates.append([AND, in_nets, out_net])


def input_name(i):
    return chr(ord('a') + i)


def assign(gates, info):
    not_gates = [g for g in gates if g[0] == NOT]
    and2_gates = [g for g in gates if g[0] == AND and len(g[1]) == 2]
    and3_gates = [g for g in gates if g[0] == AND and len(g[1]) == 3]
    and4_gates = [g for g in gates if g[0] == AND and len(g[1]) == 4]
    or_gates = [g for g in gates if g[0] == OR]

    chips_7404, dum = alloc_div('74HC04', not_gates, 6)
    chips_7432, dum = alloc_div('74HC32', or_gates, 4)
    chips_7421, left_over = alloc_div('74HC21', and4_gates, 2)
    if left_over > 0:
        chips_7421, and3_gates = use_left_overs(chips_7421, left_over, and3_gates, 4, 3)
    chips_7411, left_over = alloc_div('74HC11', and3_gates, 3)
    if left_over > 0:
        chips_7411, and2_gates = use_left_overs(chips_7411, left_over, and2_gates, 3, 2)
    chips_7408, dum = alloc_div('74HC08', and2_gates, 4)

    gate_info = {}
    if len(chips_7404) > 0:
        print('7404 hex inverter:      ', len(chips_7404))
        gate_info['7404 hex inverter'] = len(chips_7404)
    if len(chips_7408) > 0:
        print('7408 quad 2-input AND:  ', len(chips_7408))
        gate_info['7408 quad 2-input AND'] = len(chips_7408)
    if len(chips_7411) > 0:
        print('7411 triple 3-input AND:', len(chips_7411))
        gate_info['7408 triple 3-input AND'] = len(chips_7411)
    if len(chips_7421) > 0:
        print('7421 dual 4-input AND:  ', len(chips_7421))
        gate_info['7421 dual 4-input AND'] = len(chips_7421)
    if len(chips_7432) > 0:
        print('7432 quad 2-input OR:   ', len(chips_7432))
        gate_info['7432 quad 2-input OR'] = len(chips_7432)
    info['chips'] = gate_info

    return chips_7404 + chips_7432 + chips_7408 + chips_7411 + chips_7421


def alloc_div(chip, gates, per_chip):
    chips = []
    left_over = 0
    while len(gates) > 0:
        alloc = gates[:per_chip]
        if len(alloc) < per_chip:
            left_over = per_chip - len(alloc)
        chips.append([chip, alloc])
        gates = gates[per_chip:]
    return chips, left_over


def use_left_overs(chips_in, left_over, gates_in, chip_gate_size, gate_size):
    extras = (chip_gate_size - gate_size) * ['vdd']
    adjust_gates = gates_in[:left_over]
    for i in range(len(adjust_gates)):
        adjust_gates[i][1] += extras
    chips_in[-1][1] += adjust_gates
    gates_out = gates_in[left_over:]
    return chips_in, gates_out
