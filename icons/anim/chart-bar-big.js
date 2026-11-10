import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('rect[y="13"]'), el.querySelector('rect[y="5"]')];
  animate(bars, { scaleX: [0, 1.1, 1] }, { duration: τmid, delay: stagger(0.12, { startDelay: τmin }) });
}
