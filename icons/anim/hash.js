import { animate, stagger } from "../motion.js";

export default function (el) {
  const strokes = ['[x1="10"]', '[x1="16"]', '[y1="9"]', '[y1="15"]'].map((at) => el.querySelector(`line${at}`));
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
