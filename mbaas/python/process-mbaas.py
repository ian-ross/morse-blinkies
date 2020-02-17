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
        seq = sequence.text_to_bit_sequence(text)
    except:
        logging.error('Error generating bit sequence', exc_info=sys.exc_info())
        return
    seqs = sequence.to_string(seq)
    print('Bit sequence:', seqs)
    info['sequence'] = seqs
    print('')

    # Padding: either a fixed padding or to next power of two.
    length = len(seq)
    if length > 256:
        print('Sequence too long: maximum length is 256')
        sys.exit(1)
    if 'explicit_padding' in rules:
        print('Explicit padding not supported: must be power of 2!')
        sys.exit(1)
        length += rules['padding']
    else:
        length = util.next_power_of_two(length)
    npadding = length - len(seq)
    print('Length:', len(seq), '->', length)
    for i in range(npadding):
        seq.append(sequence.SPACE)
    seqs = sequence.to_string(seq)
    print('Padded bit sequence:', seqs)
    info['padded_sequence'] = seqs
    print('')

    # Convert to Espresso format.
    try:
        esp = espresso.espresso(seq)
    except:
        logging.error('Error running Espresso', exc_info=sys.exc_info())
        return

    p = placement.place(esp, rules)
    a = placement.assign(p['gates'], info)
    print('')

    netlist.skidl_build(length, a, rules)
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
