# coding:utf-8
from collections import namedtuple
import random
import time
import signal
import sys
import curses
import locale
import math
import os

# def handler(signum):
#     os.exit("Bye!")
# signal.signal(signal.SIG_IGN, handler)

locale.setlocale(locale.LC_ALL, '')
code = locale.getpreferredencoding()

stdscr = curses.initscr()
curses.noecho()
curses.cbreak()
stdscr.keypad(True)

Symbol = namedtuple('Symbol', 'symbol desc')

symbols = [
    # 母音
    Symbol("æ", "[母音] apple ([a]pple): エとアを混ぜあわせた音"),
    Symbol("ǝ",  "[母音] bird (b[ir]d) 曖昧母音。こもったタイプのア"),
    Symbol("ʌ",  "[母音] cup (c[u]p) 口を小さく開く鋭いア"),
    Symbol("ɑ", "[母音] god (g[o]d) 口を大きく開くア"),
    Symbol("ɔ",  "[母音] god (g[o]d) オに当たる母音。アメリカ英語では通常 a と発音される"),
    Symbol("ɛǝ", "[母音] chair (ch[ai]r) エとアが合わさった二重母音"),
    Symbol("r",  "[母音] door (do[or]) r音"),
    Symbol("ː", "[母音] 音を伸ばす記号 door (do[or])"),
    Symbol("i, u, e, ou",  "[母音] そのまま (母音)"),
    # 子音
    Symbol("tʃ", "[子音] watch (wa[tch]) チに当たる音"),
    Symbol("dʒ", "[子音] judge (ju[dge]) 鋭いヂ"),
    Symbol("θ", "[子音] three ([th]ree) 舌を前歯で軽く挟む澄んだタイプの音"),
    Symbol("ð", "[子音] these ([th]ese) 舌を前歯で軽く挟む濁ったタイプの音"),
    Symbol("ʃ", "[子音] dish (di[sh]) シュに当たる音"),
    Symbol("ʒ", "[子音] measure (mea[s]ure) ジュに当たる音"),
    Symbol("ŋ", "[子音] thing (thi[ng]) ン(グ)/グははっきり発音しない"),
    Symbol("j", "[子音] yes ([y]es) y に当たる音。ローマ字の j とは違う"),
    Symbol("p, b, t, d, k, g, s, z, h, l, r, w, f, v, m, n",  "[子音] そのまま"),
]

def wait_animation(stdscr, sec, y, x):
    prompt = "待ち時間 {} 秒: ".format(sec)
    stdscr.addstr(y, x, prompt)
    stdscr.refresh()
    x = x + 15
    for i in range(sec):
        time.sleep(1)
        stdscr.addstr(y, x, str(i+1))
        stdscr.refresh()
        x += 1


def main(stdscr):
    _symbols = symbols[:]
    total_count = 1
    while True:
        # stdscr.clear()
        # stdscr.addstr(0, 0, "Exit: q, Next: hit any key.")
        # stdscr.refresh()

        # c = stdscr.getch()
        # if c == ord('q'):
        #     break  # Exit the while loop

        stdscr.clear()

        # question and answer
        random.shuffle(_symbols)
        q_symbol = _symbols.pop()
        if len(_symbols) == 0:
            # copy working list from original list again
            _symbols = symbols[:]
            stdscr.addstr(0, 0, str(total_count)+"周終わり！") # display question
            total_count += 1
            stdscr.getch()
            stdscr.refresh()

        # determine which is question
        qa_flg = random.randrange(2)
        if qa_flg == 1:
            q = q_symbol.symbol
            a = q_symbol.desc
        else:
            q = q_symbol.desc
            a = q_symbol.symbol

        stdscr.addstr(0, 0, q) # display question
        stdscr.refresh()
        wait_animation(stdscr, sec=8, y=1, x=0) # waiting
        stdscr.addstr(2, 4, a)
        stdscr.refresh()
        wait_animation(stdscr, sec=8, y=3, x=0) # interval before next question

if __name__ == '__main__':
    try:
        curses.wrapper(main)
    except KeyboardInterrupt:
        print("Bye!")
