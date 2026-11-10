import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const strokes = el.querySelectorAll('path[d$="v16"]');
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const slash = el.querySelector('path[d^="M22 6"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
