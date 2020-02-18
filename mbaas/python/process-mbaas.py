#!/usr/bin/env python

import json
import logging
import os
import sys

import skidl

import espresso
import netlist
import placement
import sequence
import util


# make_board is the main driver for constructing a board.
def make_board(text, rules):
    print('Making board for:', text)
    print('')

    info = {}

    # Convert to Morse mark/space sequence.
    text = text.upper()
    info['text'] = text
    try:
        seqs = sequence.text_to_bit_sequence(text, rules['type'])
    except:
        logging.error('Error generating bit sequence', exc_info=sys.exc_info())
        return
    print('Bit sequence:')
    info['sequence'] = []
    for seq in seqs:
        s = sequence.to_string(seq)
        print(' ', s)
        info['sequence'].append(s)
    nseqs = len(seqs)
    print('')

    # Padding: either a fixed padding or to next power of two.
    length = len(seqs[0])
    if length > 256:
        print('Sequence too long: maximum length is 256')
        sys.exit(1)
    length = util.next_power_of_two(length)
    npadding = length - len(seqs[0])
    print('Length:', len(seqs[0]), '->', length)
    for i in range(npadding):
        for j in range(nseqs):
            seqs[j].append(sequence.SPACE)
    print('Padded bit sequence:')
    info['padded_sequence'] = []
    for seq in seqs:
        s = sequence.to_string(seq)
        print(' ', s)
        info['padded_sequence'].append(s)
    print('')

    # Convert to Espresso format.
    try:
        esp = espresso.espresso(seqs)
    except:
        logging.error('Error running Espresso', exc_info=sys.exc_info())
        return

    p = placement.place(esp, rules)
    print(p)
    c, a = placement.assign(p['gates'], info)
    print('')

    netlist.skidl_build(nseqs, length, c, a, rules)
    skidl.ERC()
    skidl.generate_netlist()

    return info


# Command line:
#
#  process-mbaas <text> [<rules.json file>]

if __name__ == '__main__':
    if len(sys.argv) < 2 or len(sys.argv) > 3:
        usage()
    text = sys.argv[1]
    rules_file = os.path.join(os.getenv('MBAAS_BASE_DIR'), 'lib/default_rules.json')
    if len(sys.argv) == 3:
        rules_file = sys.argv[2]
    try:
        with open(rules_file) as fp:
            rules = json.load(fp)
    except:
        logging.error('Error reading rules file', exc_info=sys.exc_info())
        sys.exit(1)

    info = make_board(text, rules)
    with open("process-mbaas.info", "w") as fp:
        json.dump(info, fp, indent=2)
