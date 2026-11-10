import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll('path[d^="M18 8"], path[d^="m6 8"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
  const line = el.querySelector('path[d^="M18 4"]');
  animate(line, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
