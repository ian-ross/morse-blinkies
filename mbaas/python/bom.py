import csv
import string
from kinparse import parse_netlist

def refkey(ref):
    letters = [l for l in ref if l in string.ascii_letters]
    digits = [d for d in ref if d in string.digits]
    if len(digits) == 1:
        digits = ['0'] + digits
    return ''.join(letters + digits)

def make_bom():
    nl = parse_netlist('process-mbaas.net')

    parts = list(map(process_part, nl.parts))

    part_groups = {}
    for p in parts:
        ref = ''.join([l for l in p['ref'][0] if l in string.ascii_letters])
        k = ref + '-' + p['val'] + '-' + p['fp']
        if p['desc'] is not None:
            k += '-' + p['desc']
        if k in part_groups:
            part_groups[k] += [p]
        else:
            part_groups[k] = [p]

    bom_rows = {}
    for k in part_groups:
        pg = part_groups[k]
        row_key = '-'.join(sorted([refkey(p['ref']) for p in pg]))
        refs = ', '.join(sorted([p['ref'] for p in pg], key=refkey))
        if pg[0]['desc'] is None:
            name = pg[0]['val']
            val = ''
        else:
            name = pg[0]['desc']
            val = pg[0]['val']
        bom_rows[row_key] = dict(refs=refs, name=name, val=val, qty=len(pg), fp=pg[0]['fp'])

    with open('bom.csv', 'w', newline='') as csvfile:
        w = csv.writer(csvfile)
        w.writerow(['References', 'Component', 'Value', 'Qty', 'Footprint'])
        for k in sorted(bom_rows.keys()):
            row = bom_rows[k]
            w.writerow([row['refs'], row['name'], row['val'], row['qty'], row['fp']])

def process_part(p):
    ref = p.ref # Reference designator.
    desc = None # Description.
    for f in p.fields:
        if f[0] == 'description':
            desc = f[1]
            break
    val = p.value # Value.
    fp = p.footprint # PCB footprint.
    return dict(ref=ref, desc=desc, val=val, fp=fp)
