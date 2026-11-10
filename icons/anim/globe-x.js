import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelector('path[d^="M2 12h20"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m16 3 "]'), el.querySelector('path[d^="m21 3-"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
