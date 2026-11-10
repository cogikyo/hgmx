import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M4 8"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d^="m9.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });

  const lid = el.querySelector("rect");
  animate(lid, { translateY: [0, -2, 0] }, { duration: τmin, delay: τmid });
}
