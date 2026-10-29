import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h8"]');
  animate(lines, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const ticks = [el.querySelector('path[d^="m3 7"]'), el.querySelector('path[d^="m3 17"]')];
  animate(ticks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });
  animate(ticks, { scale: [1, 1.25, 1] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmid }) });
}
