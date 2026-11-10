import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelector('path[d^="M19 4"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const stems = el.querySelectorAll('path[d$="v16"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
