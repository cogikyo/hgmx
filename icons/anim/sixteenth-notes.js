import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const notes = el.querySelector("g");
  animate(notes, { rotate: [0, -10, 8, -6, 0] }, { duration: τmax });

  const heads = [el.querySelector('circle[cx="6"]'), el.querySelector('circle[cx="18"]')];
  animate(heads, { scale: [1, 1.25, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
