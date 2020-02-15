# Remember:
#
# export KICAD_SYMBOL_DIR=/usr/share/kicad/library

from skidl import *

# LIBRARIES, FOOTPRINTS, TEMPLATES

lib_search_paths[KICAD].append('.')
lib = 'morse-blinky'

fp = {
    'R': 'Resistor_SMD:R_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'C': 'Capacitor_SMD:C_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'LED': 'LED_SMD:LED_0805_2012Metric_Pad1.15x1.40mm_HandSolder',
    'BATTERY': 'morse-blinkies:BatteryHolder_LINK_BAT-HLD-001-SMT',
    'SOIC-8': 'Package_SO:SOIC-8_3.9x4.9mm_P1.27mm',
    'SOIC-14': 'Package_SO:SOIC-14_3.9x8.7mm_P1.27mm',
    'SOIC-16': 'Package_SO:SOIC-16_3.9x9.9mm_P1.27mm'
}

r = Part('Device', 'R', TEMPLATE, footprint=fp['R'])
c = Part('Device', 'C', TEMPLATE, footprint=fp['C'])


# POWER

vdd, gnd = Net('VDD'), Net('GND')
vdd.drive = POWER
gnd.drive = POWER

batt = Part('Device', 'Battery', value='9V', footprint=fp['BATTERY'])
batt['+', '-'] += vdd, gnd


# 555 ASTABLE

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


# COUNTERS

counter_lo = Part(lib, '74HC193', footprint=fp['SOIC-16'])
counter_hi = Part(lib, '74HC193', footprint=fp['SOIC-16'])
vdd += counter_lo['VCC'], counter_hi['VCC']
gnd += counter_lo['GND'], counter_hi['GND']

# Counter load pins grounded.
print('Counter load pin grounding:')
for i in ['A', 'B', 'C', 'D']:
    print('i =', i)
    gnd += counter_lo[i], counter_hi[i]

# Active low load and down clock held high.
for i in ['LOAD', 'DOWN']:
    vdd += counter_lo[i], counter_hi[i]

# Reset drive by "32" pin on high counter.
counter_lo['CLR'] += counter_hi['QB']
counter_hi['CLR'] += counter_hi['QB']

# Counter carry/borrow.
counter_hi['CO'] += NC
counter_hi['BO'] += NC
counter_lo['CO'] += counter_hi['UP']
counter_lo['BO'] += NC

# Counter clocking.
counter_lo['UP'] += osc['OUT']

# Unconnected output pins.
counter_hi['QC'] += NC
counter_hi['QD'] += NC

# Counter bits.
bits = {
    'A': counter_lo['QA'],
    'B': counter_lo['QB'],
    'C': counter_lo['QC'],
    'D': counter_lo['QD'],
    'E': counter_hi['QA']
}
bit_names = ['A', 'B', 'C', 'D', 'E']


# INVERTERS

inverters = Part('morse-blinky', '74HC04', footprint=fp['SOIC-14'])
inverters['VCC'] += vdd
inverters['GND'] += gnd

for i in bit_names:
    inverters[i + 'IN'] += bits[i]
    bits['~' + i] = Net()
    bits['~' + i] += inverters[i + 'OUT']

inverters['FIN'] += gnd
inverters['FOUT'] += NC


# LOGIC

and4 = Part('morse-blinky', '74HC21', footprint=fp['SOIC-14'])
and4['VCC'] += vdd
and4['GND'] += gnd

and3 = Part('morse-blinky', '74HC15', footprint=fp['SOIC-14'])
and3['VCC'] += vdd
and3['GND'] += gnd

or2 = Part('morse-blinky', '74HC32', footprint=fp['SOIC-14'])
or2['VCC'] += vdd
or2['GND'] += gnd


and4['AIN1'] += bits['B']
and4['AIN2'] += bits['C']
and4['AIN3'] += bits['D']
and4['AIN4'] += bits['~E']

or2['AIN1'] += and4['AOUT']

and3['AIN1'] += bits['~A']
and3['AIN2'] += bits['~C']
and3['AIN3'] += bits['~D']

or2['AIN2'] += and3['AOUT']

and4['BIN1'] += bits['~B']
and4['BIN2'] += bits['~C']
and4['BIN3'] += bits['D']
and4['BIN4'] += bits['~E']

or2['BIN1'] += and4['BOUT']

and3['BIN1'] += bits['~A']
and3['BIN2'] += bits['B']
and3['BIN3'] += bits['~E']

or2['BIN2'] += and3['BOUT']

or2['CIN1'] += or2['AOUT']
or2['CIN2'] += or2['BOUT']


gnd += and3['CIN1'], and3['CIN2'], and3['CIN3']
and3['COUT'] += NC

gnd += or2['DIN1'], or2['DIN2']
or2['DOUT'] += NC


# BLINKY

r_led = r(value='1K')
led = Part('Device', 'LED', footprint=fp['LED'])
led['K'] += r_led[1]
r_led[2] += gnd

led['A'] += or2['COUT']
