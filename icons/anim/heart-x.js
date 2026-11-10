import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const heart = el.querySelector('path[d^="M21.955"]');
  animate(heart, { scale: [1, 1.25, 0.85, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m20.5"]'), el.querySelector('path[d^="m15.5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
