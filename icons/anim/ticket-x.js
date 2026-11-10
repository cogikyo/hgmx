import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticket = el.querySelector('path[d^="M2 9"]');
  animate(ticket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="m9.5 14.5 5-5"]'), el.querySelector('path[d="m9.5 9.5 5 5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}
