import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M11 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d^="m15 17"], path[d^="m20 17"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmin }) });
}
