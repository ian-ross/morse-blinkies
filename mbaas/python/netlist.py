# Remember:
#
# export KICAD_SYMBOL_DIR=/usr/share/kicad/library

from skidl import *

# LIBRARIES, FOOTPRINTS, TEMPLATES

lib_search_paths[KICAD].append('../libraries')
lib = 'morse-blinky'

fp = {
    'R': 'Resistor_SMD:R_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'C': 'Capacitor_SMD:C_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'LED': 'LED_SMD:LED_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'BATTERY': 'morse-blinky:BatteryHolder_LINK_BAT-HLD-001-SMT',
    'SOIC-8': 'Package_SO:SOIC-8_3.9x4.9mm_P1.27mm',
    'SOIC-14': 'Package_SO:SOIC-14_3.9x8.7mm_P1.27mm',
    'SOIC-16': 'Package_SO:SOIC-16_3.9x9.9mm_P1.27mm'
}

units = {
    '74HC04': 6,
    '74HC08': 4,
    '74HC11': 3,
    '74HC21': 2,
    '74HC32': 4
}

r = Part('Device', 'R', TEMPLATE, footprint=fp['R'])
c = Part('Device', 'C', TEMPLATE, footprint=fp['C'])


def skidl_build(length, gates, rules):
    vdd, gnd = power()
    osc = oscillator(vdd, gnd, rules['blink_rate_ms'])
    nets, bit_names = counters(vdd, gnd, osc, length)
    out = logic(vdd, gnd, gates, nets)
    blinky(vdd, gnd, out)


# POWER

def power():
    print('SKIDL: power')
    vdd, gnd = Net('VDD'), Net('GND')
    vdd.drive = POWER
    gnd.drive = POWER

    batt = Part('Device', 'Battery', value='9V', footprint=fp['BATTERY'])
    batt['+', '-'] += vdd, gnd

    return vdd, gnd


# 555 ASTABLE

# TODO: PARAMETERISE OSCILLATOR RATE
def oscillator(vdd, gnd, blink_rate_ms):
    print('SKIDL: oscillator')
    osc = Part(lib, '7555', footprint=fp['SOIC-8'])
    osc['VDD'] += vdd
    osc['GND'] += gnd
    osc['RESET'] += vdd
    osc['DIS'] += NC

    c_ctrl = c(value='100n')
    c_ctrl[1, 2] += osc['CTRL'], gnd

    r_trig = r(value='470K')
    c_trig = c(value='1U')

    c_trig[1] += gnd
    c_trig[2] += r_trig[1], osc['THRESH'], osc['TRIG']
    r_trig[2] += osc['OUT']

    return osc['OUT']


# COUNTERS

# TODO: PARAMETERISE BASED ON SEQUENCE LENGTH.
#
#  * IF LENGTH <= 16, USE ONE TIMER, ELSE USE TWO.
#  * TIMER RESET

def counters(vdd, gnd, osc, length):
    print('SKIDL: counters')
    counter_lo = Part(lib, '74HC193', footprint=fp['SOIC-16'])
    vdd += counter_lo['VCC']
    gnd += counter_lo['GND']
    counter_hi = None
    if length > 16:
        counter_hi = Part(lib, '74HC193', footprint=fp['SOIC-16'])
        vdd += counter_hi['VCC']
        gnd += counter_hi['GND']

    # Counter load pins grounded.
    for i in ['A', 'B', 'C', 'D']:
        gnd += counter_lo[i]
        if counter_hi is not None:
            gnd += counter_hi[i]

    # Active low load and down clock held high.
    for i in ['LOAD', 'DOWN']:
        vdd += counter_lo[i]
        if counter_hi is not None:
            vdd += counter_lo[i], counter_hi[i]

    # Counter carry/borrow.
    if counter_hi is not None:
        counter_hi['CO'] += NC
        counter_hi['BO'] += NC
        counter_lo['CO'] += counter_hi['UP']
        counter_lo['BO'] += NC
    else:
        counter_lo['CO'] += NC
        counter_lo['BO'] += NC

    # Counter clocking.
    counter_lo['UP'] += osc

    # Work out which pin to use to drive the counter reset. The
    # sequence length is a power of two, so only one bit is needed to
    # trigger the reset.
    reset_bit = 0
    reset_power = 1
    while length > reset_power:
        reset_bit += 1
        reset_power *= 2
    bit_count = reset_bit
    # Skip cases where the single or double counter is allowed to roll
    # over.
    if reset_bit != 4 and reset_bit != 8:
        reset_counter = counter_lo
        if reset_bit >= 4:
            reset_counter = counter_hi
            reset_bit -= 4
        reset_bit_name = 'Q' + chr(ord('A') + reset_bit)
        counter_lo['CLR'] += reset_counter[reset_bit_name]
        if counter_hi is not None:
            counter_hi['CLR'] += reset_counter[reset_bit_name]

        # Unconnected output pins.
        bit = reset_bit + 1
        while bit < 4:
            bit_name = 'Q' + chr(ord('A') + bit)
            reset_counter[bit_name] += NC
            bit += 1

    # Counter bits.
    bits = {}
    bit_names = []
    for bit in range(bit_count):
        if bit <= 3:
            counter = counter_lo
        else:
            counter = counter_hi
        bit_name = chr(ord('a') + bit)
        counter_bit_name = 'Q' + chr(ord('A') + bit % 4)
        bits[bit_name] = counter[counter_bit_name]
        bit_names.append(bit_name)

    return bits, bit_names


