import subprocess
import tempfile

import util


DONTCARE = 0
POS = 1
NEG = 2

def espresso(seqs):
    # Generate Espresso file.
    with tempfile.NamedTemporaryFile(mode='w') as fp:
        write_input_file(fp, seqs)

        # Run Espresso and read output.
        cp = subprocess.run(['espresso', fp.name],
                            stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if cp.returncode != 0:
            raise Exception('error from Espresso: ' + cp.stderr.decode())
        espresso = cp.stdout.decode()

        # Process Espresso output.
        return convert_output(espresso.split("\n"))


def write_input_file(fp, seqs):
    nbits = util.number_of_bits(len(seqs[0]) - 1)

    print('.i', nbits, file=fp)
    print('.o', len(seqs), file=fp)

    fmt = '{0:0' + str(nbits) + 'b}'
    for i in range(len(seqs[0])):
        print(fmt.format(i), ''.join(['1' if s[i] else '0' for s in seqs]), file=fp)

    fp.flush()


FACTORS = {
    '-': DONTCARE,
    '1': POS,
    '0': NEG,
}

def convert_output(lines):
    nseqs = len([l for l in lines if l.strip() != '' and l[0] != '.'][0].split(' ')[1])
    result = [[] for i in range(nseqs)]
    for line in lines:
        if line.strip() == '':
            continue
        if line[0] == '.':
            continue
        bits, places = line.split(' ')
        rep = [FACTORS[ch] for ch in bits]
        for i in range(nseqs):
            if places[i] == '1':
                result[i].append(rep)

    return result
