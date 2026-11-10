import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const upper = el.querySelector('path[d^="M14 4.5"]');
  animate(upper, { translateX: [0, 1.5, 0], translateY: [0, -1.5, 0] }, { duration: τmid });
  const lower = el.querySelector('path[d^="M9.698"]');
  animate(lower, { translateX: [0, -1.5, 0], translateY: [0, 1.5, 0] }, { duration: τmid });

  const chips = el.querySelectorAll('path[d^="m"]');
  animate(chips, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
