import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const meat = el.querySelector('path[d^="M15.3"]');
  animate(meat, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const marks = el.querySelectorAll('path[d^="m10.3"], path[d^="m9 15"]');
  animate(marks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
