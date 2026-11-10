import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelector('path[d^="m6 8"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const straw = el.querySelector('path[d^="m12 8"]');
  animate(straw, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
