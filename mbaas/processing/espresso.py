import subprocess
import tempfile

import util


DONTCARE = 0
POS = 1
NEG = 2

def espresso(seq):
    # Generate Espresso file.
    with tempfile.NamedTemporaryFile(mode='w') as fp:
        write_input_file(fp, seq)

        # Run Espresso and read output.
        cp = subprocess.run(['espresso', fp.name], capture_output=True)
        if cp.returncode != 0:
            raise Exception('error from Espresso: ' + cp.stderr.decode())
        espresso = cp.stdout.decode()

        # Process Espresso output.
        return convert_output(espresso.split("\n"))


def write_input_file(fp, seq):
    nbits = util.number_of_bits(len(seq) - 1)

    print('.i', nbits, file=fp)
    print('.o 1', file=fp)

    fmt = '{0:0' + str(nbits) + 'b}'
    for i in range(len(seq)):
        print(fmt.format(i), '1' if seq[i] else '0', file=fp)

    fp.flush()


FACTORS = {
    '-': DONTCARE,
    '1': POS,
    '0': NEG,
}

def convert_output(lines):
    result = []
    for line in lines:
        if line.strip() == '':
            continue
        if line[0] == '.':
            continue
        bits = line.split(' ')[0]
        result.append([FACTORS[ch] for ch in bits])

    return result
