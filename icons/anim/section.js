import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const strokes = [el.querySelector('path[d^="M16 5"]'), el.querySelector('path[d^="M8 19"]')];
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.12) });
}
