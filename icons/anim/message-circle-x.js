import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M2.992"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m15 9"]'), el.querySelector('path[d^="m9 9"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
