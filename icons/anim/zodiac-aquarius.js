import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="m2 10 "]'), el.querySelector('path[d^="m2 18.002"]')];
  animate(waves, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });
}
