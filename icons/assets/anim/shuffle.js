import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const strand = el.querySelector('path[d^="M2 18"]');
  animate(strand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const heads = [el.querySelector('path[d^="m18 14"]'), el.querySelector('path[d^="m18 2"]')];
  animate(heads, { translateX: [0, 4, -1, 0] }, { duration: τmid, delay: stagger(0.12, { startDelay: τmin }) });
}
