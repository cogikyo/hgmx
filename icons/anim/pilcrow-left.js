import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelector('path[d^="M14 9"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const stems = el.querySelectorAll('path[d$="v11"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
  const arrow = el.querySelectorAll('path[d^="M22 18"], path[d^="m6 22"]');
  animate(arrow, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
