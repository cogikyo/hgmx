import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="m10 7"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const marks = el.querySelectorAll('path[d="m7 10 2 2"], path[d="m4 13 2 2"]');
  animate(marks, { scale: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
