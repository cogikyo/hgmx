import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const water = el.querySelector('path[d^="M2 18"]');
  animate(water, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: 0.3 });
  const rungs = [el.querySelector('path[d^="M7 13"]'), el.querySelector('path[d^="M7 9"]')];
  animate(rungs, { scaleX: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
