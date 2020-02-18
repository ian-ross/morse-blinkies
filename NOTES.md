 - CE Getting to Blinky: idea to do something a bit extra

 - How to do Morse blinkies?
    * 3-cent microcontroller (boring)
    * Serial EEPROM (or something like, also boring)
    * 555s + discretes only (hard! maybe impossible?)
    * 555 + 74xx logic (silly, but let's give it a go)

 - Manual examples:
    * Repeated letter A: repeating 8-bit sequence => counter + logic
      on counter bits
    * "IAN": logic analysis is already pretty heavy => use Espresso to
      simplify
    * "KICAD": practical number of 74xx chips?

 - Thoughts:
    * Need to restrict to easily available 74xx components.
    * Write code to do:
        + ASCII input -> Morse bit sequence
        + Bit sequence -> minimised logic
        + Minimised logic -> allocation of gates to 74xx devices
        + Allocation -> Kicad netlist
    * Then do layout by hand...
    * Provide it all as a web app? MBaaS?

 - What's easy?
    * ASCII -> Morse bit sequence
    * Logic minimisation (Espresso)
    * Netlist generation (SKIDL)
    * Web app (Go)

 - What's harder?
    * Placement of gates on 74xx devices

 - To do placement, inputs needed are:
    * Espresso output with minimised logic in DNF.
    * What devices to use in generated SKIDL code (footprints,
      discretes for 555 timer, etc.) => have a rules.json file,
      provide a couple of predefined ones, and allow users to download
      template, edit and upload their own. [FIX THIS TO BEGIN WITH]
    * Also single LED/one LED per letter flag? [DO LATER]


Tasks:

1. Start Go processing code: word -> bit sequence, Espresso format,
   Espresso minimisation.
2. Plan placement code.
3. Look at counter setup using 74193.
4. Replicate "IAN" blinky using SKIDL by hand.
5. Basic placement code.
6. Netlist generation: SKIDL code generation, run SKIDL.
7. Improve placement code.
8. Options for placement code.
9. Web app front-end.
10. Features: multi-LED, timing, etc.
11. Queuing for backend jobs.




Process:

 - Hand calculation experiments
 - Logisim simulations
 - Start with single LED idea
 - Morse code stuff: text -> bit sequence
 - SKIDL experiments and framework
 - Component selection: 74xx gates and counters
 - Logic placement first pass
 - Web app first pass
 - Semi-working...
 - Multi-LED idea: web app, logic placement, ...
 - Oops.  Drive current for 5 LEDs is ca. 100 mA!  MOSFETs?
 - Mouser has lots of MOSFETS: parametric selection => jellybean