def logic(vdd, gnd, gates, nets):
    print('SKIDL: logic')
    parts_with_gates = []
    for chip in gates:
        part = Part('morse-blinky', chip[0], footprint=fp['SOIC-14'])
        part['VCC'] += vdd
        part['GND'] += gnd
        added = 0
        for i in range(len(chip[1])):
            if len(chip[1][i]) == 0:
                continue
            parts_with_gates.append([part, i, chip[1][i][1], chip[1][i][2], chip[0]])
            added += 1
        dummy = [None for i in range(len(chip[1][0][1]))]
        unit = added
        while unit < units[chip[0]]:
            parts_with_gates.append([part, unit, dummy, None, chip[0]])
            unit += 1

    last_len = len(parts_with_gates) + 1
    did_another = False
    while len(parts_with_gates) != last_len and not did_another:
        if len(parts_with_gates) == last_len:
            did_another = True
        last_len = len(parts_with_gates)

        for iattempt in range(len(parts_with_gates)):
            attempt = parts_with_gates[iattempt]
            part = attempt[0]
            unit = attempt[1]
            in_nets = attempt[2]
            ninputs = len(in_nets)
            out_net = attempt[3]
            chip = attempt[4]

            if out_net is None:
                if ninputs == 1:
                    part[chr(ord('A') + unit) + 'IN'] += NC
                else:
                    for i in range(ninputs):
                        part[chr(ord('A') + unit) + 'IN' + str(i+1)] += NC
                part[chr(ord('A') + unit) + 'OUT'] += NC
            else:
                if not all([net in nets or net == 'vdd' or net == 'gnd'
                            for net in in_nets]):
                    continue

                if ninputs == 1:
                    in_pins = [chr(ord('A') + unit) + 'IN']
                else:
                    in_pins = [chr(ord('A') + unit) + 'IN' + str(i+1) for i in range(ninputs)]
                out_pin = chr(ord('A') + unit) + 'OUT'
                for i in range(ninputs):
                    if in_nets[i] == 'vdd':
                        vdd += part[in_pins[i]]
                    elif in_nets[i] == 'gnd':
                        gnd += part[in_pins[i]]
                    else:
                        nets[in_nets[i]] += part[in_pins[i]]
                    print('  ', in_nets[i], '->', chip, part.ref, in_pins[i])
                nets[out_net] = Net()
                nets[out_net] += part[out_pin]
                print('  ', out_net, '->', chip, part.ref, out_pin)

            del parts_with_gates[iattempt]
            break

    if len(parts_with_gates) != 0:
        print("Something went wrong -- didn't place all gates!")
        print(nets)
        print('---')
        print(parts_with_gates)
        sys.exit(1)

    return nets['output']


# BLINKY

def blinky(vdd, gnd, output):
    print('SKIDL: blinky')
    r_led = r(value='1K')
    led = Part('Device', 'LED', footprint=fp['LED'])
    led['K'] += r_led[1]
    r_led[2] += gnd

    led['A'] += output
